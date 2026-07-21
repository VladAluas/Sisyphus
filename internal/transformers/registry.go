package transformers

import "fmt"

type CleanersRegistry struct {
	process map[string]Cleaner
}

func NewRegistry() *CleanersRegistry {
	return &CleanersRegistry{
		process: make(map[string]Cleaner),
	}
}

func (r *CleanersRegistry) Register(name string, e Cleaner) {
	r.process[name] = e
}

func (r *CleanersRegistry) Get(name string) (Cleaner, error) {
	e, ok := r.process[name]
	if !ok {
		return nil, fmt.Errorf("unknown extractor %s", name)
	}

	return e, nil
}

