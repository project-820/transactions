package workerpool

import (
	"context"
	"errors"
	"fmt"
	"runtime/debug"
	"sync"
	"sync/atomic"
)

var (
	ErrClosed     = errors.New("workerpool: closed")
	ErrNotStarted = errors.New("workerpool: not started")
)

type Task func(ctx context.Context) error

type Options struct {
	Workers   int
	QueueSize int

	OnPanic func(recovered any, stack []byte)
}

type Stats struct {
	Workers  int
	QueueLen int
	QueueCap int

	InFlight  int64
	Submitted uint64
	Completed uint64
	Failed    uint64
	Panics    uint64
	Dropped   uint64

	Closed  bool
	Started bool
}

type Pool struct {
	opts Options

	tasks chan Task

	startOnce sync.Once
	closeOnce sync.Once

	ctx    context.Context
	cancel context.CancelFunc

	wg sync.WaitGroup

	started atomic.Bool
	closed  atomic.Bool

	inFlight  atomic.Int64
	submitted atomic.Uint64
	completed atomic.Uint64
	failed    atomic.Uint64
	panics    atomic.Uint64
	dropped   atomic.Uint64
}

func NewPool(opts Options) *Pool {
	if opts.Workers <= 0 {
		opts.Workers = 1
	}

	if opts.QueueSize <= 0 {
		opts.QueueSize = opts.Workers * 4
	}

	return &Pool{
		opts:  opts,
		tasks: make(chan Task, opts.QueueSize),
	}
}

func (p *Pool) Start(parent context.Context) {
	p.startOnce.Do(func() {
		p.ctx, p.cancel = context.WithCancel(parent)
		p.started.Store(true)

		for i := 0; i < p.opts.Workers; i++ {
			p.wg.Add(1)
			go p.workerLoop()
		}
	})
}

func (p *Pool) Submit(ctx context.Context, task Task) error {
	if task == nil {
		return fmt.Errorf("workerpool: nil task")
	}
	if !p.started.Load() {
		return ErrNotStarted
	}
	if p.closed.Load() {
		return ErrClosed
	}

	p.submitted.Add(1)

	select {
	case <-ctx.Done():
		p.dropped.Add(1)
		return ctx.Err()
	case <-p.ctx.Done():
		p.dropped.Add(1)
		return ErrClosed
	case p.tasks <- task:
		return nil
	}
}

func (p *Pool) TrySubmit(task Task) bool {
	if task == nil || !p.started.Load() || p.closed.Load() {
		p.dropped.Add(1)
		return false
	}

	select {
	case <-p.ctx.Done():
		p.dropped.Add(1)
		return false
	case p.tasks <- task:
		p.submitted.Add(1)
		return true
	default:
		p.dropped.Add(1)
		return false
	}
}

// Shutdown gracefully drains the queue and waits for workers.
// If ctx is canceled first, it calls StopNow and returns ctx.Err().
func (p *Pool) Shutdown(ctx context.Context) error {
	if !p.started.Load() {
		return ErrNotStarted
	}

	p.close()

	done := make(chan struct{})
	go func() {
		p.wg.Wait()
		close(done)
	}()

	select {
	case <-ctx.Done():
		p.StopNow()
		return ctx.Err()
	case <-done:
		// workers exited
		if p.cancel != nil {
			p.cancel()
		}
		return nil
	}
}

// StopNow cancels pool context immediately (tasks may stop early if they respect ctx),
// closes the queue, and waits for workers to exit.
func (p *Pool) StopNow() {
	if !p.started.Load() {
		return
	}

	if p.cancel != nil {
		p.cancel()
	}

	p.close()
	p.wg.Wait()
}

func (p *Pool) Stats() Stats {
	return Stats{
		Workers:  p.opts.Workers,
		QueueLen: len(p.tasks),
		QueueCap: cap(p.tasks),

		InFlight:  p.inFlight.Load(),
		Submitted: p.submitted.Load(),
		Completed: p.completed.Load(),
		Failed:    p.failed.Load(),
		Panics:    p.panics.Load(),
		Dropped:   p.dropped.Load(),

		Closed:  p.closed.Load(),
		Started: p.started.Load(),
	}
}

func (p *Pool) workerLoop() {
	defer p.wg.Done()

	for task := range p.tasks {
		if task == nil {
			continue
		}

		p.inFlight.Add(1)
		func() {
			defer p.inFlight.Add(-1)

			defer func() {
				if r := recover(); r != nil {
					p.panics.Add(1)
					if p.opts.OnPanic != nil {
						p.opts.OnPanic(r, debug.Stack())
					}
					p.failed.Add(1)
				}
			}()

			err := task(p.ctx)
			if err != nil {
				p.failed.Add(1)

				return
			}

			p.completed.Add(1)
		}()
	}
}

func (p *Pool) close() {
	p.closeOnce.Do(func() {
		p.closed.Store(true)
		close(p.tasks)
	})
}
