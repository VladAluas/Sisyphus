// Package extractors provides a generic extractor interface for different connectors
package extractors

import (
	"context"
	"fmt"

	"github.com/VladAluas/Sisyphus/internal/domain"
)

type APIExtractor struct{}

func (e *APIExtractor) Extract(
	ctx context.Context,
	task domain.ExecutionUnit,
) error {
	fmt.Printf("Extracting Data for API: %s", task.Module.ModuleID)

	return nil
}
