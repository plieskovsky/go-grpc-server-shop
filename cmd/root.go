package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/pkg/errors"
	"github.com/plieskovsky/go-grpc-server-shop/internal/config"
	"github.com/plieskovsky/go-grpc-server-shop/internal/repository"
	"github.com/plieskovsky/go-grpc-server-shop/internal/server"
	"github.com/plieskovsky/go-grpc-server-shop/internal/service"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

const serviceName = "go-grpc-server-shop"

var (
	cfgFile string
	cfg     config.Configuration
)

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "./config.yaml", "Path to the config file")
}

// Execute run root command (main entry-point).
func Execute() error {
	return rootCmd.Execute()
}

var rootCmd = &cobra.Command{
	Use:               serviceName,
	DisableAutoGenTag: true,
	Short:             "o-grpc-server-shop",
	Long:              "GO gRPC Server with simple shop like CRUD API",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		cfg = config.MustParse(cfgFile)
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		repo := repository.NewInMemoryRepo()
		mTLSCfg, err := createMTLSCfg(cfg)
		if err != nil {
			return err
		}

		grpcServer := createGrpcServer(cfg.Server.Grpc, mTLSCfg, repo)
		go func() {
			if err := grpcServer.ListenAndServe(); err != nil {
				log.Panicf("Failed to listen or serve: %v", err)
			}
		}()
		defer grpcServer.GracefulShutdown()

		// Block until we receive the signal.
		<-sigs
		return nil
	},
}

func createMTLSCfg(cfg config.Configuration) (*tls.Config, error) {
	srvTlsCfg, err := createSrvTlsCfg(cfg)
	if err != nil {
		return nil, err
	}

	cp, err := loadClientCACerts(cfg)
	if err != nil {
		return nil, err
	}

	tlsCfg := srvTlsCfg
	// this forces mTLS - client has to provide its certificate
	tlsCfg.ClientAuth = tls.RequireAndVerifyClientCert
	tlsCfg.ClientCAs = cp
	return tlsCfg, nil
}

func createSrvTlsCfg(cfg config.Configuration) (*tls.Config, error) {
	srvKP, err := tls.LoadX509KeyPair(cfg.Server.Grpc.CertFilename, cfg.Server.Grpc.KeyFilename)
	if err != nil {
		return nil, err
	}
	log.Infof("certificate and key loaded")
	return &tls.Config{Certificates: []tls.Certificate{srvKP}}, nil
}

func loadClientCACerts(cfg config.Configuration) (*x509.CertPool, error) {
	log.Infof("Loading CA certificate for gRPC from path '%v'.", cfg.Server.Grpc.ClientCACert)
	b, err := ioutil.ReadFile(cfg.Server.Grpc.ClientCACert)
	if err != nil {
		return nil, errors.Wrap(err, "client CA certificate read")
	}
	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(b) {
		return nil, errors.New("failed to append client CA certificate")
	}

	return cp, nil
}

func createGrpcServer(opts server.Config, tls *tls.Config, r *repository.InMemoryRepo) *server.ShopServer {
	server := server.New(opts, tls)

	service := service.ShopService{ItemsRepo: r}
	service.Register(server)

	return server
}
