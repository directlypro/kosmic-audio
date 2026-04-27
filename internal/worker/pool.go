package worker

import (
	"context"
	"errors"
	"log/slog"
)

var ErrQueueFull = errors.New("worker queue is full")

type Processor interface {
	Process(ctx context.Context, job Job) error
}

type Pool struct {
	jobs      chan Job
	workers   int
	processor Processor
}

func NewPool(workers, queueDepth int, p Processor) *Pool {
	return &Pool{
		jobs:      make(chan Job, queueDepth),
		workers:   workers,
		processor: p,
	}
}

func (p *Pool) Start(ctx context.Context) {
	for i := range p.workers {
		go func(id int) {
			slog.Info("worker started", "id", id)
			for {
				select {
				case job, ok := <-p.jobs:
					if !ok {
						slog.Info("worker shutting down", "id", id)
						return
					}
					if err := p.processor.Process(ctx, job); err != nil {
						slog.Error("job failed", "id", job.ID, "err", err)
					}
				case <-ctx.Done():
					return
				}
			}
		}(i)
	}
}

func (p *Pool) Submit(job Job) error {
	select {
	case p.jobs <- job:
		return nil
	default:
		return ErrQueueFull
	}
}

func (p *Pool) Stop() {
	close(p.jobs)
}
