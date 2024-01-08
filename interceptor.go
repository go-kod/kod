package kod

import "context"

// CallInfo contains information about the call.
type CallInfo struct {
	// The component name of the called method.
	Component string
	// The full name of the called method, in the format of "package/service.method".
	FullMethod string
	// The name of the caller.
	Caller string
}

// HandleFunc is the type of the function invoked by Components.
type HandleFunc func(ctx context.Context, info CallInfo, req, reply []any) error

// Interceptor is the type of the function used to intercept Components.
type Interceptor func(ctx context.Context, info CallInfo, req, reply []any, invoker HandleFunc) error
