package strategy

import (
	"context"

	"github.com/VladAluas/Sisyphus/internal/domain"
	"github.com/VladAluas/Sisyphus/internal/extractors"
)

type BronzeStrategy struct {
	registry *extractors.ExtractorsRegistry
}

func (b *BronzeStrategy) Execute(
	ctx context.Context,
	unit domain.ExecutionUnit,
) error {
	extractor, err := b.registry.Get(unit.Source.SourceSystem)
	if err != nil {
		return err
	}

	return extractor.Extract(ctx, unit)
}
