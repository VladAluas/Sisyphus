package strategy

import (
	"context"

	"github.com/VladAluas/Sisyphus/internal/domain"
	"github.com/VladAluas/Sisyphus/internal/extractors"
)

type BronzeStrategy struct {}

func (b *BronzeStrategy) Execute(
	ctx context.Context,
	unit domain.ExecutionUnit,
) error {
	e := extractors.NewExtractorsRegistry()
	extractor, err := e.Get(unit.Source.SourceSystem)
	if err != nil {
		return err
	}

	return extractor.Extract(ctx, unit)
}
