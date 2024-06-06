package database

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"metrix/internal/exceptions"
	"metrix/internal/logger"
	"os"
	"strings"
	"sync"
	"time"
)

type MemoryStorage struct {
	mux           *sync.RWMutex
	s             map[string]float64
	saveSync      bool
	storeInterval int64
	filePath      string
}

func NewStorage(
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
		s:             make(map[string]float64),
		storeInterval: storeInterval,
		saveSync:      saveSync,
		filePath:      filePath,
	}

	if restore && ms.filePath != "" {
		err := ms.Restore()
		if err != nil {
			logger.Log.Warnf("failed to restore: db %s", err)
		}
	}

	if !saveSync {
		go ms.PeriodicBackup(ctx)
	}

	return ms
}

func (s *MemoryStorage) PeriodicBackup(ctx context.Context) {
	logger.Log.Info("starting backing up")
	ticker := time.NewTicker(time.Duration(int64(time.Second) * s.storeInterval))
	for {
		select {
		case <-ticker.C:
			err := s.BackUp()
			if err != nil {
				logger.Log.Warnf("failed to backup db %s", err)
			} else {
				logger.Log.Info("db backed up")
			}
		case <-ctx.Done():
			err := s.BackUp()
			if err != nil {
				logger.Log.Warnf("failed to backup db %s", err)
			} else {
				logger.Log.Info("db backed up")
			}
			ticker.Stop()
			return
		}
	}
}

func (s *MemoryStorage) BackUp() error {
	if s.filePath == "" {
		return nil
	}
	s.mux.RLock()
	defer s.mux.RUnlock()

	data, err := json.Marshal(s.s)
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

	logger.Log.Info("restoring db")

	s.mux.Lock()
	defer s.mux.Unlock()

	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	newStorage := make(map[string]float64)
	err = json.Unmarshal(data, &newStorage)
	if err != nil {
		return fmt.Errorf("failed to unmarshal storage: %w", err)
	}

	s.s = newStorage

	logger.Log.Info("restored db")

	return nil
}

func (s *MemoryStorage) Get(key string) (float64, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	res, ok := s.s[key]
	if !ok {
		return 0, exceptions.RecordNotFoundError{
			Msg: "not found value for key: " + key,
		}
	}

	return res, nil
}

func (s *MemoryStorage) Set(key string, value float64) (float64, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.s[key] = value

	fmt.Printf(">>>> %+v", s.s)

	if s.saveSync && s.filePath != "" {
		data, err := json.Marshal(s.s)
		if err != nil {
			return value, fmt.Errorf("failed to marshal storage: %w", err)
		}

		err = os.WriteFile(s.filePath, data, fs.ModePerm)
		if err != nil {
			return value, fmt.Errorf("failed to write file: %w", err)
		}
	}

	return value, nil
}

func (s *MemoryStorage) Keys(namespace string) ([]string, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	keys := make([]string, 0)

	for k := range s.s {
		if strings.HasPrefix(k, namespace) {
			keys = append(keys, k)
		}
	}

	return keys, nil
}
