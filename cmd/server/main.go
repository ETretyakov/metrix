package main

import (
	"metrix/internal/infrastructure"
)

func main() {
	logger := infrastructure.NewLogger()

	// infrastructure.Load(logger)  to load envs

	storageHandler, err := infrastructure.NewStorageHandler()
	if err != nil {
		logger.LogError("%s", err)
	}

	infrastructure.Dispatch(logger, storageHandler)
}
