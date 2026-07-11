// Package service helps run the repository
package service

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/VladAluas/Sisyphus/internal/domain"
	"github.com/VladAluas/Sisyphus/internal/repository"
)

type BatchService struct {
	db   *sql.DB
	repo *repository.Repository
}

func NewBatchService(db *sql.DB, repo *repository.Repository) *BatchService {
	return &BatchService{db, repo}
}

func (s *BatchService) StartBatchRun(ctx context.Context, batchCode string) (domain.ExecutionPlan, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.ExecutionPlan{}, err
	}
	defer tx.Rollback()

	// Check if the batch is already RUNNING
	err = s.repo.CheckRunningInstance(ctx, tx, batchCode)
	if err != nil {
		return domain.ExecutionPlan{}, err
	}

	brID, err := s.repo.CheckReprocessInstance(ctx, tx, batchCode)
	if err != nil {
		return domain.ExecutionPlan{}, err
	}
	//
	// Check for failed Batches from previous runs
	fmt.Printf("Checking for previous failed runs of %s\n", batchCode)
	if brID != "" {
		fmt.Printf("There was a failed batch: %s\n", brID)
	} else {
		fmt.Print("No failed Batch found. Continuing...\n")
	}

	// If there is no batch running at the time, create a new batchRun
	fmt.Printf("%s is not running. Proceeding with new run\n", batchCode)

	executionPlan, err := s.repo.InsertBatchRun(ctx, tx, batchCode)
	if err != nil {
		return domain.ExecutionPlan{}, err
	}

	// Insert the layers for the new batch
	layers, err := s.repo.FindLayers(ctx, tx, batchCode)
	if err != nil {
		return domain.ExecutionPlan{}, err
	}

	var layerRuns []domain.ExecutionLayer

	for _, layer := range layers {
		var layerRun domain.ExecutionLayer
		layerRun, err := s.repo.InsertLayerRun(ctx, tx, layer, executionPlan.Batch.BatchRunID)
		if err != nil {
			return domain.ExecutionPlan{}, err
		}

		// Insert the Modules for each Layer
		modules, err := s.repo.FindModules(ctx, tx, executionPlan.Batch.BatchCode, layerRun)
		if err != nil {
			return domain.ExecutionPlan{}, err
		}

		var moduleRuns []domain.ExecutionUnit
		for _, module := range modules {
			module, err := s.repo.InsertModuleRun(ctx, tx, executionPlan.Batch.BatchRunID, layerRun.Layer.LayerRunID, module)
			if err != nil {
				return domain.ExecutionPlan{}, err
			}
			moduleRuns = append(moduleRuns, module)
		}
		layerRun.Modules = moduleRuns
		layerRuns = append(layerRuns, layerRun)
	}

	executionPlan.Layers = layerRuns

	// COMMIT transactions
	if err = tx.Commit(); err != nil {
		return domain.ExecutionPlan{}, err
	}
	return executionPlan, nil
}

func (s *BatchService) UpdateBatchRunStatus(ctx context.Context, batch domain.ExecutionPlan) error {
	b := batch
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, ly := range b.Layers {
		status, err := s.repo.CheckLayerStatus(ctx, tx, b.Batch.BatchRunID, ly.Layer.LayerCode)
		if err != nil {
			return err
		}
		err = s.repo.UpdateLayerStatus(ctx, tx, ly.Layer.LayerRunID, status)
		if err != nil {
			return err
		}
	}

	status, err := s.repo.CheckBatchStatus(ctx, tx, b.Batch.BatchRunID)
	if err != nil {
		return err
	}

	err = s.repo.UpdateBatchStatus(ctx, tx, b.Batch.BatchRunID, status)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
