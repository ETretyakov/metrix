package watcher

import (
	"context"
	"fmt"
	"math/rand"
	"metrix/pkg/agent/config"
	"metrix/pkg/client"
	"metrix/pkg/logger"
	"reflect"
	"runtime"
	"sync"
	"time"
)

const (
	stageLogKey       = "stage"
	stageLogWatchVal  = "reading-metrics"
	stageLogReportVal = "reporting-metrics"
)

type Watcher struct {
	mux         *sync.RWMutex
	Metrics     map[string]float64
	PollCount   float64
	RandomValue float64
}

func NewWatcher(cfg config.Config) *Watcher {
	metrics := map[string]float64{}

	for _, v := range cfg.Metrics {
		metrics[v] = 0
	}

	return &Watcher{
		mux:         &sync.RWMutex{},
		PollCount:   0,
		RandomValue: 0,
		Metrics:     metrics,
	}
}

func (w *Watcher) Start(ctx context.Context, interval time.Duration) {
	var stats runtime.MemStats

	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			w.mux.Lock()
			w.PollCount++
			w.RandomValue = rand.Float64()

			runtime.ReadMemStats(&stats)

			reflect.ValueOf(stats)
			for k := range w.Metrics {
				reflectValue := reflect.ValueOf(stats)
				field := reflect.Indirect(reflectValue).FieldByName(k)
				switch field.Kind() {
				case reflect.Uint32:
					val, ok := field.Interface().(uint32)
					if !ok {
						logger.Warn(
							ctx,
							"failed to assert filetype uint32",
							stageLogKey, stageLogWatchVal,
						)
					}
					w.Metrics[k] = float64(val)
				case reflect.Uint64:
					val, ok := field.Interface().(uint64)
					if !ok {
						logger.Warn(
							ctx,
							"failed to assert filetype uint64",
							stageLogKey, stageLogWatchVal,
						)
					}
					w.Metrics[k] = float64(val)
				case reflect.Float64:
					val, ok := field.Interface().(float64)
					if !ok {
						logger.Warn(
							ctx,
							"failed to assert filetype float64",
							stageLogKey, stageLogWatchVal,
						)
					}
					w.Metrics[k] = val
				default:
					logger.Warn(
						ctx,
						fmt.Sprintf("unsupported metric field type: %s", field.Kind()),
						stageLogKey, stageLogWatchVal,
					)
				}
			}
			w.mux.Unlock()
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (w *Watcher) Report(
	ctx context.Context,
	baseURL string,
	interval time.Duration,
	useBatching bool,
) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			var err error

			metrics := []*client.Metrics{}

			w.mux.RLock()
			for k, v := range w.Metrics {
				metrics = append(
					metrics,
					&client.Metrics{
						ID:    k,
						MType: "gauge",
						Value: func() *float64 { i := float64(v); return &i }(),
					},
				)
			}

			metrics = append(
				metrics,
				&client.Metrics{
					ID:    "PollCount",
					MType: "counter",
					Delta: func() *int64 { i := int64(w.PollCount); return &i }(),
				},
			)

			metrics = append(
				metrics,
				&client.Metrics{
					ID:    "RandomValue",
					MType: "gauge",
					Value: func() *float64 { i := float64(w.RandomValue); return &i }(),
				},
			)
			w.mux.RUnlock()

			if useBatching {
				err = client.SendMetricBatch(ctx, baseURL, metrics)
				if err != nil {
					logger.Error(
						ctx,
						fmt.Sprintf("failed to report metrics as batch: %s", err),
						err,
						stageLogKey, stageLogReportVal,
					)
				}
			} else {
				err = client.SendMetric(ctx, baseURL, metrics)
				if err != nil {
					logger.Error(
						ctx,
						fmt.Sprintf("failed to report metrics: %s", err),
						err,
						stageLogKey, stageLogReportVal,
					)
				}
			}

			if err == nil {
				w.PollCount = 0
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
