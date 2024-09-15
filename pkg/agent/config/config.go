// Module "config" holds configuration structures for agent configrations.
package config

import (
	"fmt"
	"slices"
	"time"

	"github.com/caarlos0/env"
	"github.com/pkg/errors"
)

// Config - the structure that keeps main agent config.
type Config struct {
	Address          string        `env:"ADDRESS"              envDefault:"localhost:8080" flag:"address"         flagShort:"a"  flagDescription:"http address"`
	PollInterval     int64         `env:"POLL_INTERVAL"        envDefault:"2"              flag:"poll_interval"   flagShort:"p"  flagDescription:"interval between polling"`
	ReportInterval   int64         `env:"REPORT_INTERVAL"      envDefault:"10"             flag:"report_interval" flagShort:"r"  flagDescription:"interval between reporting"`
	Goroutines       int64         `env:"RATE_LIMIT"           envDefault:"5"              flag:"goroutines"      flagShort:"l"  flagDescription:"number of goroutines"`
	LogLevel         string        `env:"LOG_LEVEL"            envDefault:"info"           flag:"log_level"       flagShort:"o"  flagDescription:"level for logging"`
	SignKey          string        `env:"KEY"                  envDefault:""               flag:"sign_key"        flagShort:"k"  flagDescription:"a key using for signing requests body"`
	Metrics          []string      `env:"AGT_METRICS"          envDefault:"*"`
	UseBatching      bool          `env:"USE_BATCHING"         envDefault:"true"`
	RetryCount       int64         `env:"RETRY_COUNT"          envDefault:"3"`
	RetryWaitTime    time.Duration `env:"RETRY_WAIT_TIME"      envDefault:"1s"`
	RetryMaxWaitTime time.Duration `env:"RETRY_MAX_WAIT_TIME"  envDefault:"5s"`
	CryptoKey        string        `env:"CRYPTO_KEY"                                       flag:"crypto-key"       flagShort:"i" flagDescription:"crypto key"`
	ConfigFile       string        `env:"CONFIG"`
}

// NewConfig - the builder function for Config.
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
			"TotalMemory",
			"FreeMemory",
			"CPUutilization",
		}
	}

	parseFlags(cfg)

	if cfg.ConfigFile != "" {
		if err := readFromFile(cfg.ConfigFile, cfg); err != nil {
			return nil, errors.Wrap(err, "failed to read from file")
		}
	}

	return cfg, nil
}
