package retrier

type noopRetrier struct{}

func (r *noopRetrier) Do(fn Fn) error {
	return fn()
}
