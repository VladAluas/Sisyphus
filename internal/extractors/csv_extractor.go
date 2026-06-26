// Package extractors provides a generic extractor interface for different connectors
package extractors

import (
	"context"
	"fmt"

	"github.com/VladAluas/Sisyphus/internal/orchestrator"
)

type CsvExtractor struct{}

func (e *CsvExtractor) Extract(
	ctx  context.Context,
	task orchestrator.ModuleTask,
) error {
	fmt.Printf("Extracting CSV Data for %s", task.ModuleID)

	return nil
}
