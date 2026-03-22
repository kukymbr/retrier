package retrier

import (
	"fmt"
	"math/rand/v2"
	"time"
)

// Option configuring the Retrier instance.
type Option func(*Retrier)

// CalcDelayFunc is a function to calculate a delay before next attempt.
type CalcDelayFunc func(attemptN int, lastDelay time.Duration) time.Duration

// CalcDelayJitterFunc is a function returning a delay jitter coefficient (0..1 is recommended).
type CalcDelayJitterFunc func(attemptN int, lastDelay time.Duration) float64

// AllowNextAttemptFunc is a function to decide if next attempt is allowed.
type AllowNextAttemptFunc func(attemptN int, lastDelay time.Duration) bool

// LogFunc is a function logging debug messages.
type LogFunc func(format string, args ...any)

// WithName sets a Retrier name. Used for debug messages.
func WithName(name string) Option {
	return func(r *Retrier) {
		r.name = name
	}
}

// WithAllowNextAttemptFn adds a function deciding run next attempt or not to the chain.
func WithAllowNextAttemptFn(fn AllowNextAttemptFunc) Option {
	return func(r *Retrier) {
		r.allowNextAttemptFn = append(r.allowNextAttemptFn, fn)
	}
}

// WithCalcDelayFn sets a function calculating timeout before next attempt.
func WithCalcDelayFn(fn CalcDelayFunc) Option {
	return func(r *Retrier) {
		r.calcDelayFn = fn
	}
}

// WithCalcDelayJitterFn sets function calculating the random jitter to the delay.
func WithCalcDelayJitterFn(calcJitterFn CalcDelayJitterFunc) Option {
	return func(r *Retrier) {
		r.calcDelayJitterFn = calcJitterFn
	}
}

// WithDelayJitter adds random 0..1 delay jitter multiplier.
func WithDelayJitter() Option {
	return WithCalcDelayJitterFn(func(_ int, _ time.Duration) float64 {
		//nolint:gosec // no cryptographic purposes.
		return rand.Float64()
	})
}

// WithMaxAttemptsCount sets a function limiting the max attempts count.
func WithMaxAttemptsCount(attempts int) Option {
	if attempts < 1 {
		attempts = 1
	}

	return WithAllowNextAttemptFn(func(attemptN int, _ time.Duration) bool {
		return attemptN < attempts
	})
}

// WithoutAttempts sets a function disabling any attempts after fail (no retries).
func WithoutAttempts() Option {
	return WithAllowNextAttemptFn(func(attemptN int, _ time.Duration) bool {
		return false
	})
}

// WithoutDelay sets a function disabling delay.
func WithoutDelay() Option {
	return WithFixedDelay(0)
}

// WithFixedDelay sets a function returning the same timeout for every attempt.
func WithFixedDelay(delay time.Duration) Option {
	return WithCalcDelayFn(func(_ int, _ time.Duration) time.Duration {
		return delay
	})
}

// WithProgressiveDelay sets a function returning a progressively increasing delay.
func WithProgressiveDelay(initialDelay time.Duration, multiplier float64) Option {
	return WithCalcDelayFn(func(attemptN int, _ time.Duration) time.Duration {
		return initialDelay * time.Duration(float64(attemptN)*multiplier)
	})
}

// WithMaxDelay sets the maximum possible delay between attempts.
func WithMaxDelay(maxDelay time.Duration) Option {
	return func(r *Retrier) {
		r.maxDelay = maxDelay
	}
}

// WithLogFunc sets a function to log debug messages with.
func WithLogFunc(fn LogFunc) Option {
	return func(r *Retrier) {
		r.logFn = fn
	}
}

// WithPrintLogger sets a function printing debug messages to the std out.
func WithPrintLogger() Option {
	return func(r *Retrier) {
		r.logFn = func(format string, args ...any) {
			fmt.Printf(format, args...)
		}
	}
}

// WithNoopLogger disables debug messages.
func WithNoopLogger() Option {
	return func(r *Retrier) {
		r.logFn = func(format string, args ...any) {}
	}
}
