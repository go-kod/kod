package kaccesslog

import (
	"context"
	"log/slog"

	"github.com/go-kod/kod/interceptor"
)

// Interceptor returns a server interceptor that logs requests and responses.
func Interceptor() interceptor.Interceptor {
	return func(ctx context.Context, info interceptor.CallInfo, req, reply []any, invoker interceptor.HandleFunc) error {
		err := invoker(ctx, info, req, reply)

		attrs := []slog.Attr{
			slog.Any("req", req),
			slog.Any("reply", reply),
			slog.String("method", info.Method),
		}

		level := slog.LevelInfo
		if err != nil {
			level = slog.LevelError
			attrs = append(attrs, slog.String("error", err.Error()))
		}

		// check if impl L(ctx context.Context) method
		if l, ok := info.Impl.(interface {
			L(context.Context) *slog.Logger
		}); ok {
			l.L(ctx).LogAttrs(ctx, level, "accesslog", attrs...)
		}

		return err
	}
}
