package retrier

import (
	"fmt"
	"time"
)

func doWithRetries(fn Fn, delayFn CalcDelayFunc, allowNextFn AllowNextAttemptFunc) error {
	var (
		lastErr   error
		lastDelay time.Duration
	)

	start := time.Now()
	attemptN := -1

	for {
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
		attemptN+1, time.Since(start).String(), lastErr,
	)
}
