// Package extractors provides a generic extractor interface for different connectors
package extractors

import (
	"context"
	"fmt"

	"github.com/VladAluas/Sisyphus/internal/domain"
)

type PostgresqlExtractor struct{}

func (p *PostgresqlExtractor) Process(
	ctx context.Context,
	task domain.ExecutionUnit,
) error {
	fmt.Printf("POSTGRESQL_EXTRACTOR: Extracting data for %s\n", task.Module.ModuleCode)

	return nil
}
