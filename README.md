# 🤦‍♂️ Retrier

[![License](https://img.shields.io/github/license/kukymbr/retrier.svg)](https://github.com/kukymbr/retrier/blob/main/LICENSE)
[![Release](https://img.shields.io/github/release/kukymbr/retrier.svg)](https://github.com/kukymbr/retrier/releases/latest)
[![GoDoc](https://godoc.org/github.com/kukymbr/retrier?status.svg)](https://godoc.org/github.com/kukymbr/retrier)
[![GoReport](https://goreportcard.com/badge/github.com/kukymbr/retrier)](https://goreportcard.com/report/github.com/kukymbr/retrier)

Retrier is a tiny zero-dependency [golang](https://go.dev) library 
to do some valuable things which may fail with retrying policy.

## Installation

The `go get` is a path.

```shell
go get github.com/kukymbr/retrier
```

## Usage

### Linear retrier

Linear retrier attempts to execute function with a fixed delay.

```go
rt := retrier.NewLinear(5, 10*time.Second)
err := rt.Do(context.Background(), someNiceButUnstableFunc)
```

### Progressive retrier

Progressive retrier attempts to execute a function with a progressively increasing delay and a random jitter.

```go
rt := retrier.NewProgressive(5, 10*time.Second, 1.5)
err := rt.Do(context.Background(), someNiceButUnstableFunc)
```

### No-op retrier

Noop retrier allow to disable retrying policy.

```go
rt := retrier.NewNoop()
err := rt.Do(context.Background(), someFuncToRunWithoutRetries)
```

### Options

To make a custom Retrier, use a `retrier.New` function:

```go
rt := retrier.New(
    // Function, calculating the delay before next attempt.
    WithCalcDelayFn(func(attemptN int, lastDelay time.Duration) time.Duration {
	    return 10*time.Second 	
    }),
	// Function, deciding give it another try or not.
    WithAllowNextAttemptFn(func(attemptN int, lastDelay time.Duration) bool {
	    return lastDelay > time.Minute
    }),
)
err := rt.Do(context.Background(), someNiceButUnstableFunc)
```

Also, options are available in the `NewLinear` and `NewProgressive` functions too:

```go
rt := retrier.NewLinear(
	5, 10*time.Second,
	WithName("test_retrier"),
	WithLogFn(t.Logf),
)
```

See the [options.go](options.go) file for an available options.

## License

[MIT](LICENSE).