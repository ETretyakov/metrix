// Module "storages" comprises functionality to build custom memory storages.
package storages

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"sync"
	"time"

	"metrix/internal/model"
	"metrix/pkg/logger"
)

// MemoryStorage - is the structure to manage inmemory storage.
type MemoryStorage struct {
	mux           *sync.RWMutex
	storage       map[string]model.Metric
	saveSync      bool
	storeInterval int64
	filePath      string
}

// Create - the method to add metric record to memory storage.
func (s *MemoryStorage) Create(
	ctx context.Context,
	metric *model.Metric,
) (*model.Metric, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.storage[metric.ID] = *metric

	if s.saveSync {
		err := s.writeToFile()
		if err != nil {
			return nil, fmt.Errorf("failed to backup storage: %w", err)
		}
	}

	return metric, nil
}

// Read - the method to read a metric record from memory storage.
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

// ReadIDs - the method to read stored metrics IDs.
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

// Update - the method to update stored metrics.
func (s *MemoryStorage) Update(
	ctx context.Context,
	metric *model.Metric,
) (*model.Metric, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.storage[metric.ID] = *metric

	if s.saveSync {
		err := s.writeToFile()
		if err != nil {
			return nil, fmt.Errorf("failed to backup storage: %w", err)
		}
	}

	return metric, nil
}

// Delete - the method to delete stored metrics.
func (s *MemoryStorage) Delete(
	ctx context.Context,
	metricID string,
) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	delete(s.storage, metricID)

	if s.saveSync {
		err := s.writeToFile()
		if err != nil {
			return fmt.Errorf("failed to backup storage: %w", err)
		}
	}

	return nil
}

// NewInMemmoryStorage - the building function for InMemoryStorage.
func NewInMemoryStorage(
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
		ms.restore()
	}

	if !saveSync {
		go ms.PeriodicBackup(ctx)
	}

	return ms
}

// PeriodicBackup - the method that initiates periodic inmemory backups.
func (s *MemoryStorage) PeriodicBackup(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(int64(time.Second) * s.storeInterval))
	for {
		select {
		case <-ticker.C:
			err := s.muxWriteToFile()
			if err != nil {
				logger.Warn(ctx, fmt.Sprintf("failed to backup db %s", err))
			}
		case <-ctx.Done():
			err := s.muxWriteToFile()
			if err != nil {
				logger.Warn(ctx, fmt.Sprintf("failed to backup db %s", err))
			}
			ticker.Stop()
			return
		}
	}
}

func (s *MemoryStorage) muxWriteToFile() error {
	s.mux.RLock()
	defer s.mux.RUnlock()

	err := s.writeToFile()
	if err != nil {
		return fmt.Errorf("failed to back up: %w", err)
	}

	return nil
}

func (s *MemoryStorage) writeToFile() error {
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

func (s *MemoryStorage) restore() error {
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

// ReadMany - the method to read many metrics records.
func (s *MemoryStorage) ReadMany(ctx context.Context, metricIDs []string) (*[]model.Metric, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	metrics := []model.Metric{}
	for _, id := range metricIDs {
		metric, ok := s.storage[id]
		if ok {
			metrics = append(metrics, metric)
		}
	}

	return &metrics, nil
}

// UpsertMany - the method to insert/read many metric records.
func (s *MemoryStorage) UpsertMany(
	ctx context.Context,
	metrics []model.Metric,
) (bool, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	for _, m := range metrics {
		s.storage[m.ID] = m
	}

	if s.saveSync {
		err := s.writeToFile()
		if err != nil {
			return false, fmt.Errorf("failed to backup storage: %w", err)
		}
	}

	return true, nil
}

// PingDB - the method to ping inmemory storage.
func (s *MemoryStorage) PingDB(ctx context.Context) bool {
	return true
}
