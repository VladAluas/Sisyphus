// Package repository manages the metadata connection with the DB
package repository

import (
	"context"

	"github.com/VladAluas/Sisyphus/internal/domain"
)

func (r *Repository) FindModules(ctx context.Context,
                                 db DBTX,
                                 batch domain.Batch,
                                 layer domain.Layer,
                                 ) ([]domain.Module, error) {
	var modules []domain.Module

	query := `
		SELECT DISTINCT MODULE_ID, MODULE_CODE
		FROM dp_monitor.v_Config_BatchModuleParameters
		WHERE BATCH_CODE = $1
		  AND MODULE_LAYER_ID = $2
		`
	rows, err := db.QueryContext(ctx, query, batch.BatchCode, layer.LayerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var module domain.Module

		err = rows.Scan(&module.ModuleID, &module.ModuleCode)
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

func (r *Repository) InserModuleRun(ctx    context.Context,
                                    db     DBTX,
                                    batch  domain.Batch,
                                    layer  domain.Layer,
                                    module domain.Module,
                                    ) (domain.Module, error) {
	moduleRunID := module
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
	VALUES ($1, $2, $3, 'PENDING', 'N', Now())
	RETURNING MODULE_RUN_ID
	`
	err := db.QueryRowContext(ctx, query, batch.BatchRunID, layer.LayerRunID, module.ModuleID).Scan(&moduleRunID.ModuleRunID)
	if err != nil {
		return domain.Module{}, err
	}

	return moduleRunID, nil
}
