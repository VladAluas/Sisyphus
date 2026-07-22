// Package main is the main entry point for the application
package main

import (
	"context"
	"log"
	"os"

	"github.com/VladAluas/Sisyphus/internal/config"
	"github.com/VladAluas/Sisyphus/internal/db"
	"github.com/VladAluas/Sisyphus/internal/orchestrator"
	"github.com/VladAluas/Sisyphus/internal/repository"
	"github.com/VladAluas/Sisyphus/internal/service"
	"github.com/VladAluas/Sisyphus/internal/strategy"
	"github.com/VladAluas/Sisyphus/internal/worker"
)

func main() {
	// CLI arguments
	args := os.Args
	cfg := config.Load()

	if len(args) < 2 {
		log.Fatal("batch code missing")
	}
	batchCode := args[1]

	// Database
	pg, err := db.NewDatabase(
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
		cfg.DBUser,
		cfg.DBPass,
	)
	if err != nil {
		log.Fatal(err)
	}

	// Context
	ctx := context.Background()

	// Worker Pool
	pool := worker.New(2)


	// Repository
	repo := repository.NewRepository(pg)

	// Service
	srv := service.NewBatchService(pg, repo)

	//Create registers
	stratReg := strategy.NewStrategyRegistry()

	// Start Batch Run
	plan, err := srv.StartBatchRun(ctx, batchCode)
	if err != nil {
		log.Fatal(err)
	}

	// Orchestrator
	orch := orchestrator.New(pool, stratReg)
	err = orch.Run(ctx, plan)
	if err != nil {
		log.Fatal(err)
	}

	// Batch Status Update
	err = srv.UpdateBatchRunStatus(ctx, plan)
	if err != nil {
		log.Fatal(err)
	}
}
