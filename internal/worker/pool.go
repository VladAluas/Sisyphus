package worker

import (
	"context"
	"sync"

	"github.com/VladAluas/Sisyphus/internal/domain"
	"github.com/VladAluas/Sisyphus/internal/extractors"
)

type Pool struct {
	workers  int
	registry *extractors.Registry
}

func New(workers int, registry *extractors.Registry) *Pool {
	return &Pool{
		workers,
		registry,
	}
}

func (p *Pool) Run(ctx context.Context, units []domain.ExecutionUnit) error {
	jobs := make(chan domain.ExecutionUnit)
	errs := make(chan error)

	var wg sync.WaitGroup

	for i := range p.workers {

		wg.Add(1)

		go func(workerID int) {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return

				case unit, ok := <-jobs:
					if !ok {
						return
					}

					extractor, err := p.registry.Get(unit.Source.SourceSystem)
					if err != nil {
						errs <- err
						continue
					}

					err = extractor.Extract(ctx, unit)
					errs <- err
				}
			}
		}(i)
	}

	go func() {
		defer close(jobs)

		for _, unit := range units {
			select {
			case <-ctx.Done():
				return

			case jobs <- unit:
			}
		}
	}()

	go func() {
		wg.Wait()
		close(errs)
	}()

	for err := range errs {
		if err != nil {
			return err
		}
	}

	return nil
}
