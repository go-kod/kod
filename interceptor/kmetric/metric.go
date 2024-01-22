package kmetric

import (
	"context"
	"time"

	"github.com/go-kod/kod"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

var (
	methodCounts = lo.Must(otel.Meter(kod.PkgPath).Int64Counter("kod.component.count",
		metric.WithDescription("Count of Kod component method invocations"),
	))

	methodErrors = lo.Must(otel.Meter(kod.PkgPath).Int64Counter("kod.component.error",
		metric.WithDescription("Count of Kod component method invocations that result in an error"),
	))

	methodDurations = lo.Must(otel.Meter(kod.PkgPath).Float64Histogram("kod.component.duration",
		metric.WithDescription("Duration, in microseconds, of Kod component method execution"),
		metric.WithUnit("ms"),
	))
)

// Interceptor returns an interceptor that adds OpenTelemetry metrics to the context.
func Interceptor() kod.Interceptor {
	return func(ctx context.Context, info kod.CallInfo, req, reply []any, invoker kod.HandleFunc) (err error) {
		now := time.Now()

		err = invoker(ctx, info, req, reply)

		as := attribute.NewSet(
			attribute.String("component", info.Component),
			attribute.String("full_method", info.FullMethod),
		)

		if err != nil {
			methodErrors.Add(ctx, 1, metric.WithAttributeSet(as))
		}

		methodCounts.Add(ctx, 1, metric.WithAttributeSet(as))
		methodDurations.Record(ctx, float64(time.Since(now).Milliseconds()), metric.WithAttributeSet(as))

		return err
	}
}
