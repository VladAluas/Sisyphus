// Package repository manages the metadata connection with the DB
package repository

import (
	"context"

	"github.com/VladAluas/Sisyphus/internal/domain"
)

func (r *Repository) FindLayers(ctx context.Context, db DBTX, batchCode string) ([]domain.Layer, error) {
	var layerIDs []domain.Layer

	query := `
	SELECT DISTINCT MODULE_LAYER_ID, MODULE_LAYER_CODE
	FROM dp_monitor.v_Config_BatchModuleParameters
	WHERE BATCH_CODE = $1
	`
	rows, err := db.QueryContext(ctx, query, batchCode)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var layerID domain.Layer

		if err = rows.Scan(&layerID.LayerID, &layerID.LayerCode); err != nil {
			return nil, err
		}
		layerIDs = append(layerIDs, layerID)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return layerIDs, nil
}

func (r *Repository) InsertLayerRun(ctx context.Context, db DBTX, batchRunID domain.Batch, layerID domain.Layer) (domain.Layer, error) {
  layerRunID := layerID
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
	err := db.QueryRowContext(ctx, query, batchRunID.BatchRunID, layerRunID.LayerID).Scan(&layerRunID.LayerRunID)
	if err != nil {
		return domain.Layer{}, err
	}

	return layerRunID, nil
}

func (r *Repository) UpdateLayerStatus(ctx context.Context, db DBTX, layerRunID domain.Layer, status string) error {
	query := `
	UPDATE dp_ctrl.LOG_LAYER_RUN
	SET STATUS = $1
	WHERE LAYER_RUN_ID = $2
	`
	rows, err := db.QueryContext(ctx, query, status, layerRunID.LayerID)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}
