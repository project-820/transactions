package workerpool

import (
	"context"
)

type WorkerPool struct {
}

func NewWorkerPool() WorkerPool {
	return WorkerPool{}
}

func (w *WorkerPool) TaskAdd(ctx context.Context, task Task) error {
	return nil
}
