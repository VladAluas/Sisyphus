// Package orchestrator manages the tasks for the ETL
package orchestrator

import (
	"context"

	"github.com/VladAluas/Sisyphus/internal/domain"
	"github.com/VladAluas/Sisyphus/internal/strategy"
	"github.com/VladAluas/Sisyphus/internal/worker"
)

type Orchestrator struct {
	pool *worker.Pool
}

func New(pool *worker.Pool) *Orchestrator {
	return &Orchestrator{pool}
}

func (o *Orchestrator) executeLayer(ctx context.Context, layer domain.ExecutionLayer, strategy strategy.Strategy) error {
	return o.pool.Run(ctx, layer.Modules, strategy)
}

func (o *Orchestrator) Run(ctx context.Context, plan domain.ExecutionPlan) error {
	strat := strategy.NewStrategyRegistry()
	for _, layer := range plan.Layers {

		strategy, err := strat.Get(layer.Layer.LayerCode)
		if err != nil {
			return err
		}
		err = o.executeLayer(ctx, layer, strategy)
		if err != nil {
			return err
		}
	}

	return nil
}
