package extractors

import "fmt"

type ExtractorsRegistry struct {
	extract map[string]Extractor
}

func NewRegistry() *ExtractorsRegistry {
	return &ExtractorsRegistry{
		extract: make(map[string]Extractor),
	}
}

func (r *ExtractorsRegistry) Register(name string, e Extractor) {
	r.extract[name] = e
}

func (r *ExtractorsRegistry) Get(name string) (Extractor, error) {
	e, ok := r.extract[name]
	if !ok {
		return nil, fmt.Errorf("unknown extractor %s", name)
	}

	return e, nil
}
