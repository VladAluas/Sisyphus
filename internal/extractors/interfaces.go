// Package extractors provides a generic extractor interface for different connectors
package extractors

import (
	"context"

	"github.com/VladAluas/Sisyphus/internal/domain"
)

type Processor interface {
	Process(ctx context.Context, task domain.ExecutionUnit) error
}

func Get(unit domain.ExecutionUnit) (Processor, error) {
	switch unit.Source.SourceSystem {
	case "Postgresql":
		return &PostgresqlExtractor{}, nil
	case "CSV":
		return &CsvExtractor{}, nil
	case "API":
		return &APIExtractor{}, nil
	default:
		return &PostgresqlExtractor{}, nil
	}
}
