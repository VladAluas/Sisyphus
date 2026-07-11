// Package repository manages the metadata connection with the DB
package repository

import (
	"context"
	"fmt"

	"github.com/VladAluas/Sisyphus/internal/domain"
)

// CheckRunningInstance makes sure that no other instance of a particular batch is running at the moment the process is triggered
func (r *Repository) CheckRunningInstance(ctx context.Context, db DBTX, batchCode string) error {
	query := `
		SELECT BATCH_CODE
		FROM dp_ctrl.LOG_BATCH_RUN       lbr
		LEFT JOIN dp_ctrl.LOG_LAYER_RUN  llr ON llr.BATCH_RUN_ID = lbr.BATCH_RUN_ID
		LEFT JOIN dp_ctrl.LOG_MODULE_RUN lmr ON lmr.BATCH_RUN_ID = lbr.BATCH_RUN_ID
		LEFT JOIN dp_ctrl.BATCH            b ON b.BATCH_ID       = lbr.BATCH_ID
		LEFT JOIN dp_ctrl.LAYER            l ON l.LAYER_ID       = llr.LAYER_ID
		LEFT JOIN dp_ctrl.MODULE           m ON m.MODULE_ID      = lmr.MODULE_ID
		WHERE 1 = 1
			AND b.BATCH_CODE = $1
			AND (lbr.STATUS = 'RUNNING'
				OR llr.STATUS = 'RUNNING'
				OR lmr.STATUS = 'RUNNING'
					)
	`

	rows, err := db.QueryContext(ctx, query, batchCode)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var b domain.Batch

		err := rows.Scan(&b.BatchCode)
		if err != nil {
			return err
		}
		if b.BatchCode != "" {
			return fmt.Errorf("%s is still running, retry later", b.BatchCode)
		}
	}
	return nil
}

// CheckReprocessInstance checks to see if there are any instances that need to be run once more with the same setup
// The reprocessing can be due to a failure of a manual setup due to business logic
func (r *Repository) CheckReprocessInstance(ctx context.Context, db DBTX, batchCode string) (string, error) {
	// I need to return some kind of value if this is finds something
	// Not sure what to return just now
	// I would need a separate function to return the execution plan. I want to be able to check first, and if there are failed instances, then return the plan
	var b domain.Batch

	query := `
	WITH T AS
	(
		SELECT lbr.BATCH_RUN_ID, lbr.START_TIME
		FROM dp_ctrl.LOG_BATCH_RUN       lbr
		LEFT JOIN dp_ctrl.LOG_LAYER_RUN  llr ON llr.BATCH_RUN_ID = lbr.BATCH_RUN_ID
		LEFT JOIN dp_ctrl.LOG_MODULE_RUN lmr ON lmr.BATCH_RUN_ID = lbr.BATCH_RUN_ID
		LEFT JOIN dp_ctrl.BATCH            b ON b.BATCH_ID       = lbr.BATCH_ID
		LEFT JOIN dp_ctrl.LAYER            l ON l.LAYER_ID       = llr.LAYER_ID
		LEFT JOIN dp_ctrl.MODULE           m ON m.MODULE_ID      = lmr.MODULE_ID
		WHERE 1 = 1
			AND b.BATCH_CODE = $1
			AND (lbr.STATUS        = 'FAILED'
				OR llr.STATUS        = 'FAILED'
				OR lmr.STATUS        = 'FAILED'
				OR lbr.REPROCESS_RUN = 'Y'
				OR llr.REPROCESS_RUN = 'Y'
				OR lmr.REPROCESS_RUN = 'Y'
					)
		ORDER BY lbr.START_TIME
		LIMIT 1
	)
	SELECT BATCH_RUN_ID FROM T
	`

	rows, err := db.QueryContext(ctx, query, batchCode)
	if err != nil {
		return "", err
	}

	if rows.Next() {
		err := rows.Scan(&b.BatchRunID)
		if err != nil {
			return "", err
		}
	}

	return b.BatchRunID, nil
}

func (r *Repository) ReprocessInstance(ctx context.Context, db DBTX, batchRunID string) (domain.ExecutionPlan, error) {
	var exec domain.ExecutionPlan
	
	exec.Batch.BatchRunID = batchRunID


	return exec, nil
}


// InsertBatchRun initialises a new batch run and returns the ID
func (r *Repository) InsertBatchRun(ctx context.Context, db DBTX, batchCode string) (domain.ExecutionPlan, error) {
	var batchRunID domain.ExecutionPlan
	query := `
	INSERT INTO dp_ctrl.LOG_BATCH_RUN
	(
		 BATCH_ID
		,STATUS
		,REPROCESS_RUN
		,START_TIME
	)
	SELECT DISTINCT
		 BATCH_ID
		,'RUNNING'
		,'N'
		,Now()
	FROM dp_monitor.v_Config_BatchModuleParameters
	WHERE 1 = 1
	  AND BATCH_CODE = $1
	RETURNING BATCH_ID, BATCH_RUN_ID
	`
	err := db.QueryRowContext(ctx, query, batchCode).Scan(&batchRunID.Batch.BatchID, &batchRunID.Batch.BatchRunID)
	if err != nil {
		return domain.ExecutionPlan{}, err
	}

	err = db.QueryRowContext(ctx, `SELECT $1 AS BATCH_CODE`, batchCode).Scan(&batchRunID.Batch.BatchCode)
	if err != nil {
		return domain.ExecutionPlan{}, err
	}

	return batchRunID, nil
}

func (r *Repository) CheckBatchStatus(ctx context.Context, db DBTX, batchRunID string) (string, error) {
	status := "SUCCESS"
	query := `
	SELECT 
		LAYER_CODE
	FROM dp_monitor.v_Log_BatchModuleRun
	WHERE 1 = 1
	  AND BATCH_RUN_ID = $1
	  AND LAYER_RUN_STATUS = 'FAILED'
	`
	rows, err := db.QueryContext(ctx, query, batchRunID)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	if rows.Next() {
		var ly domain.Layer

		err := rows.Scan(&ly.LayerCode)
		if err != nil {
			return "", err
		}
		if ly.LayerCode != "" {
			status = "FAILED"
		}
	}
	return status, nil
}

func (r *Repository) UpdateBatchStatus(ctx context.Context, db DBTX, batchRunID string, status string) error {
	query := `
	UPDATE dp_ctrl.LOG_BATCH_RUN
	SET STATUS = $1
	WHERE BATCH_RUN_ID = $2
	`
	rows, err := db.QueryContext(ctx, query, status, batchRunID)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}
