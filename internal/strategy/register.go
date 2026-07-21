package strategy

func NewStrategyRegistry() *StrategyRegistry {
	r := NewRegistry()

	r.Register("BRONZE_LAYER", &BronzeStrategy{})
	r.Register("SILVER_LAYER", &SilverStrategy{})

	return r
}

