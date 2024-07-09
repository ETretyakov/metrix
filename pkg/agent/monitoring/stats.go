package monitoring

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"metrix/pkg/logger"
	"reflect"
	"runtime"
	"sync"
)

type Stats struct {
	Alloc           float64 `metricType:"gauge"`
	BuckHashSys     float64 `metricType:"gauge"`
	Frees           float64 `metricType:"gauge"`
	GCCPUFraction   float64 `metricType:"gauge"`
	GCSys           float64 `metricType:"gauge"`
	HeapAlloc       float64 `metricType:"gauge"`
	HeapIdle        float64 `metricType:"gauge"`
	HeapInuse       float64 `metricType:"gauge"`
	HeapObjects     float64 `metricType:"gauge"`
	HeapReleased    float64 `metricType:"gauge"`
	HeapSys         float64 `metricType:"gauge"`
	LastGC          float64 `metricType:"gauge"`
	Lookups         float64 `metricType:"gauge"`
	MCacheInuse     float64 `metricType:"gauge"`
	MCacheSys       float64 `metricType:"gauge"`
	MSpanInuse      float64 `metricType:"gauge"`
	MSpanSys        float64 `metricType:"gauge"`
	Mallocs         float64 `metricType:"gauge"`
	NextGC          float64 `metricType:"gauge"`
	NumForcedGC     float64 `metricType:"gauge"`
	NumGC           float64 `metricType:"gauge"`
	OtherSys        float64 `metricType:"gauge"`
	PauseTotalNs    float64 `metricType:"gauge"`
	StackInuse      float64 `metricType:"gauge"`
	StackSys        float64 `metricType:"gauge"`
	Sys             float64 `metricType:"gauge"`
	TotalAlloc      float64 `metricType:"gauge"`
	TotalMemmory    float64 `metricType:"gauge"`
	FreeMemory      float64 `metricType:"gauge"`
	CPUutilization1 float64 `metricType:"gauge"`
	RandomValue     float64 `metricType:"gauge"`
	PollCount       int64   `metricType:"counter"`
	mux             *sync.RWMutex
}

type Metric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

func NewStats() *Stats {
	return &Stats{
		mux: &sync.RWMutex{},
	}
}

func (rs *Stats) Read(ctx context.Context, metrics ...string) error {
	if rs == nil {
		return errors.New("failed to read metrics for nil pointer")
	}

	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	rs.mux.Lock()
	defer rs.mux.Unlock()

	for _, metric := range metrics {
		field := reflect.
			Indirect(reflect.ValueOf(stats)).
			FieldByName(metric)

		switch field.Kind() {
		case reflect.Uint32:
			if val, ok := field.Interface().(uint32); !ok {
				logger.Warn(
					ctx,
					"failed to assert filetype Uint32",
				)
			} else {
				reflect.
					ValueOf(rs).
					Elem().
					FieldByName(metric).
					SetFloat(float64(val))
			}
		case reflect.Uint64:
			if val, ok := field.Interface().(uint64); !ok {
				logger.Warn(
					ctx,
					"failed to assert filetype Uint64",
				)
			} else {
				reflect.
					ValueOf(rs).
					Elem().
					FieldByName(metric).
					SetFloat(float64(val))
			}
		case reflect.Float64:
			if val, ok := field.Interface().(float64); !ok {
				logger.Warn(
					ctx,
					"failed to assert filetype Float64",
				)
			} else {
				reflect.
					ValueOf(rs).
					Elem().
					FieldByName(metric).
					SetFloat(float64(val))
			}
		default:
			logger.Warn(
				ctx,
				fmt.Sprintf(
					"unsupported metric field type: %s",
					field.Kind(),
				),
			)
		}
	}

	rs.RandomValue = rand.Float64()

	return nil
}

func (rs *Stats) IncrementPollCount() {
	rs.mux.Lock()
	rs.PollCount++
	rs.mux.Unlock()
}

func (rs *Stats) ResetPollCount() {
	rs.mux.Lock()
	rs.PollCount = 0
	rs.mux.Unlock()
}

func (rs *Stats) AsMapOfMetrics(metrics ...string) ([]*Metric, error) {
	if rs == nil {
		return nil, errors.New("failed to read metrics for nil pointer")
	}

	m := []*Metric{}

	rs.mux.RLock()
	defer rs.mux.RUnlock()

	v := reflect.ValueOf(rs).Elem()
	for _, metric := range metrics {
		if metricTypeField, ok := v.Type().FieldByName(metric); ok {
			metricType := metricTypeField.Tag.Get("metricType")

			switch metricType {
			case "gauge":
				metricValField := reflect.Indirect(v).FieldByName(metric)
				metricVal := metricValField.Float()
				m = append(
					m,
					&Metric{
						ID:    metric,
						MType: metricType,
						Value: &metricVal,
					},
				)
			case "counter":
				metricValField := reflect.Indirect(v).FieldByName(metric)
				metricVal := metricValField.Int()
				m = append(
					m,
					&Metric{
						ID:    metric,
						MType: metricType,
						Delta: &metricVal,
					},
				)
			}
		}
	}

	return m, nil
}
