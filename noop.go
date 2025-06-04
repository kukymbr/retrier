package retrier

import "context"

type noopRetrier struct{}

func (r *noopRetrier) Do(_ context.Context, fn Fn) error {
	return fn()
}
