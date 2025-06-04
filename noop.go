package retrier

import "context"

type noopRetrier struct{}

func (r *noopRetrier) Do(fn Fn) error {
	return r.DoContext(context.Background(), fn)
}

func (r *noopRetrier) DoContext(_ context.Context, fn Fn) error {
	return fn()
}
