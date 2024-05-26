package infrastructure

import (
	"context"
	"fmt"
	"metrix/internal/database"
	"metrix/internal/interfaces"
	"strings"
	"time"
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

func NewStorageHandler(
	ctx context.Context,
	filePath string,
	storeInterval time.Duration,
	restore bool,
) (interfaces.StorageHandler, error) {
	storageHandler := &StorageHandler{
		Storage: database.NewStorage(ctx, filePath, storeInterval, restore),
	}

	return storageHandler, nil
}

func (s *StorageHandler) BackUp() error {
	err := s.Storage.BackUp()
	if err != nil {
		return fmt.Errorf("failed to backup memory storage: %w", err)
	}

	return nil
}

func (s *StorageHandler) Restore() error {
	err := s.Storage.Restore()
	if err != nil {
		return fmt.Errorf("failed to restore memory storage: %w", err)
	}

	return nil
}
