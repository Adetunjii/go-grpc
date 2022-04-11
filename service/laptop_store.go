package service

import (
	"errors"
	"fmt"
	"github.com/Adetunjii/go-grpc/pb"
	"github.com/jinzhu/copier"
	"sync"
)

var DuplicateException = errors.New("resource already exists")

type LaptopStore interface {
	Save(laptop *pb.Laptop) error
	FindById(laptopId string) (*pb.Laptop, error)
}

// store laptops in memory
type InMemoryLaptopStore struct {
	mutex sync.RWMutex // to handle concurrency while saving
	data  map[string]*pb.Laptop
}

func NewInMemoryLaptopStore() *InMemoryLaptopStore {
	return &InMemoryLaptopStore{
		data: make(map[string]*pb.Laptop),
	}
}

func (store *InMemoryLaptopStore) Save(laptop *pb.Laptop) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.data[laptop.Id] != nil {
		return DuplicateException
	}

	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return fmt.Errorf("cannot copy laptop data: %v", err)
	}

	store.data[other.Id] = other
	return nil
}

func (store *InMemoryLaptopStore) FindById(id string) (*pb.Laptop, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	laptop := store.data[id]
	if laptop == nil {
		return nil, nil
	}

	other := &pb.Laptop{}
	err := copier.Copy(other, laptop)
	if err != nil {
		return nil, fmt.Errorf("cannot copy laptop data: %v", err)
	}

	return other, nil
}
