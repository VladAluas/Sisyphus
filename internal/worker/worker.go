// Package worker provides the interface for the orchestrator to organise e
package worker

import (
	"context"

	"github.com/VladAluas/Sisyphus/internal/domain"
)

type Processor interface {
	Process(ctx context.Context, layer domain.ExecutionUnit) error
}
