// Package retrier designed to execute some valuable functionality
// with retrying policy.
package retrier

import (
	"context"
	"time"
)

// Fn is a function to call.
type Fn func() error

// Retrier is a tool to execute Fn with retries.
type Retrier interface {
	// Do execute a Fn and retry it in case of error.
	Do(fn Fn) error

	// DoContext execute a Fn and retry it in case of error, respecting the context.
	DoContext(ctx context.Context, fn Fn) error
}

// New returns custom Retrier.
func New(calcDelayFn CalcDelayFunc, allowNextAttemptFn AllowNextAttemptFunc) Retrier {
	if calcDelayFn == nil {
		panic("calcDelayFn cannot be nil")
	}

	if allowNextAttemptFn == nil {
		panic("allowNextAttemptFn cannot be nil")
	}

	return &retrier{
		calcDelayFn:        calcDelayFn,
		allowNextAttemptFn: allowNextAttemptFn,
	}
}

// NewLinear returns a Retrier, attempting to execute a Fn with a fixed delay.
func NewLinear(attempts uint, delay time.Duration) Retrier {
	//nolint:gosec
	return New(FixedDelay(delay), LimitAttemptsCount(int(attempts)))
}

// NewProgressive returns a Retrier, attempting to execute a Fn
// with a progressively increasing delay and a random jitter.
func NewProgressive(attempts uint, initialDelay time.Duration, multiplier float64) Retrier {
	//nolint:gosec
	return New(
		WithDelayJitter(ProgressiveDelay(initialDelay, multiplier)),
		LimitAttemptsCount(int(attempts)),
	)
}

// NewNoop returns a Retrier without retries.
func NewNoop() Retrier {
	return &noopRetrier{}
}
