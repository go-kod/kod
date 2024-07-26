package kslog

import (
	"context"
	"log/slog"
)

// levelHandler is a slog.Handler implementation that observes logs.
type levelHandler struct {
	level slog.Level
	next  slog.Handler
}

// NewLogLevelHandler returns a new LevelHandler and the observed logs.
func NewLevelHandler(level slog.Level) func(slog.Handler) slog.Handler {
	return func(h slog.Handler) slog.Handler {
		return &levelHandler{
			next:  h,
			level: level,
		}
	}
}

// Handle implements slog.Handler.
func (h *levelHandler) Handle(ctx context.Context, r slog.Record) error {
	return h.next.Handle(ctx, r)
}

// WithAttrs returns a new slog.handler with the provided attributes.
func (h *levelHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	h.next.WithAttrs(attrs)
	return h
}

// WithGroup returns a slog.handler with a group, provided the group's name.
func (h *levelHandler) WithGroup(name string) slog.Handler {
	h.next.WithGroup(name)
	return h
}

// Enabled returns true if the provided level is enabled.
func (h levelHandler) Enabled(_ context.Context, lv slog.Level) bool {
	return lv >= h.level
}
