package interceptor

import (
	"context"
)

// CallInfo contains information about the call.
type CallInfo struct {
	// The impl of the called component.
	Impl any
	// The component name of the called method.
	Component string
	// The full name of the called method, in the format of "package/service.method".
	FullMethod string
}

// HandleFunc is the type of the function invoked by Components.
type HandleFunc func(ctx context.Context, info CallInfo, req, reply []any) error

// Interceptor is the type of the function used to intercept Components.
type Interceptor func(ctx context.Context, info CallInfo, req, reply []any, invoker HandleFunc) error

// Condition is the type of the function used to determine whether an interceptor should be used.
type Condition func(ctx context.Context, info CallInfo) bool

// Chain converts a slice of Interceptors into a single Interceptor.
func Chain(interceptors []Interceptor) Interceptor {
	if len(interceptors) == 0 {
		return nil
	}

	return func(ctx context.Context, info CallInfo, req, reply []any, invoker HandleFunc) error {
		// Build the interceptor chain.
		chain := buildInterceptorChain(invoker, interceptors, 0)
		return chain(ctx, info, req, reply)
	}
}

// buildInterceptorChain recursively constructs a chain of interceptors.
func buildInterceptorChain(invoker HandleFunc, interceptors []Interceptor, current int) HandleFunc {
	if current == len(interceptors) {
		return invoker
	}

	return func(ctx context.Context, info CallInfo, req, reply []any) error {
		return interceptors[current](ctx, info, req, reply, buildInterceptorChain(invoker, interceptors, current+1))
	}
}

// If returns an Interceptor that only invokes the given interceptor if the given condition is true.
func If(interceptor Interceptor, condition Condition) Interceptor {
	return func(ctx context.Context, info CallInfo, req, reply []any, invoker HandleFunc) error {
		if condition(ctx, info) {
			return interceptor(ctx, info, req, reply, invoker)
		}

		return invoker(ctx, info, req, reply)
	}
}

// And groups conditions with the AND operator.
func And(first, second Condition, conditions ...Condition) Condition {
	return func(ctx context.Context, info CallInfo) bool {
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
	return func(ctx context.Context, info CallInfo) bool {
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
	return func(ctx context.Context, info CallInfo) bool {
		return !condition(ctx, info)
	}
}

// IsMethod returns a condition that checks if the method name matches the given method.
func IsMethod(method string) Condition {
	return func(_ context.Context, info CallInfo) bool {
		return info.FullMethod == method
	}
}
