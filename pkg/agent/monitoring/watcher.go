package monitoring

import (
	"context"
	"fmt"
	"metrix/pkg/agent/config"
	"metrix/pkg/logger"
	"strings"
	"time"
)

type Watcher struct {
	stats   *Stats
	metrics *[]string
	ch      chan struct{}
}

func NewWatcher(metrics []string) *Watcher {
	return &Watcher{
		stats:   NewStats(),
		metrics: &metrics,
		ch:      make(chan struct{}),
	}
}

func (w Watcher) watch(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			if err := w.stats.Read(ctx, *w.metrics...); err != nil {
				logger.Error(ctx, "failed to read metrics", err)
			} else {
				w.stats.IncrementPollCount()
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (w Watcher) report(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			w.ch <- struct{}{}
		case <-ctx.Done():
			ticker.Stop()
			close(w.ch)
			return
		}
	}
}

func (w Watcher) worker(ctx context.Context, id int, c *Client) {
	logger.Info(ctx, fmt.Sprintf("Worker №%d has been started", id))

	for {
		select {
		case <-w.ch:
			metrics, err := w.stats.AsMapOfMetrics(*w.metrics...)
			if err != nil {
				logger.Error(ctx, "failed to get metrics", err)
			}

			if len(metrics) == 0 {
				logger.Warn(ctx, "metrics are empty, nothing to send")
			}

			err = c.SendMetrics(ctx, metrics)
			if err != nil {
				logger.Error(ctx, "failed to send metrics", err)
			} else {
				w.stats.ResetPollCount()
				logger.Info(ctx, fmt.Sprintf("Worker №%d has sent metrics", id))
			}
		case <-ctx.Done():
			logger.Info(ctx, fmt.Sprintf("Worker №%d shutting down", id))
			return
		}
	}
}

func (w Watcher) Run(
	ctx context.Context,
	cfg *config.Config,
) {
	address := cfg.Address
	if !strings.HasPrefix(address, "http") {
		address = "http://" + address
	}

	client := NewClient(
		ctx,
		address,
		cfg.SignKey,
		cfg.UseBatching,
		int(cfg.RetryCount),
		cfg.RetryWaitTime,
		cfg.RetryMaxWaitTime,
	)

	for i := 1; i <= int(cfg.Goroutines); i++ {
		go w.worker(ctx, i, client)
	}

	go w.watch(ctx, time.Duration(cfg.PollInterval*int64(time.Second)))
	go w.report(ctx, time.Duration(cfg.ReportInterval*int64(time.Second)))

	<-ctx.Done()
}
