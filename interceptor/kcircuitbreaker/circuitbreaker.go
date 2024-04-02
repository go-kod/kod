package kcircuitbreaker

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/interceptor/internal/circuitbreaker"
	"github.com/go-kod/kod/internal/singleton"
)

var (
	once sync.Once
	pool *singleton.Singleton[circuitbreaker.CircuitBreaker]
)

// Interceptor returns an interceptor do circuit breaker.
func Interceptor() interceptor.Interceptor {
	once.Do(func() {
		pool = singleton.NewSingleton[circuitbreaker.CircuitBreaker]()
	})

	return func(ctx context.Context, info interceptor.CallInfo, req, reply []any, invoker interceptor.HandleFunc) error {
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
