package server

import (
	"crypto/tls"
	log "github.com/sirupsen/logrus"
	"net"
	"time"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcprom "github.com/grpc-ecosystem/go-grpc-prometheus"
	grpcot "github.com/opentracing-contrib/go-grpc"
	ot "github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const maxConnectionAge = 60 * time.Second

// Config gRPC server options.
type Config struct {
	Address              string
	CertFilename         string
	KeyFilename          string
	ClientCACert         string
	ReflectionAPIEnabled bool
	TraceEnabled         bool
}

// DefaultConfig default gRPC server options.
var DefaultConfig = Config{
	Address:              "localhost:8443",
	CertFilename:         "test-certs/server-cert.pem",
	KeyFilename:          "test-certs/server-key.pem",
	ClientCACert:         "test-certs/ca-cert.pem",
	ReflectionAPIEnabled: true,
	TraceEnabled:         false,
}

// ShopServer is server where gRPC services can be registered in.
type ShopServer struct {
	Addr       string
	grpcServer *grpc.Server
}

// New returns initialized grpc server.
func New(opts Config, tls *tls.Config) *ShopServer {
	s := new(ShopServer)
	s.Addr = opts.Address

	grpcprom.EnableHandlingTimeHistogram()
	logEntry := log.NewEntry(log.New())

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		grpcprom.UnaryServerInterceptor,
		grpc_logrus.UnaryServerInterceptor(logEntry),
	}
	streamInterceptors := []grpc.StreamServerInterceptor{
		grpcprom.StreamServerInterceptor,
		grpc_logrus.StreamServerInterceptor(logEntry),
	}

	if opts.TraceEnabled {
		tracer := ot.GlobalTracer()
		unaryInterceptors = append(unaryInterceptors, grpcot.OpenTracingServerInterceptor(tracer))
		streamInterceptors = append(streamInterceptors, grpcot.OpenTracingStreamServerInterceptor(tracer))
	}

	serverOptions := []grpc.ServerOption{
		grpc.Creds(credentials.NewTLS(tls)),
		grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(unaryInterceptors...)),
		grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(streamInterceptors...)),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			// sends each 60s kill msg so client reloads the servers and can add new instances.
			MaxConnectionAge: maxConnectionAge,
		}),
	}
	s.grpcServer = grpc.NewServer(serverOptions...)

	grpcprom.Register(s.grpcServer)

	if opts.ReflectionAPIEnabled {
		reflection.Register(s.grpcServer)
		log.Info("Reflection API is active.")
	}

	return s
}

// RegisterService implements grpc.ServiceRegistrar interface so internals of this type does not need to be exposed.
func (s *ShopServer) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	s.grpcServer.RegisterService(desc, impl)
}

// ListenAndServe gRPC server starts listening on given address including the port.
func (s *ShopServer) ListenAndServe() error {
	log.Infof("Starting gRPC server on address '%s'.", s.Addr)

	lis, err := net.Listen("tcp", s.Addr)
	if err != nil {
		return err
	}

	if err := s.grpcServer.Serve(lis); err != nil {
		return err
	}
	return nil
}

// GracefulShutdown gracefully shutdowns the gRPC server.
func (s *ShopServer) GracefulShutdown() {
	log.Infof("Shutting down gRPC server on address '%s'.", s.Addr)
	s.grpcServer.GracefulStop()
}
