package kvalidate

import (
	"context"
	"fmt"

	"github.com/go-playground/validator/v10"

	"github.com/go-kod/kod/interceptor"
)

// Interceptor returns a interceptor that validates the call specified by info.
func Interceptor() interceptor.Interceptor {
	validate := validator.New()

	return func(ctx context.Context, info interceptor.CallInfo, req, reply []any, invoker interceptor.HandleFunc) error {
		for _, v := range req {
			if err := validate.Struct(v); err != nil {
				return fmt.Errorf("validate failed: %w", err)
			}
		}

		return invoker(ctx, info, req, reply)
	}
}
