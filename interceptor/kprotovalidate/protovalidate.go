package kprotovalidate

import (
	"context"
	"fmt"

	"github.com/bufbuild/protovalidate-go"
	"github.com/go-kod/kod/interceptor"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
)

// Interceptor returns a interceptor that validates the call specified by info.
func Interceptor() interceptor.Interceptor {
	validate := lo.Must(protovalidate.New(
		protovalidate.WithFailFast(true),
	))

	return func(ctx context.Context, info interceptor.CallInfo, req, reply []any, invoker interceptor.HandleFunc) error {
		for _, v := range req {
			if obj, ok := v.(proto.Message); ok {
				if err := validate.Validate(obj); err != nil {
					return fmt.Errorf("proto validate failed: %w", err)
				}
			}
		}

		return invoker(ctx, info, req, reply)
	}
}
