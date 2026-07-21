package strategy

import (
	"context"

	"github.com/VladAluas/Sisyphus/internal/domain"
	"github.com/VladAluas/Sisyphus/internal/extractors"
)

type BronzeStrategy struct {
	extractor extractors.ExtractorsRegistry
}

func (b *BronzeStrategy) Execute(
	ctx context.Context,
	unit domain.ExecutionUnit,
) error {
	extractor, err := b.extractor.Get(unit.Source.SourceSystem)
	if err != nil {
		return err
	}

	return extractor.Extract(ctx, unit)
}
