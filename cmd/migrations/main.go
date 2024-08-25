package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"metrix/internal/config"
	"metrix/pkg/logger"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/pressly/goose/v3"
)

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
)

func main() {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	flags.Parse(os.Args[1:])
	args := flags.Args()

	if len(args) < 2 {
		flags.Usage()

		return
	}

	command := args[1]

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Error(ctx, "failed to read config", err)
	}

	logger.InitDefault(cfg.LogLevel)

	if cfg.Postgres.DSN == "" {
		logger.Fatal(ctx, "postgres DSN is empty", err)
	}

	db, err := goose.OpenDBWithDriver("pgx", cfg.Postgres.DSN)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	fmt.Printf("db: %+v\n", db)

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	arguments := []string{}
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	if err := goose.RunContext(
		ctx,
		command,
		db,
		cfg.Postgres.MigrationFolder,
		arguments...,
	); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
