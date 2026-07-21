package extractors

import "fmt"

type ExtractorsRegistry struct {
	process map[string]Extractor
}

func NewRegistry() *ExtractorsRegistry {
	return &ExtractorsRegistry{
		process: make(map[string]Extractor),
	}
}

func (r *ExtractorsRegistry) Register(name string, e Extractor) {
	r.process[name] = e
}

func (r *ExtractorsRegistry) Get(name string) (Extractor, error) {
	e, ok := r.process[name]
	if !ok {
		return nil, fmt.Errorf("unknown extractor %s", name)
	}

	return e, nil
}
