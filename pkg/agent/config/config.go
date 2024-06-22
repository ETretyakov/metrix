package config

import (
	"errors"
	"fmt"
	"os"
	"slices"
	"time"

	"github.com/caarlos0/env"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/pflag"
)

const (
	stageLogKey = "Stage"
	stageLogVal = "loading-config"
)

type Config struct { //nolint:govet // I want it be pretty
	Address        string   `env:"ADDRESS"         mapstructure:"ADDRESS"         envDefault:"localhost:8080"` //nolint:lll // I want it be pretty
	PollInterval   int      `env:"POLL_INTERVAL"   mapstructure:"POLL_INTERVAL"   envDefault:"2"`
	ReportInterval int      `env:"REPORT_INTERVAL" mapstructure:"REPORT_INTERVAL" envDefault:"10"`
	Metrics        []string `env:"AGT_METRICS"     mapstructure:"AGT_METRICS"     envDefault:"*"`
}

func LoadConfig() (*Config, error) {
	log.Logger = log.Output( //nolint:reassign // documentation approach
		zerolog.ConsoleWriter{
			Out:        os.Stderr,
			TimeFormat: time.RFC3339,
		},
	)

	log.Info().Caller().Str(stageLogKey, stageLogVal).
		Msg("started config loading")
	config := &Config{}
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("env.Parse: %w", err)
	}

	if len(config.Metrics) == 0 {
		return nil, errors.New("metrics were not provided")
	}

	if slices.Contains(config.Metrics, "*") {
		log.Info().Caller().Str(stageLogKey, stageLogVal).
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

	var addr string
	pflag.StringVarP(&addr, "address", "a", "", "the address for the api to listen on. Host and port separated by ':'")

	var pollInterval int
	pflag.IntVarP(&pollInterval, "poll interval", "r", 0, "the number of seconds - interval between polling")

	var reportInterval int
	pflag.IntVarP(&reportInterval, "report interval", "p", 0, "the number of seconds - interval between reporting")
	pflag.Parse()

	envAddress := os.Getenv("ADDRESS")
	if len(envAddress) == 0 && addr != "" {
		config.Address = addr
	}

	envPollInterval := os.Getenv("POLL_INTERVAL")
	if len(envPollInterval) == 0 && pollInterval != 0 {
		config.PollInterval = pollInterval
	}

	envReportInterval := os.Getenv("REPORT_INTERVAL")
	if len(envReportInterval) == 0 && reportInterval != 0 {
		config.ReportInterval = reportInterval
	}

	log.Info().Caller().Str(stageLogKey, stageLogVal).
		Msg(fmt.Sprintf("config: %+v", config))

	return config, nil
}
