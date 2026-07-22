package strategy

import (
	"github.com/VladAluas/Sisyphus/internal/extractors"
	"github.com/VladAluas/Sisyphus/internal/transformers"
)

func NewStrategyRegistry() *StrategyRegistry {
	r := NewRegistry()

	// BRONZE STRATEGY
	brRegistry := extractors.NewExtractorsRegistry()
	brStrategy := BronzeStrategy{brRegistry}

	// SILVER STRATEGY
	slRegistry := transformers.NewCleanRegistry()
	slStrategy := SilverStrategy{slRegistry}

	r.Register("BRONZE_LAYER", &brStrategy)
	r.Register("SILVER_LAYER", &slStrategy)

	return r
}
