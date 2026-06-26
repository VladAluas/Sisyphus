// Package extractors provides a generic extractor interface for different connectors
package extractors

import (
	"context"
	"log"

	"github.com/VladAluas/Sisyphus/internal/orchestrator"
)

type Extractor interface {
	Extract(ctx context.Context, task orchestrator.ModuleTask) error
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
	jobs <-chan orchestrator.ModuleTask,
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
