package service

import (
	"context"
	"fmt"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/plieskovsky/go-grpc-server-shop/internal/repository"
	"github.com/plieskovsky/go-grpc-server-shop/internal/server"
	"github.com/plieskovsky/go-grpc-server-shop/proto"
	log "github.com/sirupsen/logrus"
	"github.com/twinj/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ItemsRepo provides functions to manage Items in repository.
type ItemsRepo interface {
	Get(id string) (*proto.Item, error)
	GetAll() (*proto.ItemsList, error)
	Upsert(i *proto.Item) (*proto.Item, error)
	Remove(id string) error
}

// ShopService provides CRUD on Items.
type ShopService struct {
	proto.UnimplementedShopServiceServer
	ItemsRepo ItemsRepo
}

// Register registers the service to gRPC server.
func (s *ShopService) Register(server *server.ShopServer) {
	proto.RegisterShopServiceServer(server, s)
}

func (s *ShopService) GetAll(_ context.Context, _ *empty.Empty) (*proto.ItemsList, error) {
	log.Info("Get all items item request.")

	i, err := s.ItemsRepo.GetAll()
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (s *ShopService) Get(_ context.Context, id *proto.ItemRequestId) (*proto.Item, error) {
	log.Infof("Get item request '%+v'.", id)

	i, err := s.ItemsRepo.Get(id.GetId())
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (s *ShopService) Create(_ context.Context, item *proto.CreateItemRequest) (*proto.Item, error) {
	log.Infof("Create item request '%+v'.", item)

	uuid := uuid.NewV4().String()
	i := &proto.Item{
		Id:    uuid,
		Name:  item.Name,
		Price: item.Price,
	}
	i, err := s.ItemsRepo.Upsert(i)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (s *ShopService) Update(_ context.Context, i *proto.Item) (*proto.Item, error) {
	log.Infof("Update item request '%+v'.", i)

	_, err := s.ItemsRepo.Get(i.GetId())
	if errors.Is(err, repository.NotFoundErr) {
		return nil, status.Errorf(codes.InvalidArgument, fmt.Sprintf("Item with id '%s' doesn't exist.", i.GetId()))
	}

	i, err = s.ItemsRepo.Upsert(i)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (s *ShopService) Remove(_ context.Context, id *proto.ItemRequestId) (*empty.Empty, error) {
	log.Infof("Remove item request '%+v'.", id)

	err := s.ItemsRepo.Remove(id.GetId())
	if err != nil {
		return nil, err
	}

	return &empty.Empty{}, nil
}
