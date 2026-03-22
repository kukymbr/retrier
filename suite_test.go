package retrier_test

import (
	"context"
	"errors"
	"testing"
	"time"
)

var errTest = errors.New("test error")

func failingFn(_ context.Context) error {
	return errTest
}

func successFn(_ context.Context) error {
	return nil
}

func errorIs(t *testing.T, err error, target error) {
	t.Helper()

	if !errors.Is(err, target) {
		t.Errorf("got error %s, want %s", err, target)
	}
}

func noError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Errorf("got error %s, want no error", err)
	}
}

func less(t *testing.T, v1, v2 time.Duration) {
	t.Helper()

	if v1 >= v2 {
		t.Errorf("expected %s less than %s, but it is not", v1.String(), v2.String())
	}
}

func greaterOrEqual(t *testing.T, v1, v2 time.Duration) {
	t.Helper()

	if v1 < v2 {
		t.Errorf("expected %s greater or equal to %s, but it is not", v1.String(), v2.String())
	}
}
