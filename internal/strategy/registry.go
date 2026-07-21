package strategy

import "fmt"

type StrategyRegistry struct {
	process map[string]Strategy
}

func NewRegistry() *StrategyRegistry {
	return &StrategyRegistry{
		process: make(map[string]Strategy),
	}
}

func (r *StrategyRegistry) Register(name string, s Strategy) {
	r.process[name] = s
}

func (r *StrategyRegistry) Get(name string) (Strategy, error) {
	e, ok := r.process[name]
	if !ok {
		return nil, fmt.Errorf("unknown strategy %s", name)
	}

	return e, nil
}
