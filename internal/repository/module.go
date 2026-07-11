// Package repository manages the metadata connection with the DB
package repository

import (
	"context"

	"github.com/VladAluas/Sisyphus/internal/domain"
)

func (r *Repository) FindModules(ctx context.Context, db DBTX, batchCode string, layer domain.ExecutionLayer) ([]domain.ExecutionUnit, error) {
	var modules []domain.ExecutionUnit

	query := `
		SELECT DISTINCT MODULE_ID, MODULE_CODE
		FROM dp_monitor.v_Config_BatchModuleParameters
		WHERE BATCH_CODE = $1
		  AND MODULE_LAYER_ID = $2
		`
	rows, err := db.QueryContext(ctx, query, batchCode, layer.Layer.LayerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var module domain.ExecutionUnit

		err = rows.Scan(&module.Module.ModuleID, &module.Module.ModuleCode)
		if err != nil {
			return nil, err
		}
		modules = append(modules, module)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return modules, nil
}

func (r *Repository) InsertModuleRun(ctx context.Context, db DBTX, batchRunID string, layerRunID string, module domain.ExecutionUnit) (domain.ExecutionUnit, error) {
	moduleRun := module
	query := `
	INSERT INTO dp_ctrl.LOG_MODULE_RUN
	(
		 BATCH_RUN_ID
		,LAYER_RUN_ID
		,MODULE_ID
		,STATUS
		,REPROCESS_RUN
		,START_TIME
	)
	VALUES ($1, $2, $3, 'SUCCESS', 'N', Now())
	RETURNING MODULE_RUN_ID
	`
	err := db.QueryRowContext(ctx, query, batchRunID, layerRunID, module.Module.ModuleID).Scan(&moduleRun.Module.ModuleRunID)
	if err != nil {
		return domain.ExecutionUnit{}, err
	}

	return moduleRun, nil
}
