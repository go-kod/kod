package kvalidate

import (
	"context"
	"fmt"

	"github.com/go-kod/kod"
	"github.com/go-playground/validator/v10"
)

// Interceptor returns a interceptor that validates the call specified by info.
func Interceptor() kod.Interceptor {
	validate := validator.New()

	return func(ctx context.Context, info kod.CallInfo, req, reply []any, invoker kod.HandleFunc) error {
		for _, v := range req {
			if err := validate.Struct(v); err != nil {
				return fmt.Errorf("validate failed: %w", err)
			}
		}

		return invoker(ctx, info, req, reply)
	}
}
