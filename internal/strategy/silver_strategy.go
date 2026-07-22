package strategy

import (
	"context"

	"github.com/VladAluas/Sisyphus/internal/domain"
	"github.com/VladAluas/Sisyphus/internal/transformers"
)

type SilverStrategy struct {
	registry *transformers.CleanersRegistry
}

func NewSilverStrategy(registry *transformers.CleanersRegistry) *SilverStrategy {
	return &SilverStrategy{registry}
}

func (s SilverStrategy) Execute(
	ctx context.Context,
	unit domain.ExecutionUnit,
) error {
	cleaner, err := s.registry.Get("Clean")
	if err != nil {
		return err
	}

	return cleaner.Clean(ctx, unit)
}
