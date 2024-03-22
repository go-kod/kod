package kprotovalidate

import (
	"context"
	"fmt"

	"github.com/bufbuild/protovalidate-go"
	"github.com/go-kod/kod"
	"github.com/samber/lo"
	"google.golang.org/protobuf/proto"
)

// Interceptor returns a interceptor that validates the call specified by info.
func Interceptor() kod.Interceptor {
	validate := lo.Must(protovalidate.New(
		protovalidate.WithFailFast(true),
	))

	return func(ctx context.Context, info kod.CallInfo, req, reply []any, invoker kod.HandleFunc) error {
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
