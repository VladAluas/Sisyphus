// Package repository manages the metadata connection with the DB
package repository

import (
	"context"

	"github.com/VladAluas/Sisyphus/internal/domain"
)

func (r *Repository) FindLayers(ctx context.Context, db DBTX, batchCode string) ([]domain.ExecutionLayer, error) {
	var layers []domain.ExecutionLayer

	query := `
	WITH T AS 
	(SELECT DISTINCT MODULE_LAYER_ID, MODULE_LAYER_CODE, PRIORITY
	FROM dp_monitor.v_Config_BatchModuleParameters
	WHERE BATCH_CODE = $1
	ORDER BY PRIORITY ASC
	)
	SELECT DISTINCT MODULE_LAYER_ID, MODULE_LAYER_CODE
	FROM T
	`
	rows, err := db.QueryContext(ctx, query, batchCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var layer domain.ExecutionLayer

		if err = rows.Scan(&layer.Layer.LayerID, &layer.Layer.LayerCode); err != nil {
			return nil, err
		}
		layers = append(layers, layer)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return layers, nil
}

func (r *Repository) InsertLayerRun(ctx context.Context, db DBTX, layer domain.ExecutionLayer, batchRunID string) (domain.ExecutionLayer, error) {
	layerRunID := layer
	query := `
	INSERT INTO dp_ctrl.LOG_LAYER_RUN
	(
		 BATCH_RUN_ID
		,LAYER_ID
		,STATUS
		,REPROCESS_RUN
		,START_TIME
	)
	VALUES ($1, $2, 'PENDING', 'N', Now())
	RETURNING LAYER_RUN_ID
	`
	err := db.QueryRowContext(ctx, query, batchRunID, layerRunID.Layer.LayerID).Scan(&layerRunID.Layer.LayerRunID)
	if err != nil {
		return domain.ExecutionLayer{}, err
	}

	return layerRunID, nil
}

func (r *Repository) CheckLayerStatus(ctx context.Context, db DBTX, batchRunID string, layerCode string) (string, error) {
	status := "SUCCESS"
	query := `
	SELECT
		MODULE_CODE
	FROM dp_monitor.v_Log_BatchModuleRun
	WHERE 1 = 1
	  AND BATCH_RUN_ID = $1
	  AND LAYER_CODE = $2
	  AND MODULE_RUN_STATUS = 'FAILED'
	`
	rows, err := db.QueryContext(ctx, query, batchRunID, layerCode)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		var md domain.Module

		err := rows.Scan(&md.ModuleCode)
		if err != nil {
			return "", err
		}
		if md.ModuleCode != "" {
			status = "FAILED"
		}
	}
	return status, nil
}

func (r *Repository) UpdateLayerStatus(ctx context.Context, db DBTX, layerRunID string, status string) error {
	query := `
	UPDATE dp_ctrl.LOG_LAYER_RUN
	SET STATUS = $1
	WHERE LAYER_RUN_ID = $2
	`

	rows, err := db.QueryContext(ctx, query, status, layerRunID)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}
