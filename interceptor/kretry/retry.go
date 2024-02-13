package kretry

import (
	"context"
	"fmt"

	"github.com/avast/retry-go/v4"
	"github.com/go-kod/kod"
)

// Interceptor returns a interceptor that retries the call specified by info.
func Interceptor(opts ...retry.Option) kod.Interceptor {
	return func(ctx context.Context, info kod.CallInfo, req, reply []any, invoker kod.HandleFunc) error {
		err := retry.Do(func() error {
			return invoker(ctx, info, req, reply)
		}, opts...)
		if err != nil {
			return fmt.Errorf("retry failed: %w", err)
		}

		return nil
	}
}
