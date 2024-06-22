package database

import (
	"metrix/internal/exceptions"
	"strings"
	"sync"
)

type MemoryStorage struct {
	mux *sync.RWMutex
	s   map[string]float64
}

func NewStorage() *MemoryStorage {
	return &MemoryStorage{
		mux: &sync.RWMutex{},
		s:   make(map[string]float64),
	}
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

	return value, nil
}

func (s *MemoryStorage) Keys(namespace string) ([]string, error) {
	keys := make([]string, 0)

	for k := range s.s {
		if strings.HasPrefix(k, namespace) {
			keys = append(keys, k)
		}
	}

	return keys, nil
}
