// Package extractors provides a generic extractor interface for different connectors
package extractors

import (
	"context"
	"fmt"

	"github.com/VladAluas/Sisyphus/internal/orchestrator"
)

type PostgresqlExtractor struct{}

func (e *PostgresqlExtractor) Extract(
	ctx context.Context,
	task orchestrator.ModuleTask,
) error {
	fmt.Printf("Extracting data for %s", task.ModuleID)

	return nil
}
