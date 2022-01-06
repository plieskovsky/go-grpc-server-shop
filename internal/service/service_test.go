package service

import (
	"context"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/pkg/errors"
	"github.com/plieskovsky/go-grpc-server-shop/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"reflect"
	"testing"
)

type repoMock struct {
	mock.Mock
}

func (m *repoMock) Get(id string) (*proto.Item, error) {
	args := m.Called(id)
	return args.Get(0).(*proto.Item), args.Error(1)
}

func (m *repoMock) GetAll() (*proto.ItemsList, error) {
	args := m.Called()
	return args.Get(0).(*proto.ItemsList), args.Error(1)
}

func (m *repoMock) Upsert(i *proto.Item) (*proto.Item, error) {
	args := m.Called(i)
	return args.Get(0).(*proto.Item), args.Error(1)
}

func (m *repoMock) Remove(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

var i1 = proto.Item{
	Id:    "id-1",
	Name:  "name-1",
	Price: 1.1,
}
var i2 = proto.Item{
	Id:    "id-2",
	Name:  "name-2",
	Price: 2.2,
}

func TestShopService_GetAll(t *testing.T) {
	type repoResponse struct {
		items *proto.ItemsList
		err   error
	}
	tests := []struct {
		name         string
		repoResponse repoResponse
		want         *proto.ItemsList
		wantErr      bool
	}{
		{
			name: "repository returns items",
			repoResponse: repoResponse{
				items: &proto.ItemsList{Items: []*proto.Item{&i1, &i2}},
				err:   nil,
			},
			want: &proto.ItemsList{Items: []*proto.Item{&i1, &i2}},
		},
		{
			name: "repository returns error",
			repoResponse: repoResponse{
				err: errors.New("repo error"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := new(repoMock)
			r.On("GetAll").Return(tt.repoResponse.items, tt.repoResponse.err)
			s := &ShopService{ItemsRepo: r}

			got, err := s.GetAll(context.Background(), &empty.Empty{})

			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.True(t, reflect.DeepEqual(tt.want.Items, got.Items))
			}
		})
	}
}
