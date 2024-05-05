package infrastructure

import (
	"fmt"
	"metrix/internal/database"
	"metrix/internal/interfaces"
	"strings"
)

type StorageHandler struct {
	Storage *database.MemoryStorage
}

func (s *StorageHandler) Key(sections ...string) string {
	return strings.Join(sections, ":")
}

func (s *StorageHandler) Get(key string) (float64, error) {
	val, err := s.Storage.Get(key)
	if err != nil {
		return 0, fmt.Errorf("failed to get value from memory storage: %w", err)
	}

	return val, nil
}

func (s *StorageHandler) Set(key string, value float64) (float64, error) {
	val, err := s.Storage.Set(key, value)
	if err != nil {
		return 0, fmt.Errorf("failed to set value from memory storage: %w", err)
	}

	return val, nil
}

func (s *StorageHandler) Keys(namespace string) ([]string, error) {
	val, err := s.Storage.Keys(namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get keys from memory storage: %w", err)
	}

	return val, nil
}

func NewStorageHandler() (interfaces.StorageHandler, error) {
	storageHandler := &StorageHandler{
		Storage: database.NewStorage(),
	}

	return storageHandler, nil
}
