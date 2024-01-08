package kcircuitbreaker

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/internal/circuitbreaker"
	"github.com/go-kod/kod/internal/singleton"
)

var (
	once sync.Once
	pool *singleton.Singleton[circuitbreaker.CircuitBreaker]
)

// Interceptor returns an interceptor do circuit breaker.
func Interceptor() kod.Interceptor {
	once.Do(func() {
		pool = singleton.NewSingleton[circuitbreaker.CircuitBreaker]()
	})

	return func(ctx context.Context, info kod.CallInfo, req, reply []any, invoker kod.HandleFunc) error {
		breaker := pool.Get(info.FullMethod, func() *circuitbreaker.CircuitBreaker {
			return circuitbreaker.NewCircuitBreaker(ctx, info.FullMethod)
		})

		done, err := breaker.Allow()
		if err != nil {
			return fmt.Errorf("kcircuitbreaker: %w", err)
		}

		err = invoker(ctx, info, req, reply)

		done(err)

		return err
	}
}
