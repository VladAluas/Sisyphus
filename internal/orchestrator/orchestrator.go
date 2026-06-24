// Package orchestrator manages the tasks for the ETL
package orchestrator

import "context"

type Orchestrator struct {
	repo *Repository
}

func New(repo *Repository) *Orchestrator {
	return &Orchestrator{repo}
}

func (o *Orchestrator) Run(
			ctx context.Context,
			batchRunID string,
) (<-chan ModuleTask, error) {
	tasks, err := o.repo.GetBatchModules(batchRunID)
	if err != nil {
		return nil, err
	}

	out := make(chan ModuleTask)

	go func() {
		defer close(out)

		for _, task := range tasks {
			select {
			case <-ctx.Done():
				return
			case out <- task:
			}
		}
	}()

	return out, nil
}
