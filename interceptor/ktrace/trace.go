package ktrace

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
)

// Interceptor returns an interceptor that adds OpenTelemetry tracing to the context.
func Interceptor() interceptor.Interceptor {
	return func(ctx context.Context, info interceptor.CallInfo, req, reply []any, invoker interceptor.HandleFunc) error {
		span := trace.SpanFromContext(ctx)
		if span.SpanContext().IsValid() {
			// Create a child span for this method.
			ctx, span = otel.Tracer(kod.PkgPath).Start(ctx,
				info.FullMethod,
				trace.WithSpanKind(trace.SpanKindInternal),
				trace.WithAttributes(
					attribute.String("method", info.FullMethod),
				),
			)
		}

		err := invoker(ctx, info, req, reply)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}

		span.End()

		return err
	}
}
