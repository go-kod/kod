package ktimeout

import (
	"context"
	"time"

	"github.com/go-kod/kod/interceptor"
)

type Options struct {
	Timeout time.Duration
}

// WithTimeout sets the timeout for the interceptor.
func WithTimeout(timeout time.Duration) func(*Options) {
	return func(options *Options) {
		options.Timeout = timeout
	}
}

const defaultTimeout = time.Second * 5

// Interceptor returns an interceptor that adds OpenTelemetry tracing to the context.
func Interceptor(opts ...func(*Options)) interceptor.Interceptor {
	options := Options{
		Timeout: defaultTimeout,
	}

	for _, o := range opts {
		o(&options)
	}

	return func(ctx context.Context, info interceptor.CallInfo, req, reply []any, invoker interceptor.HandleFunc) error {
		ctx, cancel := context.WithTimeout(ctx, options.Timeout)
		defer cancel()
		return invoker(ctx, info, req, reply)
	}
}
