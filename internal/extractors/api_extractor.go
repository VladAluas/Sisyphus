// Package extractors provides a generic extractor interface for different connectors
package extractors

import (
	"context"
	"fmt"

	"github.com/VladAluas/Sisyphus/internal/orchestrator"
)

type APIExtractor struct{}

func (e *APIExtractor) Extract(
	ctx context.Context,
	task orchestrator.ModuleTask,
) error {
	fmt.Printf("Extracting Data for API: %s", task.ModuleID)

	return nil
}
