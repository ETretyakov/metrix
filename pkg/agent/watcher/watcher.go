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

func (w *Watcher) Report(ctx context.Context, baseURL string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			metricsVal := map[string]float64{}
			w.mux.RLock()
			for k, v := range w.Metrics {
				metricsVal[k] = v
			}
			w.mux.RUnlock()

			logger.Info(
				ctx,
				fmt.Sprintf("sending metrics: %+v", metricsVal),
				stageLogKey, stageLogReportVal,
			)

			for k, v := range metricsVal {
				err := client.SendMetric(ctx, baseURL, client.GaugeType, k, v)
				if err != nil {
					logger.Error(
						ctx,
						fmt.Sprintf("failed to report metrics: %s", err),
						err,
						stageLogKey, stageLogReportVal,
					)
				}
			}

			err := client.SendMetric(ctx, baseURL, client.CounterType, "PollCount", w.PollCount)
			if err != nil {
				logger.Error(
					ctx,
					fmt.Sprintf("failed to report PollCount metric: %s", err),
					err,
					stageLogKey, stageLogReportVal,
				)
			}

			err = client.SendMetric(ctx, baseURL, client.GaugeType, "RandomValue", w.RandomValue)
			if err != nil {
				logger.Error(
					ctx,
					fmt.Sprintf("failed to report RandomValue metric: %s", err),
					err,
					stageLogKey, stageLogReportVal,
				)
			}

			w.PollCount = 0
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
