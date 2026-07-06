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

func (s *BatchService) StartBatchRun(ctx context.Context, batchCode string) (domain.Batch, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return domain.Batch{}, err
	}
	defer tx.Rollback()

	// Check if the batch is already RUNNING
	err = s.repo.CheckRunningInstance(ctx, tx, batchCode)
	if err != nil {
		return domain.Batch{}, err
	}
	fmt.Printf("%s is not running. Proceeding with new run\n", batchCode)

	// If there is no batch running at the time, create a new batchRun
	batchRunID, err := s.repo.InsertBatchRun(ctx, tx, batchCode)
	if err != nil {
		return domain.Batch{}, err
	}

	// Insert the layers for the new batch
	layerIDs, err := s.repo.FindLayers(ctx, tx, batchCode)
	if err != nil {
		return domain.Batch{}, err
	}

	var layerRunIDs []domain.Layer
	for _, layerID := range layerIDs {
		lrID, err := s.repo.InsertLayerRun(ctx, tx, batchRunID, layerID)
		if err != nil {
			return domain.Batch{}, err
		}
		layerRunIDs = append(layerRunIDs, lrID)
	}

	// Insert the Modules for each Layer
	for _, layerRun := range layerRunIDs {
		modules, err := s.repo.FindModules(ctx, tx, batchRunID, layerRun)
		if err != nil {
			return domain.Batch{}, err
		}
		for _, module := range modules {
			module, err := s.repo.InserModuleRun(ctx, tx, batchRunID, layerRun, module)
			if err != nil {
				return domain.Batch{}, err
			}
			fmt.Printf("BatchCode: %s; BatchID: %s; LayerID: %s; LayerCode: %s; LayerRun: %s; ModuleID: %s; ModuleCode: %s; ModuleRun: %s;\n", batchRunID.BatchCode, batchRunID.BatchID, layerRun.LayerID, layerRun.LayerCode, layerRun.LayerRunID, module.ModuleID, module.ModuleCode, module.ModuleRunID)
		}
	}

	// COMMIT transactions
	if err = tx.Commit(); err != nil {
		return domain.Batch{}, err
	}
	return batchRunID, nil
}

func (s *BatchService) UpdateBatchRunStatus(ctx context.Context, batch domain.Batch) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	status, err := s.repo.UpdateLayerStatus(ctx, tx, )

	err = s.repo.UpdateBatchStatus(ctx, tx, batch, "SUCCESS")
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}
