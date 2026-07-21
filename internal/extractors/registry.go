package extractors

import "fmt"

type Registry struct {
	process map[string]Processor
}

func NewRegistry() *Registry {
	return &Registry{
		process: make(map[string]Processor),
	}
}

func (r *Registry) Register(name string, e Processor) {
	r.process[name] = e
}

func (r *Registry) Get(name string) (Processor, error) {
	e, ok := r.process[name]
	if !ok {
		return nil, fmt.Errorf("unknown extractor %s", name)
	}

	return e, nil
}
