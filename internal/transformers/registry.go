package transformers

import "fmt"

type CleanersRegistry struct {
	clean map[string]Cleaner
}

func NewRegistry() *CleanersRegistry {
	return &CleanersRegistry{
		clean: make(map[string]Cleaner),
	}
}

func (r *CleanersRegistry) Register(name string, e Cleaner) {
	r.clean[name] = e
}

func (r *CleanersRegistry) Get(name string) (Cleaner, error) {
	e, ok := r.clean[name]
	if !ok {
		return nil, fmt.Errorf("unknown extractor %s", name)
	}

	return e, nil
}

