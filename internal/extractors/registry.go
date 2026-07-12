package extractors

import "fmt"

type Registry struct {
	extractors map[string]Extractor
}

func NewRegistry() *Registry {
	return &Registry{
		extractors: make(map[string]Extractor),
	}
}

func (r *Registry) Register(name string, e Extractor) {
	r.extractors[name] = e
}

func (r *Registry) Get(name string) (Extractor, error) {
	e, ok := r.extractors[name]
	if !ok {
		return nil, fmt.Errorf("unknown extractor %s", name)
	}

	return e, nil
}
