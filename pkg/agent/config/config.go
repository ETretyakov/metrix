package config

import (
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/caarlos0/env"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	ServerUrl      string        `env:"AGT_SERVER_URL"      mapstructure:"AGT_SERVER_URL"      envDefault:"http://localhost:8080"`
	PollInterval   time.Duration `env:"AGT_POLL_INTERVAL"   mapstructure:"AGT_POLL_INTERVAL"   envDefault:"2s"`
	ReportInterval time.Duration `env:"AGT_REPORT_INTERVAL" mapstructure:"AGT_REPORT_INTERVAL" envDefault:"10s"`
	Metrics        []string      `env:"AGT_METRICS"         mapstructure:"AGT_METRICS"         envDefault:"*"`
}

func LoadConfig() (*Config, error) {
	log.Logger = log.Output(
		zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		},
	)

	log.Info().Caller().Str("Stage", "loading-config").
		Msg("started config loading")
	config := &Config{}
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("env.Parse: %w", err)
	}

	log.Info().Caller().Str("Stage", "loading-config").
		Msg(fmt.Sprintf("parsed config: %+v", config))

	if len(config.Metrics) == 0 {
		return nil, fmt.Errorf("metrics were not provided")
	}

	if slices.Contains(config.Metrics, "*") {
		log.Info().Caller().Str("Stage", "loading-config").
			Msg("detected * in metrics - loading all metrics")

		config.Metrics = []string{
			"Alloc",
			"BuckHashSys",
			"Frees",
			"GCCPUFraction",
			"GCSys",
			"HeapAlloc",
			"HeapIdle",
			"HeapInuse",
			"HeapObjects",
			"HeapReleased",
			"HeapSys",
			"LastGC",
			"Lookups",
			"MCacheInuse",
			"MSpanSys",
			"Mallocs",
			"NextGC",
			"NumForcedGC",
			"NumGC",
			"OtherSys",
			"PauseTotalNs",
			"StackInuse",
			"StackSys",
			"Sys",
			"TotalAlloc",
		}
	}

	log.Info().Caller().Str("Stage", "loading-config").
		Msg(fmt.Sprintf("config: %+v", config))

	return config, nil
}
