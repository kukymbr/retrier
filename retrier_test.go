package retrier_test

import (
	"context"
	"testing"
	"time"

	"github.com/kukymbr/retrier"
)

func TestRetrier_Do(t *testing.T) {
	tests := []struct {
		Name       string
		GetRetrier func() retrier.Retrier
		Fn         func() error
		Assert     func(t *testing.T, took time.Duration, err error)
	}{
		{
			Name: "noop retrier",
			GetRetrier: func() retrier.Retrier {
				return retrier.NewNoop()
			},
			Fn: failingFn,
			Assert: func(t *testing.T, took time.Duration, err error) {
				errorIs(t, err, errTest)
				less(t, took, time.Millisecond)
			},
		},
		{
			Name: "linear retrier when success",
			GetRetrier: func() retrier.Retrier {
				return retrier.NewLinear(3, time.Second)
			},
			Fn: successFn,
			Assert: func(t *testing.T, took time.Duration, err error) {
				noError(t, err)
				less(t, took, time.Millisecond)
			},
		},
		{
			Name: "linear retrier when error",
			GetRetrier: func() retrier.Retrier {
				return retrier.NewLinear(3, 10*time.Millisecond)
			},
			Fn: failingFn,
			Assert: func(t *testing.T, took time.Duration, err error) {
				greaterOrEqual(t, took, 20*time.Millisecond)
				errorIs(t, err, errTest)
			},
		},
		{
			Name: "progressive retrier when success",
			GetRetrier: func() retrier.Retrier {
				return retrier.NewProgressive(3, time.Second, 2)
			},
			Fn: successFn,
			Assert: func(t *testing.T, took time.Duration, err error) {
				less(t, took, time.Millisecond)
				noError(t, err)
			},
		},
		{
			Name: "progressive retrier when error",
			GetRetrier: func() retrier.Retrier {
				return retrier.NewProgressive(3, 5*time.Millisecond, 2)
			},
			Fn: failingFn,
			Assert: func(t *testing.T, took time.Duration, err error) {
				greaterOrEqual(t, took, 15*time.Millisecond)
				errorIs(t, err, errTest)
			},
		},
		{
			Name: "progressive retrier when error with max delay",
			GetRetrier: func() retrier.Retrier {
				return retrier.New(
					retrier.WithMaxDelay(retrier.ProgressiveDelay(5*time.Millisecond, 2), 10*time.Millisecond),
					retrier.LimitAttemptsCount(5),
				)
			},
			Fn: failingFn,
			Assert: func(t *testing.T, took time.Duration, err error) {
				greaterOrEqual(t, took, 35*time.Millisecond)
				less(t, took, 55*time.Millisecond)
				errorIs(t, err, errTest)
			},
		},
		{
			Name: "when zero arguments",
			GetRetrier: func() retrier.Retrier {
				return retrier.NewLinear(0, 0)
			},
			Fn: failingFn,
			Assert: func(t *testing.T, took time.Duration, err error) {
				less(t, took, time.Millisecond)
				errorIs(t, err, errTest)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			retry := test.GetRetrier()

			start := time.Now()
			err := retry.Do(test.Fn)
			took := time.Since(start)

			test.Assert(t, took, err)
		})
	}
}

func TestRetrier_DoContext(t *testing.T) {
	tests := []struct {
		Name       string
		GetContext func() context.Context
		Fn         func() error
		Assert     func(t *testing.T, err error)
	}{
		{
			Name: "with context and error",
			GetContext: func() context.Context {
				return context.Background()
			},
			Fn: failingFn,
			Assert: func(t *testing.T, err error) {
				errorIs(t, err, errTest)
			},
		},
		{
			Name: "with context without error",
			GetContext: func() context.Context {
				return context.Background()
			},
			Fn: successFn,
			Assert: func(t *testing.T, err error) {
				noError(t, err)
			},
		},
		{
			Name: "with cancelled context with error",
			GetContext: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()

				return ctx
			},
			Fn: failingFn,
			Assert: func(t *testing.T, err error) {
				errorIs(t, err, context.Canceled)
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			t.Parallel()

			retry := retrier.NewLinear(3, 10*time.Millisecond)
			err := retry.DoContext(test.GetContext(), test.Fn)

			test.Assert(t, err)
		})
	}
}
