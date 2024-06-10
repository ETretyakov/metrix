package config

import (
	"errors"
	"fmt"
	"slices"

	"github.com/caarlos0/env"
)

type Config struct {
	Address        string   `env:"ADDRESS"         envDefault:"localhost:8080" flag:"address"         flagShort:"a" flagDescription:"http adress"`
	PollInterval   int64    `env:"POLL_INTERVAL"   envDefault:"2"              flag:"poll_interval"   flagShort:"p" flagDescription:"interval between polling"`
	ReportInterval int64    `env:"REPORT_INTERVAL" envDefault:"10"             flag:"report_interval" flagShort:"r" flagDescription:"interval between reporting"`
	LogLevel       string   `env:"LOG_LEVEL"       envDefault:"info"           flag:"log_level"       flagShort:"l" flagDescription:"level for logging"`
	Metrics        []string `env:"AGT_METRICS"     envDefault:"*"`
	UseBatching    bool     `env:"USE_BATCHING"    envDefault:"true"`
}

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse agent envs: %w", err)
	}

	if len(cfg.Metrics) == 0 {
		return nil, errors.New("metrics were not provided")
	}

	if slices.Contains(cfg.Metrics, "*") {
		cfg.Metrics = []string{
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
			"MCacheSys",
			"MSpanInuse",
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

	ParseFlags(cfg)

	return cfg, nil
}
