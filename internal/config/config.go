// Module "config" aggregates all the necessary structures and functions that
// enables the service to read environment variables and arguments.
package config

import (
	"fmt"
	"net"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

// AppMode - the string-based type that defines running mode for the service,
// allowed values are: prod, stage, dev, local.
type AppMode string

// ProdAppMode - the constant for "prod" AppMode.
// StageAppMode - the constant for "stage" AppMode.
// DevAppMode - the constant for "dev" AppMode.
// LocalAppMode - the constant for "local" AppMode.
const (
	ProdAppMode  AppMode = "prod"
	StageAppMode AppMode = "stage"
	DevAppMode   AppMode = "dev"
	LocalAppMode AppMode = "local"
)

// Postgres - the structure for postgresql config.
type Postgres struct {
	DSN             string        `env:"DSN"              envDefault:""`
	MaxOpenConn     int           `env:"MAX_OPEN_CONN"    envDefault:"10"`
	IdleConn        int           `env:"MAX_IDLE_CONN"    envDefault:"10"`
	PingInterval    time.Duration `env:"DURATION"         envDefault:"5s"`
	MigrationFolder string        `env:"MIGRATION_FOLDER" envDefault:"./migrations"`
}

// Config - the structure for general config.
type Config struct {
	AppMode              AppMode  `env:"APP_MODE"          envDefault:"local"           flag:"mode"              flagShort:"m" flagDescription:"application mode"`
	HTTPAddress          string   `env:"ADDRESS"           envDefault:"localhost:8080"  flag:"address"           flagShort:"a" flagDescription:"http address"`
	StoreInterval        int64    `env:"STORE_INTERVAL"    envDefault:"300"             flag:"store_interval"    flagShort:"s" flagDescription:"interval for storage backup"`
	FileStoragePath      string   `env:"FILE_STORAGE_PATH" envDefault:""                flag:"file_storage_path" flagShort:"f" flagDescription:"filepath storage backup"`
	Restore              bool     `env:"RESTORE"           envDefault:"false"           flag:"restore"           flagShort:"r" flagDescription:"boolean to restore from backup"`
	LogLevel             string   `env:"LOG_LEVEL"         envDefault:"info"            flag:"log_level"         flagShort:"l" flagDescription:"level for logging"`
	LogFile              string   `env:"LOG_FILE"          envDefault:"logs/logs.jsonl" flag:"log_file"          flagShort:"w" flagDescription:"filepath for logs"`
	Postgres             Postgres `envPrefix:"DATABASE_"                                flag:"pg_dsn"            flagShort:"d" flagDescription:"database dsn"`
	SignKey              string   `env:"KEY"                                            flag:"sign_key"          flagShort:"k" flagDescription:"a key using for signing"`
	CryptoKey            string   `env:"CRYPTO_KEY"                                     flag:"crypto-key"        flagShort:"i" flagDescription:"crypto key"`
	ConfigFile           string   `env:"CONFIG"`
	TrustedSubNet        string   `env:"TRUSTED_SUBNET"    envDefault:"192.168.1.0/24"   flag:"trusted-subnet"    flagShort:"t" flagDescription:"trusted subnet variable"`
	TrustedSubNetDefined *net.IPNet
}

// NewConfig - the builder function for new configuration.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	if cfg.ConfigFile != "" {
		if err := readFromFile(cfg.ConfigFile, cfg); err != nil {
			return nil, errors.Wrap(err, "failed to read from file")
		}
	}

	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("failed to parse server envs: %w", err)
	}

	parseFlags(cfg)

	if _, TrustedSubNetDefined, err := net.ParseCIDR(cfg.TrustedSubNet); err != nil {
		return cfg, errors.Wrap(err, "failed to define subnet")
	} else {
		cfg.TrustedSubNetDefined = TrustedSubNetDefined
	}

	return cfg, nil
}
