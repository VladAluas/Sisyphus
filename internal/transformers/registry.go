package transformers

import "fmt"

type CleanersRegistry struct {
	cleaners map[string]Cleaner
}

func NewRegistry() *CleanersRegistry {
	return &CleanersRegistry{
		cleaners: make(map[string]Cleaner),
	}
}

func (r *CleanersRegistry) Register(name string, e Cleaner) {
	r.cleaners[name] = e
}

func (r *CleanersRegistry) Get(name string) (Cleaner, error) {
	e, ok := r.cleaners[name]
	if !ok {
		return nil, fmt.Errorf("unknown extractor %s", name)
	}

	return e, nil
}

