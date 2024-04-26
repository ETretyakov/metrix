package db

import (
	"fmt"

	"github.com/rs/zerolog/log"
)

type Storage struct {
	Data map[string]uint64
}

func (m *Storage) Key(metricType string, name string) string {
	return metricType + ":" + name
}

func (m *Storage) Set(metricType string, name string, value uint64) {
	log.Info().Msg(
		fmt.Sprintf(
			"[memstorage:Set] setting value for [%s] %s: %d -> %d",
			metricType,
			name,
			m.Get(metricType, name),
			value,
		),
	)

	m.Data[m.Key(metricType, name)] = value
}

func (m *Storage) Get(metricType string, name string) uint64 {
	return m.Data[m.Key(metricType, name)]
}

var MemStorage = Storage{
	Data: make(map[string]uint64),
}
