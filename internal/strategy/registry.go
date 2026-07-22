package strategy

import "fmt"

type StrategyRegistry struct {
	strategy map[string]Strategy
}

func NewRegistry() *StrategyRegistry {
	return &StrategyRegistry{
		strategy: make(map[string]Strategy),
	}
}

func (r *StrategyRegistry) Register(name string, s Strategy) {
	r.strategy[name] = s
}

func (r *StrategyRegistry) Get(name string) (Strategy, error) {
	e, ok := r.strategy[name]
	if !ok {
		return nil, fmt.Errorf("unknown strategy %s", name)
	}

	return e, nil
}
