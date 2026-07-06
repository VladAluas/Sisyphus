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

func (r *Repository) CheckFailedInstance(ctx context.Context, db DBTX, batchCode string) error {
	query := `
		SELECT *
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
	`

	_, err := db.QueryContext(ctx, query, batchCode)
	if err != nil {
		return err
	}

	return nil
}

// InsertBatchRun initialises a new batch run and returns the ID
func (r *Repository) InsertBatchRun(ctx context.Context, db DBTX, batchCode string) (domain.Batch, error) {
	var batchRunID domain.Batch
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
	err := db.QueryRowContext(ctx, query, batchCode).Scan(&batchRunID.BatchID, &batchRunID.BatchRunID)
	if err != nil {
		return domain.Batch{}, err
	}

	err = db.QueryRowContext(ctx, `SELECT $1 AS BATCH_CODE`, batchCode).Scan(&batchRunID.BatchCode)
	if err != nil {
		return domain.Batch{}, err
	}

	return batchRunID, nil
}

func (r *Repository) UpdateBatchStatus(ctx context.Context, db DBTX, batchRunID domain.Batch, status string) error {
	query := `
	UPDATE dp_ctrl.LOG_BATCH_RUN
	SET STATUS = $1
	WHERE BATCH_RUN_ID = $2
	`
	rows, err := db.QueryContext(ctx, query, status, batchRunID.BatchRunID)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}
