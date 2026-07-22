package strategy

import "fmt"

type StrategyRegistry struct {
	strategies map[string]Strategy
}

func NewRegistry() *StrategyRegistry {
	return &StrategyRegistry{
		strategies: make(map[string]Strategy),
	}
}

func (r *StrategyRegistry) Register(name string, s Strategy) {
	r.strategies[name] = s
}

func (r *StrategyRegistry) Get(name string) (Strategy, error) {
	e, ok := r.strategies[name]
	if !ok {
		return nil, fmt.Errorf("unknown strategy %s", name)
	}

	return e, nil
}
