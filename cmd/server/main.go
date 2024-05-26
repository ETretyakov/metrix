package main

import (
	"context"
	"log"
	"metrix/internal/config"
	"metrix/internal/infrastructure"
	"metrix/internal/logger"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		err := logger.Initialize("debug")
		if err != nil {
			logger.Log.Fatalf("failed to load config: %w", err)
		} else {
			log.Fatalf("failed to load config: %s", err)
		}
	}

	logger.Initialize(config.LogLevel)

	storageHandler, err := infrastructure.NewStorageHandler(
		context.Background(),
		config.FileStoragePath,
		config.StoreInterval,
		config.Restore,
	)
	if err != nil {
		logger.Log.Fatalf("failed to connect to the storage: %w", err)
	}

	infrastructure.Dispatch(
		config.Address,
		storageHandler,
	)
}
