package extractors

import "fmt"

type ExtractorsRegistry struct {
	extractors map[string]Extractor
}

func NewRegistry() *ExtractorsRegistry {
	return &ExtractorsRegistry{
		extractors: make(map[string]Extractor),
	}
}

func (r *ExtractorsRegistry) Register(name string, e Extractor) {
	r.extractors[name] = e
}

func (r *ExtractorsRegistry) Get(name string) (Extractor, error) {
	e, ok := r.extractors[name]
	if !ok {
		return nil, fmt.Errorf("unknown extractor %s", name)
	}

	return e, nil
}
