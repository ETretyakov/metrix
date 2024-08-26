package monitoring

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"sync"

	"metrix/pkg/logger"

	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/v4/cpu"
)

// Stats - the function for all collactible stats.
type Stats struct {
	Alloc          float64   `metricType:"gauge"   metricGroup:"runtime"`
	BuckHashSys    float64   `metricType:"gauge"   metricGroup:"runtime"`
	Frees          float64   `metricType:"gauge"   metricGroup:"runtime"`
	GCCPUFraction  float64   `metricType:"gauge"   metricGroup:"runtime"`
	GCSys          float64   `metricType:"gauge"   metricGroup:"runtime"`
	HeapAlloc      float64   `metricType:"gauge"   metricGroup:"runtime"`
	HeapIdle       float64   `metricType:"gauge"   metricGroup:"runtime"`
	HeapInuse      float64   `metricType:"gauge"   metricGroup:"runtime"`
	HeapObjects    float64   `metricType:"gauge"   metricGroup:"runtime"`
	HeapReleased   float64   `metricType:"gauge"   metricGroup:"runtime"`
	HeapSys        float64   `metricType:"gauge"   metricGroup:"runtime"`
	LastGC         float64   `metricType:"gauge"   metricGroup:"runtime"`
	Lookups        float64   `metricType:"gauge"   metricGroup:"runtime"`
	MCacheInuse    float64   `metricType:"gauge"   metricGroup:"runtime"`
	MCacheSys      float64   `metricType:"gauge"   metricGroup:"runtime"`
	MSpanInuse     float64   `metricType:"gauge"   metricGroup:"runtime"`
	MSpanSys       float64   `metricType:"gauge"   metricGroup:"runtime"`
	Mallocs        float64   `metricType:"gauge"   metricGroup:"runtime"`
	NextGC         float64   `metricType:"gauge"   metricGroup:"runtime"`
	NumForcedGC    float64   `metricType:"gauge"   metricGroup:"runtime"`
	NumGC          float64   `metricType:"gauge"   metricGroup:"runtime"`
	OtherSys       float64   `metricType:"gauge"   metricGroup:"runtime"`
	PauseTotalNs   float64   `metricType:"gauge"   metricGroup:"runtime"`
	StackInuse     float64   `metricType:"gauge"   metricGroup:"runtime"`
	StackSys       float64   `metricType:"gauge"   metricGroup:"runtime"`
	Sys            float64   `metricType:"gauge"   metricGroup:"runtime"`
	TotalAlloc     float64   `metricType:"gauge"   metricGroup:"runtime"`
	TotalMemory    float64   `metricType:"gauge"   metricGroup:"gopsutil/mem" metricAlias:"Total"`
	FreeMemory     float64   `metricType:"gauge"   metricGroup:"gopsutil/mem" metricAlias:"Free"`
	CPUutilization []float64 `metricType:"gauge"   metricGroup:"gopsutil/cpu" metricAlias:"System"`
	RandomValue    float64   `metricType:"gauge"   metricGroup:"custom"`
	PollCount      int64     `metricType:"counter" metricGroup:"custom"`
	mux            *sync.RWMutex
}

// Metric - the structure for Metric validation.
type Metric struct {
	ID    string   `json:"id"`
	MType string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

// NewStats - the builder function for Stats.
func NewStats() *Stats {
	return &Stats{
		mux:            &sync.RWMutex{},
		CPUutilization: []float64{},
	}
}

// Read - the method for retrieving stats from runtime.
func (rs *Stats) Read(ctx context.Context, metrics ...string) error {
	if rs == nil {
		return errors.New("failed to read metrics for nil pointer")
	}

	var stats runtime.MemStats
	runtime.ReadMemStats(&stats)

	memoryStat, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("failed to retrieve virtual memory: %w", err)
	}

	timesStat, err := cpu.Times(true)
	if err != nil {
		return fmt.Errorf("failed to retrieve cpu time: %w", err)
	}

	rs.mux.Lock()
	defer rs.mux.Unlock()

	v := reflect.ValueOf(rs).Elem()
	for _, metric := range metrics {
		metricTypeField, ok := v.Type().FieldByName(metric)
		if !ok {
			logger.Warn(ctx, "unrecognised metric field: "+metric)
			continue
		}

		fields := []reflect.Value{}

		metricType := metricTypeField.Tag.Get("metricGroup")
		metricAlias := metricTypeField.Tag.Get("metricAlias")
		if metricAlias == "" {
			metricAlias = metric
		}

		switch metricType {
		case "runtime":
			fields = append(
				fields,
				reflect.
					Indirect(reflect.ValueOf(stats)).
					FieldByName(metricAlias),
			)
		case "gopsutil/mem":
			fields = append(
				fields,
				reflect.
					Indirect(reflect.ValueOf(memoryStat)).
					FieldByName(metricAlias),
			)
		case "gopsutil/cpu":
			for _, cpu := range timesStat {
				fields = append(
					fields,
					reflect.
						Indirect(reflect.ValueOf(cpu)).
						FieldByName(metricAlias),
				)
			}
		default:
			continue
		}

		if len(fields) == 1 {
			value := fieldToFloat64(ctx, fields[0])
			reflect.
				ValueOf(rs).
				Elem().
				FieldByName(metric).
				SetFloat(value)
		} else {
			slice := reflect.MakeSlice(
				reflect.TypeOf([]float64{}),
				len(fields),
				len(fields),
			)

			reflect.
				ValueOf(rs).
				Elem().
				FieldByName(metric).Set(slice)

			for i, field := range fields {
				value := fieldToFloat64(ctx, field)
				reflect.
					ValueOf(rs).
					Elem().
					FieldByName(metric).
					Index(i).
					SetFloat(value)
			}
		}
	}

	rs.RandomValue = rand.Float64()

	return nil
}

// IncrementPollCount - the method to increment poll count.
func (rs *Stats) IncrementPollCount() {
	rs.mux.Lock()
	rs.PollCount++
	rs.mux.Unlock()
}

// ResetPollCount - the method to reset poll count.
func (rs *Stats) ResetPollCount() {
	rs.mux.Lock()
	rs.PollCount = 0
	rs.mux.Unlock()
}

// AsMapOfMetrics - the method to convert metrics into metric array.
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
				value := reflect.Indirect(v).FieldByName(metric)
				slice, ok := value.Interface().([]float64)
				if !ok {
					metricVal := value.Float()
					m = append(
						m,
						&Metric{
							ID:    metric,
							MType: metricType,
							Value: &metricVal,
						},
					)
				} else {
					for i, value := range slice {
						m = append(
							m,
							&Metric{
								ID:    fmt.Sprintf("%s%d", metric, i+1),
								MType: metricType,
								Value: &value,
							},
						)
					}
				}
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
