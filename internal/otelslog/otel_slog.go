package otelslog

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// otelHandler is an implementation of slog's Handler interface.
// Its role is to ensure correlation between logs and OTel spans
// by:
//
// 1. Adding otel span and trace IDs to the log record.
// 2. Adding otel context baggage members to the log record.
// 3. Setting slog record as otel span event.
// 4. Adding slog record attributes to the otel span event.
// 5. Setting span status based on slog record level (only if >= slog.LevelError).
type otelHandler struct {
	// Next represents the next handler in the chain.
	Next slog.Handler
	// NoBaggage determines whether to add context baggage members to the log record.
	NoBaggage bool
}

// NewOtelHandler creates and returns a new OtelHandler to use with log/slog.
func NewOtelHandler(next slog.Handler) slog.Handler {
	return &otelHandler{
		Next: next,
	}
}

// Handle handles the provided log record and adds correlation between a slog record and an Open-Telemetry span.
func (h otelHandler) Handle(ctx context.Context, record slog.Record) error {
	if ctx == nil {
		return h.Next.Handle(ctx, record)
	}

	span := trace.SpanFromContext(ctx)
	if span == nil || !span.IsRecording() {
		return h.Next.Handle(ctx, record)
	}

	if !h.NoBaggage {
		// Adding context baggage members to log record.
		b := baggage.FromContext(ctx)
		for _, m := range b.Members() {
			record.AddAttrs(slog.String(m.Key(), m.Value()))
		}
	}

	// Adding log info to span event.
	eventAttrs := make([]attribute.KeyValue, 0)
	eventAttrs = append(eventAttrs, attribute.String("msg", record.Message))
	eventAttrs = append(eventAttrs, attribute.String("level", record.Level.String()))
	eventAttrs = append(eventAttrs, attribute.String("time", record.Time.Format(time.RFC3339Nano)))

	record.Attrs(func(attr slog.Attr) bool {
		otelAttr := h.slogAttrToOtelAttr(attr)
		if otelAttr.Valid() {
			eventAttrs = append(eventAttrs, otelAttr)
		}

		return true
	})

	span.AddEvent("log_record", trace.WithAttributes(eventAttrs...))

	// Adding span info to log record.
	spanContext := span.SpanContext()
	if spanContext.HasTraceID() {
		traceID := spanContext.TraceID().String()
		record.AddAttrs(slog.String("trace_id", traceID))
	}

	if spanContext.HasSpanID() {
		spanID := spanContext.SpanID().String()
		record.AddAttrs(slog.String("span_id", spanID))
	}

	// Setting span status if the log is an error.
	// Purposely leaving as codes.Unset (default) otherwise.
	if record.Level >= slog.LevelError {
		span.SetStatus(codes.Error, record.Message)
	}

	return h.Next.Handle(ctx, record)
}

// WithAttrs returns a new Otel whose attributes consists of handler's attributes followed by attrs.
func (h otelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return otelHandler{
		Next: h.Next.WithAttrs(attrs),
	}
}

// WithGroup returns a new Otel with a group, provided the group's name.
func (h otelHandler) WithGroup(name string) slog.Handler {
	return otelHandler{
		Next: h.Next.WithGroup(name),
	}
}

// Enabled reports whether the logger emits log records at the given context and level.
// Note: We handover the decision down to the next handler.
func (h otelHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Next.Enabled(ctx, level)
}

// slogAttrToOtelAttr converts a slog attribute to an OTel one.
// Note: returns an empty attribute if the provided slog attribute is empty.
func (h otelHandler) slogAttrToOtelAttr(attr slog.Attr, groupKeys ...string) attribute.KeyValue {
	attr.Value = attr.Value.Resolve()
	if attr.Equal(slog.Attr{}) {
		return attribute.KeyValue{}
	}

	key := func(k string, prefixes ...string) string {
		for _, prefix := range prefixes {
			k = fmt.Sprintf("%s.%s", prefix, k)
		}

		return k
	}(attr.Key, groupKeys...)

	value := attr.Value.Resolve()

	switch attr.Value.Kind() {
	case slog.KindAny:
		return attribute.String(key, fmt.Sprintf("%+v", value.Any()))
	case slog.KindBool:
		return attribute.Bool(key, value.Bool())
	case slog.KindFloat64:
		return attribute.Float64(key, value.Float64())
	case slog.KindInt64:
		return attribute.Int64(key, value.Int64())
	case slog.KindString:
		return attribute.String(key, value.String())
	case slog.KindTime:
		return attribute.String(key, value.Time().Format(time.RFC3339Nano))
	case slog.KindGroup:
		groupAttrs := value.Group()
		if len(groupAttrs) == 0 {
			return attribute.KeyValue{}
		}

		for _, groupAttr := range groupAttrs {
			return h.slogAttrToOtelAttr(groupAttr, append(groupKeys, key)...)
		}
	default:
		return attribute.KeyValue{}
	}

	return attribute.KeyValue{}
}
