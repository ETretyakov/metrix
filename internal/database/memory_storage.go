package database

import (
	"sync"
)

type MemoryStorage struct {
	mux *sync.Mutex
	s   map[string]uint64
}

func NewStorage() *MemoryStorage {
	return &MemoryStorage{
		mux: &sync.Mutex{},
		s:   make(map[string]uint64),
	}
}

func (s *MemoryStorage) Get(key string) (*uint64, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	res, ok := s.s[key]
	if !ok {
		return nil, nil
	}

	return &res, nil
}

func (s *MemoryStorage) Set(key string, value uint64) (*uint64, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.s[key] = value

	return &value, nil
}
