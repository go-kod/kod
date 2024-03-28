package circuitbreaker

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/sony/gobreaker"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CircuitBreaker struct {
	breaker *gobreaker.TwoStepCircuitBreaker
}

// newCircuitBreaker creates TwoStepCircuitBreaker with suitable settings.
//
// Name is the name of the CircuitBreaker.
//
// MaxRequests is the maximum number of requests allowed to pass through
// when the CircuitBreaker is half-open.
// If MaxRequests is 0, the CircuitBreaker allows only 1 request.
//
// Interval is the cyclic period of the closed state
// for the CircuitBreaker to clear the internal Counts.
// If Interval is 0, the CircuitBreaker doesn't clear internal Counts during the closed state.
//
// Timeout is the period of the open state,
// after which the state of the CircuitBreaker becomes half-open.
// If Timeout is 0, the timeout value of the CircuitBreaker is set to 60 seconds.
//
// ReadyToTrip is called with a copy of Counts whenever a request fails in the closed state.
// If ReadyToTrip returns true, the CircuitBreaker will be placed into the open state.
// If ReadyToTrip is nil, default ReadyToTrip is used.
//
// Default settings:
// MaxRequests: 3
// Interval:    5 * time.Second
// Timeout:     10 * time.Second
// ReadyToTrip: DefaultReadyToTrip
func NewCircuitBreaker(ctx context.Context, name string) *CircuitBreaker {
	return &CircuitBreaker{
		breaker: gobreaker.NewTwoStepCircuitBreaker(gobreaker.Settings{
			Name:        name,
			MaxRequests: 3,
			Interval:    5 * time.Second,
			Timeout:     10 * time.Second,
			ReadyToTrip: defaultReadyToTrip,
			OnStateChange: func(name string, from, to gobreaker.State) {
				fmt.Fprintf(os.Stderr, "circuit breaker state change: %s %s -> %s\n", name, from, to)
			},
		}),
	}
}

func (c *CircuitBreaker) Allow() (func(error), error) {
	done, err := c.breaker.Allow()
	if err != nil {
		return nil, err
	}

	return func(err error) {
		done(isSuccessful(err))
	}, nil
}

// DefaultReadyToTrip returns true when the number of consecutive failures is more than 3 and rate of failure is more than 60%.
func defaultReadyToTrip(counts gobreaker.Counts) bool {
	failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
	return counts.ConsecutiveFailures >= 3 && failureRatio >= 0.6
}

func isSuccessful(err error) bool {

	if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
		return false
	}

	switch status.Code(err) {
	case codes.Canceled,
		codes.DeadlineExceeded,
		codes.ResourceExhausted,
		codes.Aborted,
		codes.Internal,
		codes.Unavailable:
		return false
	}

	return true
}
