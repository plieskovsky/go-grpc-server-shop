package repository

import (
	"github.com/plieskovsky/go-grpc-server-shop/proto"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

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

func TestInMemoryRepo_Get(t *testing.T) {
	tests := []struct {
		name    string
		items   items
		id      string
		want    *proto.Item
		wantErr bool
	}{
		{
			name:  "get present item",
			items: map[string]*proto.Item{"id-1": &i1, "id-2": &i2},
			id:    "id-1",
			want:  &i1,
		},
		{
			name:    "get non-present item",
			items:   map[string]*proto.Item{"id-1": &i1, "id-2": &i2},
			id:      "id-3",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := InMemoryRepo{
				items: tt.items,
			}

			got, err := r.Get(tt.id)

			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestInMemoryRepo_GetAll(t *testing.T) {
	tests := []struct {
		name    string
		items   items
		want    *proto.ItemsList
		wantErr bool
	}{
		{
			name:  "items present in repo",
			items: map[string]*proto.Item{"id-1": &i1, "id-2": &i2},
			want:  &proto.ItemsList{Items: []*proto.Item{&i1, &i2}},
		},
		{
			name:  "items empty",
			items: map[string]*proto.Item{},
			want:  &proto.ItemsList{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := InMemoryRepo{
				items: tt.items,
			}

			got, err := r.GetAll()

			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.True(t, reflect.DeepEqual(tt.want.Items, got.Items))
			}
		})
	}
}

func TestInMemoryRepo_Upsert(t *testing.T) {
	tests := []struct {
		name    string
		items   items
		i       *proto.Item
		want    *proto.Item
		wantErr bool
	}{
		{
			name:  "upsert non existing",
			items: map[string]*proto.Item{"id-1": &i1},
			i:     &i2,
			want:  &i2,
		},
		{
			name:  "upsert existing",
			items: map[string]*proto.Item{"id-1": &i1},
			i: &proto.Item{
				Id:    "id-1",
				Name:  "updated-1",
				Price: 5.5,
			},
			want: &proto.Item{
				Id:    "id-1",
				Name:  "updated-1",
				Price: 5.5,
			},
		},
		{
			name:  "upsert item when items are empty",
			items: map[string]*proto.Item{},
			i:     &i1,
			want:  &i1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := InMemoryRepo{
				items: tt.items,
			}

			got, err := r.Upsert(tt.i)

			assert.Equal(t, tt.wantErr, err != nil)
			if !tt.wantErr {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestInMemoryRepo_Remove(t *testing.T) {
	tests := []struct {
		name    string
		items   items
		id      string
		want    items
		wantErr bool
	}{
		{
			name:  "Remove existing",
			items: map[string]*proto.Item{"id-1": &i1, "id-2": &i2},
			want:  map[string]*proto.Item{"id-1": &i1},
			id:    "id-2",
		},
		{
			name:  "Remove non existing",
			items: map[string]*proto.Item{"id-1": &i1, "id-2": &i2},
			want:  map[string]*proto.Item{"id-1": &i1, "id-2": &i2},
			id:    "id-5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := InMemoryRepo{
				items: tt.items,
			}

			err := r.Remove(tt.id)

			assert.Equal(t, tt.wantErr, err != nil)
			assert.Equal(t, tt.want, r.items)
		})
	}
}
