package kaccesslog

import (
	"context"
	"log/slog"

	"github.com/go-kod/kod"
)

// Interceptor returns a server interceptor that logs requests and responses.
func Interceptor() kod.Interceptor {
	return func(ctx context.Context, info kod.CallInfo, req, reply []any, invoker kod.HandleFunc) error {
		err := invoker(ctx, info, req, reply)

		attrs := []slog.Attr{
			slog.Any("req", req),
			slog.Any("reply", reply),
			slog.String("component", info.Component),
			slog.String("method", info.Method),
		}

		level := slog.LevelInfo
		if err != nil {
			level = slog.LevelError
			attrs = append(attrs, slog.String("error", err.Error()))
		}

		kod.FromContext(ctx).L(ctx).LogAttrs(ctx, level, "accesslog", attrs...)

		return err
	}
}
