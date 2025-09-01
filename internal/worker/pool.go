package worker

import (
	"context"
	"sync"

	"dtqs/internal/queue"
)

type Pool struct {
	workers []*Worker
	wg      sync.WaitGroup
}

func NewPool(size int, q queue.Queue) *Pool {
	workers := make([]*Worker, size)
	for i := 0; i < size; i++ {
		workers[i] = NewWorker(i+1, q)
	}

	return &Pool{
		workers: workers,
	}
}

func (p *Pool) Start(ctx context.Context) {
	for _, w := range p.workers {
		p.wg.Add(1)
		go func(worker *Worker) {
			defer p.wg.Done()
			worker.Start(ctx)
		}(w)
	}
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
