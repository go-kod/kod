package kslog

import (
	"context"
	"log/slog"
	"sync"

	"github.com/samber/lo"
)

// ObservedLogs is a structure that holds observed logs.
type ObservedLogs struct {
	mu      sync.Mutex
	records []slog.Record
}

// Len returns the number of observed logs.
func (l *ObservedLogs) Len() int {
	l.mu.Lock()
	defer l.mu.Unlock()

	return len(l.records)
}

// Filter returns a new observed logs with the provided filter applied.
func (l *ObservedLogs) Filter(filter func(slog.Record) bool) *ObservedLogs {
	l.mu.Lock()
	defer l.mu.Unlock()

	return &ObservedLogs{
		mu: sync.Mutex{},
		records: lo.Filter(l.records, func(r slog.Record, index int) bool {
			return filter(r)
		}),
	}
}

// Clean clears the observed logs.
func (l *ObservedLogs) Clean() *ObservedLogs {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.records = make([]slog.Record, 0)

	return l
}

// All returns all the observed logs.
func (l *ObservedLogs) All() []slog.Record {
	l.mu.Lock()
	defer l.mu.Unlock()

	return l.records
}

// add adds a record to the observed logs.
func (l *ObservedLogs) add(r slog.Record) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.records = append(l.records, r)
}

// observer is a slog.Handler implementation that observes logs.
type observer struct {
	observedLogs *ObservedLogs
	next         slog.Handler
}

// NewLogObserver returns a new observer and the observed logs.
func NewLogObserver() (func(slog.Handler) slog.Handler, *ObservedLogs) {
	logs := &ObservedLogs{
		mu:      sync.Mutex{},
		records: make([]slog.Record, 0),
	}

	return func(h slog.Handler) slog.Handler {
		return &observer{
			next:         h,
			observedLogs: logs,
		}
	}, logs
}

// Handle implements slog.Handler.
func (h *observer) Handle(ctx context.Context, r slog.Record) error {
	h.observedLogs.add(r)

	return h.next.Handle(ctx, r)
}

// WithAttrs returns a new slog.handler with the provided attributes.
func (h *observer) WithAttrs(attrs []slog.Attr) slog.Handler {

	h.next.WithAttrs(attrs)
	return h
}

// WithGroup returns a slog.handler with a group, provided the group's name.
func (h *observer) WithGroup(name string) slog.Handler {
	h.next.WithGroup(name)

	return h
}

// Enabled returns true if the provided level is enabled.
func (h *observer) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}
