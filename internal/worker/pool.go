package worker

import (
	"context"
	"sync"

	"github.com/VladAluas/Sisyphus/internal/domain"
)

type Pool struct {
	workers int
	worker  Processor
}

func New(workers int, worker Processor) *Pool {
	return &Pool{
		workers,
		worker,
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
					err := p.worker.Process(ctx, unit)
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
