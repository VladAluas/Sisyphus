// Package extractors provides a generic extractor interface for different connectors
package extractors

import (
	"context"
	"fmt"

	"github.com/VladAluas/Sisyphus/internal/domain"
)

type APIExtractor struct{}

func (e *APIExtractor) Process(
	ctx context.Context,
	task domain.ExecutionUnit,
) error {
	fmt.Printf("API_EXTRACTOR: Extracting data for %s\n", task.Module.ModuleCode)

	return nil
}
