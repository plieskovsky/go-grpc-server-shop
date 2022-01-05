package repository

import (
	"github.com/pkg/errors"
	"github.com/plieskovsky/go-grpc-server-shop/proto"
	"sync"
)

var NotFoundErr = errors.New("Item not found")

// InMemoryRepo represents repository of items protected by RW lock.
type InMemoryRepo struct {
	items items
	lock  sync.RWMutex
}

// items maps item ID to item.
type items map[string]*proto.Item

// NewInMemoryRepo creates a new empty repository that holds items in app memory.
func NewInMemoryRepo() *InMemoryRepo {
	items := make(map[string]*proto.Item)
	return &InMemoryRepo{items: items}
}

func (r *InMemoryRepo) Get(id string) (*proto.Item, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	z, ok := r.items[id]
	if !ok {
		return nil, NotFoundErr
	}

	return z, nil
}

func (r *InMemoryRepo) GetAll() (*proto.ItemsList, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	var items []*proto.Item
	for _, v := range r.items {
		items = append(items, v)
	}

	return &proto.ItemsList{Items: items}, nil
}

func (r *InMemoryRepo) Upsert(i *proto.Item) (*proto.Item, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.items[i.GetId()] = i
	return i, nil
}

func (r *InMemoryRepo) Remove(id string) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	delete(r.items, id)
	return nil
}
