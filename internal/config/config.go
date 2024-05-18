package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env"
	"github.com/spf13/pflag"
)

type Config struct { //nolint:govet // I want it be pretty
	Address  string `env:"ADDRESS" mapstructure:"ADDRESS" envDefault:"localhost:8080"` //nolint:lll // I want it be pretty
	LogLevel string `env:"LogLevel"   mapstructure:"LogLevel"   envDefault:"info"`
}

func LoadConfig() (*Config, error) {
	config := &Config{}
	if err := env.Parse(config); err != nil {
		return nil, fmt.Errorf("env.Parse: %w", err)
	}

	var addr string
	pflag.StringVarP(&addr, "address", "a", "", "the address for the api to listen on. Host and port separated by ':'")

	var logLevel string
	pflag.StringVarP(&logLevel, "loglevel", "r", "", "the level of logger")

	envAddress := os.Getenv("ADDRESS")
	if len(envAddress) == 0 && addr != "" {
		config.Address = addr
	}

	envLogLevel := os.Getenv("LogLevel")
	if len(envLogLevel) == 0 && envLogLevel != "" {
		config.LogLevel = envLogLevel
	}

	return config, nil
}
