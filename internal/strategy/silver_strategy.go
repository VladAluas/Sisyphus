package strategy

import (
	"context"

	"github.com/VladAluas/Sisyphus/internal/domain"
	"github.com/VladAluas/Sisyphus/internal/transformers"
)

type SilverStrategy struct {
	registry *transformers.CleanersRegistry
}

func (s SilverStrategy) Execute(
	ctx context.Context,
	unit domain.ExecutionUnit,
) error {

	c := transformers.NewCleanRegistry()
	cleaner, err := c.Get("Clean")
	if err != nil {
		return err
	}

	return cleaner.Clean(ctx, unit)
}
