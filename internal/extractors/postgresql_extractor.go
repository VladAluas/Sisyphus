// Package extractors provides a generic extractor interface for different connectors
package extractors

import (
	"context"
	"fmt"

	"github.com/VladAluas/Sisyphus/internal/domain"
)

type PostgresqlExtractor struct{}

func (e *PostgresqlExtractor) Extract(
	ctx context.Context,
	task domain.ExecutionUnit,
) error {
	fmt.Printf("Extracting data for %s", task.Module.ModuleCode)

	return nil
}
