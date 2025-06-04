package retrier

import (
	"context"
	"math/rand/v2"
	"time"
)

// CalcDelayFunc is a function to calculate a delay before next attempt.
type CalcDelayFunc func(attemptN int, lastDelay time.Duration) time.Duration

// AllowNextAttemptFunc is a function to decide if next attempt is allowed.
type AllowNextAttemptFunc func(attemptN int, lastDelay time.Duration) bool

// LimitAttemptsCount returns as AllowNextAttemptFunc with attempts count limit.
func LimitAttemptsCount(attempts int) AllowNextAttemptFunc {
	if attempts < 1 {
		attempts = 1
	}

	return func(attemptN int, _ time.Duration) bool {
		return attemptN < attempts
	}
}

// FixedDelay returns a CalcDelayFunc with a fixed delay value.
func FixedDelay(delay time.Duration) CalcDelayFunc {
	return func(_ int, _ time.Duration) time.Duration {
		return delay
	}
}

// ProgressiveDelay returns a CalcDelayFunc with a progressively increasing delay.
func ProgressiveDelay(initialDelay time.Duration, multiplier float64) CalcDelayFunc {
	return func(attemptN int, _ time.Duration) time.Duration {
		return initialDelay * time.Duration(float64(attemptN)*multiplier)
	}
}

// WithDelayJitter adds random jitter to the delay.
func WithDelayJitter(calcDelayFn CalcDelayFunc) CalcDelayFunc {
	return func(attemptN int, lastDelay time.Duration) time.Duration {
		delay := calcDelayFn(attemptN, lastDelay)
		//nolint:gosec
		jitter := rand.Float64()

		return time.Duration(float64(delay) * (1 + jitter))
	}
}

// WithMaxDelay adds max delay check to the existing CalcDelayFunc.
func WithMaxDelay(calcDelayFn CalcDelayFunc, maxDelay time.Duration) CalcDelayFunc {
	return func(attemptN int, lastDelay time.Duration) time.Duration {
		delay := calcDelayFn(attemptN, lastDelay)

		return min(maxDelay, delay)
	}
}

type retrier struct {
	calcDelayFn        CalcDelayFunc
	allowNextAttemptFn AllowNextAttemptFunc
}

func (r *retrier) Do(ctx context.Context, fn Fn) error {
	return doWithRetries(ctx, fn, r.calcDelayFn, r.allowNextAttemptFn)
}
