// Package orchestrator manages the tasks for the ETL
package orchestrator

import (
	"context"

	"github.com/VladAluas/Sisyphus/internal/domain"
	"github.com/VladAluas/Sisyphus/internal/worker"
)

type Orchestrator struct {
	pool *worker.Pool
}

func New(pool *worker.Pool) *Orchestrator {
	return &Orchestrator{pool}
}

func (o *Orchestrator) executeLayer(ctx context.Context, layer domain.ExecutionLayer) error {
	return o.pool.Run(ctx, layer.Modules)
}

func (o *Orchestrator) Run(ctx context.Context, plan domain.ExecutionPlan) error {
	for _, layer := range plan.Layers {

		err := o.executeLayer(ctx, layer)
		if err != nil {
			return err
		}
	}

	return nil
}
