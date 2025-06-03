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
err := rt.Do(someNiceButUnstableFunc())
```

### Progressive retrier

Linear retrier attempts to execute a function with a progressively increasing delay and a random jitter.

```go
rt := retrier.NewProgressive(5, 10*time.Second, 1.5)
err := rt.Do(someNiceButUnstableFunc())
```

### No-op retrier

Noop retrier allow to disable retrying policy using the same `Retrier` interface.

```go
rt := retrier.NewNoop()
err := rt.Do(someFuncToRunWithoutRetries())
```

### Custom retrier

To use a Retrier with some custom logic, use a `retrier.New` function:

```go
rt := retrier.New(
    func(attemptN int, lastDelay time.Duration) time.Duration {
	    return 10*time.Second 	
    },
    func(attemptN int, lastDelay time.Duration) bool {
	    return lastDelay > time.Minute
    }
)
err := rt.Do(someNiceButUnstableFunc())
```

#### Predefined sugars

There are some predefined sugar functions to create custom retriers. 
For example:

```go
rt := retrier.New(
    retrier.WithJitter(retrier.WithMaxDelay(retrier.ProgressiveDelay(10*time.Second, 1.5), 10*time.Minute)),
    retrier.LimitAttemptsCount(100),
)
err := rt.Do(someNiceButUnstableFunc())
```

##### Predefined `CalcDelayFunc` functions

* `retrier.LimitAttemptsCount` adds attempt count limit.

##### Predefined `AllowNextAttemptFunc` functions

* `retrier.FixedDelay` sets fixed delay before every attempt;
* `retrier.ProgressiveDelay` sets a progressively increasing delay;
* `retrier.WithDelayJitter` adds random jitter to the delay;
* `retrier.WithMaxDelay` adds max delay limit.

## License

[MIT](LICENSE).