// Package strategy return an execution strategy for the worker pool
package strategy

import (
	"context"

	"github.com/VladAluas/Sisyphus/internal/domain"
)

type Strategy interface {
	Execute(ctx context.Context, unit domain.ExecutionUnit) error
}
