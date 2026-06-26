// Package main is the main entry point for the application
package main

import (
	"context"
	"log"
	"sync"

	"github.com/VladAluas/Sisyphus/internal/config"
	"github.com/VladAluas/Sisyphus/internal/db"
	"github.com/VladAluas/Sisyphus/internal/extractors"
	"github.com/VladAluas/Sisyphus/internal/orchestrator"
)

func main() {
	const NumWorkers = 4
	var wg sync.WaitGroup

	cfg := config.Load()

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

	repo := orchestrator.NewRepository(pg)
	orch := orchestrator.New(repo)

	ctx := context.Background()

	tasks, err := orch.Run(ctx, "019eebcc-5b66-7d3b-a678-af05fb8eb64d")
	if err != nil {
		log.Fatal(err)
	}

	// for task := range tasks {
	// 	fmt.Printf("BatchID: %s | ModuleID: %s | LayerID: %s\n", task.BatchID, task.ModuleID, task.LayerID)
	// }

	for i := 0; i < NumWorkers; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			extractors.ExtractWorker(ctx, i, tasks)
		}(i)

		wg.Wait()
	}
}
