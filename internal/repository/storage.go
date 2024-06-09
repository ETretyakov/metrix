package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"metrix/internal/model"
	"metrix/pkg/logger"
	"os"
	"sync"
	"time"
)

type MemoryStorage struct {
	mux           *sync.RWMutex
	storage       map[string]model.Metric
	saveSync      bool
	storeInterval int64
	filePath      string
}

func (s *MemoryStorage) Create(
	ctx context.Context,
	metric *model.Metric,
) (*model.Metric, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.storage[metric.ID] = *metric

	if s.saveSync {
		err := s.backUp()
		if err != nil {
			return nil, fmt.Errorf("failed to backup storage: %w", err)
		}
	}

	return metric, nil
}

func (s *MemoryStorage) Read(
	ctx context.Context,
	metricID string,
) (*model.Metric, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	m, ok := s.storage[metricID]
	if !ok {
		return nil, nil
	}

	return &m, nil
}

func (s *MemoryStorage) ReadIDs(
	ctx context.Context,
) (*[]string, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	var ids []string

	for k := range s.storage {
		ids = append(ids, k)
	}

	return &ids, nil
}

func (s *MemoryStorage) Update(
	ctx context.Context,
	metric *model.Metric,
) (*model.Metric, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.storage[metric.ID] = *metric

	if s.saveSync {
		err := s.backUp()
		if err != nil {
			return nil, fmt.Errorf("failed to backup storage: %w", err)
		}
	}

	return metric, nil
}

func (s *MemoryStorage) Delete(
	ctx context.Context,
	metricID string,
) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	delete(s.storage, metricID)

	if s.saveSync {
		err := s.backUp()
		if err != nil {
			return fmt.Errorf("failed to backup storage: %w", err)
		}
	}

	return nil
}

func NewInMemmoryStorage(
	ctx context.Context,
	filePath string,
	storeInterval int64,
	restore bool,
) *MemoryStorage {
	saveSync := false
	if storeInterval == 0 {
		saveSync = true
	}

	ms := &MemoryStorage{
		mux:           &sync.RWMutex{},
		storage:       make(map[string]model.Metric),
		storeInterval: storeInterval,
		saveSync:      saveSync,
		filePath:      filePath,
	}

	if restore {
		ms.Restore()
	}

	if !saveSync {
		go ms.PeriodicBackup(ctx)
	}

	return ms
}

func (s *MemoryStorage) PeriodicBackup(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(int64(time.Second) * s.storeInterval))
	for {
		select {
		case <-ticker.C:
			err := s.BackUp()
			if err != nil {
				logger.Warn(ctx, fmt.Sprintf("failed to backup db %s", err))
			}
		case <-ctx.Done():
			err := s.BackUp()
			if err != nil {
				logger.Warn(ctx, fmt.Sprintf("failed to backup db %s", err))
			}
			ticker.Stop()
			return
		}
	}
}

func (s *MemoryStorage) BackUp() error {
	s.mux.RLock()
	defer s.mux.RUnlock()

	err := s.backUp()
	if err != nil {
		return fmt.Errorf("failed to back up: %w", err)
	}

	return nil
}

func (s *MemoryStorage) backUp() error {
	if s.filePath == "" {
		return nil
	}

	data, err := json.Marshal(s.storage)
	if err != nil {
		return fmt.Errorf("failed to marshal storage: %w", err)
	}

	err = os.WriteFile(s.filePath, data, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (s *MemoryStorage) Restore() error {
	if s.filePath == "" {
		return nil
	}

	s.mux.Lock()
	defer s.mux.Unlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	newStorage := make(map[string]model.Metric)
	err = json.Unmarshal(data, &newStorage)
	if err != nil {
		return fmt.Errorf("failed to unmarshal storage: %w", err)
	}

	s.storage = newStorage

	return nil
}

func (r *MemoryStorage) PingDB() bool {
	return true
}
