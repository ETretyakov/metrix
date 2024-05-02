package watcher

import (
	"context"
	"fmt"
	"math/rand"
	"metrix/pkg/agent/config"
	"metrix/pkg/client"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type Watcher struct {
	mux         *sync.RWMutex
	Metrics     map[string]float64
	PollCount   float64
	RandomValue float64
}

func NewWatcher(config config.Config) *Watcher {
	metrics := map[string]float64{}

	for _, v := range config.Metrics {
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
			w.PollCount += 1
			w.RandomValue = rand.Float64()

			runtime.ReadMemStats(&stats)

			reflect.ValueOf(stats)
			for k := range w.Metrics {
				reflectValue := reflect.ValueOf(stats)
				field := reflect.Indirect(reflectValue).FieldByName(k)
				switch field.Kind() {
				case reflect.Uint32:
					w.Metrics[k] = float64(field.Interface().(uint32))
				case reflect.Uint64:
					w.Metrics[k] = float64(field.Interface().(uint64))
				case reflect.Float64:
					w.Metrics[k] = field.Interface().(float64)
				default:
					log.Warn().Caller().Str("Stage", "reading-metrics").
						Msg(fmt.Sprintf("unsupported metric field type: %s", field.Kind()))
				}
			}
			w.mux.Unlock()
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (w *Watcher) Report(ctx context.Context, baseUrl string, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			w.mux.Lock()
			log.Info().Caller().Str("Stage", "reporting-metrics").
				Msg(fmt.Sprintf("sending metrics: %+v", w))

			for k, v := range w.Metrics {
				err := client.SendMetric(ctx, baseUrl, client.GaugeType, k, v)
				if err != nil {
					log.Error().Err(err).Str("Stage", "reporting-metrics").
						Msg(fmt.Sprintf("failed to report metrics: %s", err))
				}
			}

			err := client.SendMetric(ctx, baseUrl, client.CounterType, "PollCount", w.PollCount)
			if err != nil {
				log.Error().Err(err).Str("Stage", "reporting-metrics").
					Msg(fmt.Sprintf("failed to report PollCount metric: %s", err))
			}

			err = client.SendMetric(ctx, baseUrl, client.CounterType, "RandomValue", w.RandomValue)
			if err != nil {
				log.Error().Err(err).Str("Stage", "reporting-metrics").
					Msg(fmt.Sprintf("failed to report RandomValue metric: %s", err))
			}

			w.PollCount = 0
			w.mux.Unlock()
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}
