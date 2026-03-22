// Package retrier designed to execute some valuable functionality
// with retrying policy.
package retrier

import (
	"context"
	"fmt"
	"time"
)

// Retrier is a tool to do valuable things with retries.
type Retrier struct {
	name               string
	calcDelayFn        CalcDelayFunc
	calcDelayJitterFn  CalcDelayJitterFunc
	allowNextAttemptFn []AllowNextAttemptFunc
	maxDelay           time.Duration
	logFn              LogFunc
}

// New returns a custom Retrier.
func New(opts ...Option) *Retrier {
	r := &Retrier{}

	for _, opt := range opts {
		opt(r)
	}

	if r.calcDelayFn == nil {
		panic("calcDelayFn cannot be nil")
	}

	if len(r.allowNextAttemptFn) == 0 {
		panic("allowNextAttemptFn chain cannot be empty")
	}

	if r.logFn == nil {
		r.logFn = func(string, ...any) {}
	}

	if r.name == "" {
		r.name = "retrier"
	}

	return r
}

// NewLinear returns a Retrier, attempting to execute a function with a fixed delay.
func NewLinear(attempts int, delay time.Duration, opts ...Option) *Retrier {
	opts = append([]Option{
		WithMaxAttemptsCount(attempts),
		WithFixedDelay(delay),
	}, opts...)

	return New(opts...)
}

// NewProgressive returns a Retrier, attempting to execute a function
// with a progressively increasing delay and a random jitter.
func NewProgressive(attempts int, initialDelay time.Duration, multiplier float64, opts ...Option) *Retrier {
	opts = append([]Option{
		WithMaxAttemptsCount(attempts),
		WithProgressiveDelay(initialDelay, multiplier),
		WithDelayJitter(),
	}, opts...)

	return New(opts...)
}

// NewNoop returns a Retrier without retries.
func NewNoop() *Retrier {
	return New(WithoutAttempts(), WithoutDelay(), WithNoopLogger())
}

// Do tries to execute a function and runs the retry attempts as configured.
//
//nolint:funlen
func (r *Retrier) Do(ctx context.Context, fn func(context.Context) error) error {
	var (
		lastErr   error
		lastDelay time.Duration
	)

	start := time.Now()
	attemptN := 0

attemptsLoop:
	for {
		attemptN++

		r.logf("running attempt %d", attemptN)

		if err := ctx.Err(); err != nil {
			r.logf("attempt %d: context error: %s", attemptN, err.Error())

			return err
		}

		err := fn(ctx)
		lastErr = err

		if err == nil {
			r.logf("attempt %d: succeeded", attemptN)

			return nil
		}

		for i, allowNextFn := range r.allowNextAttemptFn {
			if !allowNextFn(attemptN, lastDelay) {
				r.logf("attempt %d: disallowed to proceed (rule #%d)", attemptN, i)

				break attemptsLoop
			}
		}

		delay := r.calcDelay(attemptN, lastDelay)

		r.logf("attempt %d: calculated delay: %s", attemptN, delay.String())

		lastDelay = delay

		select {
		case <-ctx.Done():
			err := ctx.Err()

			r.logf("attempt %d: context error: %s", attemptN, err.Error())

			return err
		case <-time.After(delay):
			continue
		}
	}

	err := fmt.Errorf(
		"failed after %d attempts (elapsed %s): %w",
		attemptN, time.Since(start).String(), lastErr,
	)

	r.logf("failed: %s", err.Error())

	return err
}

func (r *Retrier) calcDelay(attemptN int, lastDelay time.Duration) time.Duration {
	delay := r.calcDelayFn(attemptN, lastDelay)
	if delay < 0 {
		delay = 0
	}

	if r.calcDelayJitterFn != nil {
		delay = time.Duration(float64(delay) * (1 + r.calcDelayJitterFn(attemptN, lastDelay)))
	}

	if r.maxDelay > 0 {
		delay = min(r.maxDelay, delay)
	}

	return delay
}

func (r *Retrier) logf(format string, args ...any) {
	r.logFn(r.name+": "+format, args...)
}
