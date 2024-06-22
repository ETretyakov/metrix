package main

import (
	"metrix/internal/infrastructure"
	"os"

	"github.com/spf13/pflag"
)

func main() {
	logger := infrastructure.NewLogger()

	// infrastructure.Load(logger)  to load envs

	storageHandler, err := infrastructure.NewStorageHandler()
	if err != nil {
		logger.LogError("%s", err)
	}

	var addr string

	pflag.StringVarP(&addr, "address", "a", ":8080", "the address for the api to listen on. Host and port separated by ':'")
	pflag.Parse()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		addr = envRunAddr
	}

	infrastructure.Dispatch(addr, logger, storageHandler)
}
