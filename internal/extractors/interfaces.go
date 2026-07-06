// Package extractors provides a generic extractor interface for different connectors
package extractors

import (
	"context"
	"log"

	"github.com/VladAluas/Sisyphus/internal/domain"
)

type Extractor interface {
	Extract(ctx context.Context, task domain.ExecutionUnit) error
}

func GetExtractor(sourceSystem string) (Extractor, error) {
	switch sourceSystem {
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

func ExtractWorker(
	ctx context.Context,
	id int,
	jobs <-chan domain.ExecutionUnit,
) {
	for task := range jobs {
		extractor, err := GetExtractor("Postgresql")
		if err != nil {
			log.Println(err)
		}

		err = extractor.Extract(ctx, task)
		if err != nil {
			log.Println(err)
		}
	}
}
