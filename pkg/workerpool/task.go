package workerpool

import "context"

type Task struct {
	F    func(ctx context.Context, data any) error
	Data any
}

func (t *Task) Process(ctx context.Context) error {
	return t.F(ctx, t.Data)
}
