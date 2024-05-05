package database

import (
	"fmt"
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

func (s *MemoryStorage) Get(key string) (*float64, error) {
	s.mux.RLock()
	defer s.mux.RUnlock()

	res, ok := s.s[key]
	if !ok {
		return nil, nil
	}

	return &res, nil
}

func (s *MemoryStorage) Set(key string, value float64) (*float64, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.s[key] = value

	return &value, nil
}

func (s *MemoryStorage) Keys(namespace string) ([]string, error) {
	keys := make([]string, len(s.s))

	for k := range s.s {
		if strings.HasPrefix(k, namespace) {
			splitKey := strings.Split(k, ":")
			keys = append(keys, fmt.Sprintf("%s %s", splitKey[1], splitKey[2]))
		}
	}

	return keys, nil
}
