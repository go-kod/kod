package kslog

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"
)

// removeTime removes the top-level time attribute.
// It is intended to be used as a ReplaceAttr function,
// to make example output deterministic.
func removeTime(groups []string, a slog.Attr) slog.Attr {
	if a.Key == slog.TimeKey && len(groups) == 0 {
		return slog.Attr{}
	}
	return a
}

type observer struct {
	buf *bytes.Buffer
}

func (b *observer) parse() []map[string]any {
	lines := strings.Split(b.buf.String(), "\n")

	data := make([]map[string]any, 0)
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		var m map[string]any
		if err := json.Unmarshal([]byte(line), &m); err != nil {
			panic(err)
		}

		data = append(data, m)
	}

	return data
}

// String returns the observed logs as a string.
func (b *observer) String() string {
	return b.buf.String()
}

// Len returns the number of observed logs.
func (b *observer) Len() int {
	return len(b.parse())
}

// ErrorCount returns the number of observed logs with level error.
func (b *observer) ErrorCount() int {
	return b.Filter(func(r map[string]any) bool {
		return r["level"] == slog.LevelError.String()
	}).Len()
}

// Filter returns a new observed logs with the provided filter applied.
func (b *observer) Filter(filter func(map[string]any) bool) *observer {
	var filtered []map[string]any
	for _, line := range b.parse() {
		if filter(line) {
			filtered = append(filtered, line)
		}
	}

	buf := new(bytes.Buffer)
	for _, line := range filtered {
		if err := json.NewEncoder(buf).Encode(line); err != nil {
			panic(err)
		}
	}

	return &observer{
		buf: buf,
	}
}

// Clean clears the observed logs.
func (b *observer) Clean() *observer {
	b.buf.Reset()

	return b
}

func NewTestLogger() (*slog.Logger, *observer) {
	observer := &observer{
		buf: new(bytes.Buffer),
	}
	log := slog.New(slog.NewJSONHandler(observer.buf, &slog.HandlerOptions{
		ReplaceAttr: removeTime,
	}))
	slog.SetDefault(log)

	return log, observer
}
