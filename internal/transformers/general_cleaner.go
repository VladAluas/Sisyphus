package transformers

import (
	"context"
	"fmt"

	"github.com/VladAluas/Sisyphus/internal/domain"
)

type GeneralCleaner struct{}

func (p *GeneralCleaner) Process(ctx context.Context, task domain.ExecutionUnit) error {
	fmt.Printf("GENERAL_CLEANER: Cleaning data for %s\n", task.Module.ModuleCode)

	return nil
}
