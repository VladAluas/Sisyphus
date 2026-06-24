// Package main is the main entry point for the application
package main

import (
	"context"
	"log"

	"github.com/VladAluas/Sisyphus/internal/config"
	"github.com/VladAluas/Sisyphus/internal/db"
	"github.com/VladAluas/Sisyphus/internal/orchestrator"
)

func main() {
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

	for task := range tasks {
		log.Printf("MODULE_CODE = %s; LAYER_CODE = %s; BATCH_CODE = %s;\n",
			task.ModuleID,
			task.LayerID,
			task.BatchID,
			)
	}
}
