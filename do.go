package retrier

import (
	"context"
	"fmt"
	"time"
)

func doWithRetries(ctx context.Context, fn Fn, delayFn CalcDelayFunc, allowNextFn AllowNextAttemptFunc) error {
	var (
		lastErr   error
		lastDelay time.Duration
	)

	start := time.Now()
	attemptN := 0

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		attemptN++

		err := fn()
		lastErr = err

		if err == nil {
			return nil
		}

		if !allowNextFn(attemptN, lastDelay) {
			break
		}

		delay := delayFn(attemptN, lastDelay)
		if delay < 0 {
			delay = 0
		}

		lastDelay = delay

		if delay == 0 {
			continue
		}

		<-time.After(delay)
	}

	return fmt.Errorf(
		"failed after %d attempts (elapsed %s): %w",
		attemptN, time.Since(start).String(), lastErr,
	)
}
