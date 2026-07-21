// Package main is the main entry point for the application
package main

import (
	"context"
	"log"
	"os"

	"github.com/VladAluas/Sisyphus/internal/config"
	"github.com/VladAluas/Sisyphus/internal/db"
	"github.com/VladAluas/Sisyphus/internal/extractors"
	"github.com/VladAluas/Sisyphus/internal/orchestrator"
	"github.com/VladAluas/Sisyphus/internal/repository"
	"github.com/VladAluas/Sisyphus/internal/service"
	"github.com/VladAluas/Sisyphus/internal/worker"
)

func main() {
	// const NumWorkers = 4
	// var wg sync.WaitGroup

	args := os.Args
	cfg := config.Load()

	if len(args) < 2 {
		log.Fatal("batch code missing")
	}
	batchCode := args[1]

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

	registry := extractors.NewDefaultRegistry()
	ctx := context.Background()
	repo := repository.NewRepository(pg)
	pool := worker.New(2, registry)
	orch := orchestrator.New(pool)
	srv := service.NewBatchService(pg, repo)
	plan, err := srv.StartBatchRun(ctx, batchCode)
	if err != nil {
		log.Fatal(err)
	}

	err = orch.Run(ctx, plan)
	if err != nil {
		log.Fatal(err)
	}

	err = srv.UpdateBatchRunStatus(ctx, plan)
	if err != nil {
		log.Fatal(err)
	}
	// for i := 0; i < NumWorkers; i++ {
	// 	wg.Add(1)
	//
	// 	go func(id int) {
	// 		defer wg.Done()
	//
	// 		extractors.ExtractWorker(ctx, i, tasks)
	// 	}(i)
	//
	// 	wg.Wait()
	// }
}
