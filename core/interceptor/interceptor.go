package interceptor

import (
	"context"

	"github.com/go-kod/kod"
)

// Condition is the type of the function used to determine whether an interceptor should be used.
type Condition func(ctx context.Context, info kod.CallInfo) bool

// Chain converts a slice of Interceptors into a single Interceptor.
func Chain(interceptors []kod.Interceptor) kod.Interceptor {
	if len(interceptors) == 0 {
		return nil
	}

	return func(ctx context.Context, info kod.CallInfo, req, reply []any, invoker kod.HandleFunc) error {
		// Build the interceptor chain.
		chain := buildInterceptorChain(invoker, interceptors, 0)
		return chain(ctx, info, req, reply)
	}
}

// buildInterceptorChain recursively constructs a chain of interceptors.
func buildInterceptorChain(invoker kod.HandleFunc, interceptors []kod.Interceptor, current int) kod.HandleFunc {
	if current == len(interceptors) {
		return invoker
	}

	return func(ctx context.Context, info kod.CallInfo, req, reply []any) error {
		return interceptors[current](ctx, info, req, reply, buildInterceptorChain(invoker, interceptors, current+1))
	}
}

// If returns an Interceptor that only invokes the given interceptor if the given condition is true.
func If(interceptor kod.Interceptor, condition Condition) kod.Interceptor {
	return func(ctx context.Context, info kod.CallInfo, req, reply []any, invoker kod.HandleFunc) error {
		if condition(ctx, info) {
			return interceptor(ctx, info, req, reply, invoker)
		}

		return invoker(ctx, info, req, reply)
	}
}

// And groups conditions with the AND operator.
func And(first, second Condition, conditions ...Condition) Condition {
	return func(ctx context.Context, info kod.CallInfo) bool {
		if !first(ctx, info) || !second(ctx, info) {
			return false
		}
		for _, condition := range conditions {
			if !condition(ctx, info) {
				return false
			}
		}

		return true
	}
}

// Or groups conditions with the OR operator.
func Or(first, second Condition, conditions ...Condition) Condition {
	return func(ctx context.Context, info kod.CallInfo) bool {
		if first(ctx, info) || second(ctx, info) {
			return true
		}
		for _, condition := range conditions {
			if condition(ctx, info) {
				return true
			}
		}

		return false
	}
}

// Not negates the given condition.
func Not(condition Condition) Condition {
	return func(ctx context.Context, info kod.CallInfo) bool {
		return !condition(ctx, info)
	}
}
