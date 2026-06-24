// Package orchestrator manages the tasks for the ETL
package orchestrator

import "database/sql"

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) GetBatchModules(batchRunID string) ([]ModuleTask, error) {
	var tasks []ModuleTask

	query := `
		SELECT
			 bm.BATCH_ID
			,bm.MODULE_ID
			,bm.MODULE_LAYER_ID
			,bm.PRIORITY
		FROM dp_ctrl.BATCH_MODULE bm
		INNER JOIN dp_ctrl.LOG_BATCH_RUN br ON br.BATCH_ID = bm.BATCH_ID
		WHERE 1 = 1
			AND br.BATCH_RUN_ID = $1
		ORDER BY bm.PRIORITY
	`

	rows, err := r.db.Query(query, batchRunID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var t ModuleTask

		t.BatchRunID = batchRunID

		err := rows.Scan(
			&t.BatchID,
			&t.ModuleID,
			&t.LayerID,
			&t.Priority,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, t)
	}

	return tasks, nil
}
