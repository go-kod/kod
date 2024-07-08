package kresty

import (
	"time"

	"github.com/go-resty/resty/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

type Client = resty.Client

type Config struct {
	BaseURL string        `toml:"baseURL"`
	Timeout time.Duration `toml:"timeout"`
}

func (c Config) Build() *Client {
	cc := resty.New()
	cc.SetBaseURL(c.BaseURL)
	cc.SetTimeout(c.Timeout)
	cc.RemoveProxy()

	cc.OnError(func(r *resty.Request, err error) {
		span := trace.SpanFromContext(r.Context())
		if span.SpanContext().IsValid() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}
	})

	cc.OnBeforeRequest(func(_ *resty.Client, r *resty.Request) error {
		span := trace.SpanFromContext(r.Context())
		if span.SpanContext().IsValid() {
			ctx, _ := otel.Tracer("github.com/go-kod/kod").Start(r.Context(), r.URL, trace.WithSpanKind(trace.SpanKindClient))
			r.SetContext(ctx)
		}

		return nil
	})

	cc.OnAfterResponse(func(_ *resty.Client, r *resty.Response) error {
		span := trace.SpanFromContext(r.Request.Context())
		if span.SpanContext().IsValid() {
			span.SetAttributes(semconv.HTTPClientAttributesFromHTTPRequest(r.Request.RawRequest)...)
			span.SetAttributes(
				semconv.HTTPStatusCodeKey.Int64(int64(r.StatusCode())),
			)
			span.End()
		}

		return nil
	})

	return cc
}
