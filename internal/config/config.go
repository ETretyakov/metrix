package config

import (
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env"
	"github.com/spf13/pflag"
)

type Config struct { //nolint:govet // I want it be pretty
	Address         string        `env:"ADDRESS"           mapstructure:"ADDRESS"           envDefault:"localhost:8080"` //nolint:lll // I want it be pretty
	StoreInterval   time.Duration `env:"STORE_INTERVAL"    mapstructure:"STORE_INTERVAL"    envDefault:"300s"`
	FileStoragePath string        `env:"FILE_STORAGE_PATH" mapstructure:"FILE_STORAGE_PATH" envDefault:""`
	Restore         bool          `env:"RESTORE"           mapstructure:"RESTORE"           envDefault:"false"`
	LogLevel        string        `env:"LOG_LEVEL"         mapstructure:"LOG_LEVEL"         envDefault:"info"`
}

func LoadConfig() (*Config, error) {
	config := &Config{}
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("env.Parse: %w", err)
	}

	var addr string
	pflag.StringVarP(
		&addr,
		"address",
		"a",
		"",
		"the address for the api to listen on. Host and port separated by ':'",
	)

	var storeInterval time.Duration
	pflag.DurationVarP(
		&storeInterval,
		"store_interval",
		"i",
		time.Second*300,
		"the store_interval for database backup",
	)

	var fileStoragePath string
	pflag.StringVarP(
		&fileStoragePath,
		"file-storage-path",
		"f",
		"",
		"the filepath to save memory storage",
	)

	var restore bool
	pflag.BoolVarP(
		&restore,
		"restore",
		"r",
		true,
		"the bool for if restore memory storage",
	)

	var logLevel string
	pflag.StringVarP(
		&logLevel,
		"loglevel",
		"l",
		"",
		"the level of logger",
	)

	pflag.Parse()

	envAddress := os.Getenv("ADDRESS")
	if len(envAddress) == 0 && addr != "" {
		config.Address = addr
	}

	envStoreInterval := os.Getenv("STORE_INTERVAL")
	if len(envStoreInterval) == 0 {
		config.StoreInterval = storeInterval
	}

	envFileStoragePath := os.Getenv("FILE_STORAGE_PATH")
	if len(envFileStoragePath) == 0 && fileStoragePath != "" {
		config.FileStoragePath = fileStoragePath
	}

	envRestore := os.Getenv("RESTORE")
	if len(envRestore) == 0 && restore {
		config.Restore = restore
	}

	envLogLevel := os.Getenv("LOG_LEVEL")
	if len(envLogLevel) == 0 && logLevel != "" {
		config.LogLevel = envLogLevel
	}

	return config, nil
}
