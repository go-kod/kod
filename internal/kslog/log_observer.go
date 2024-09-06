package kslog

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"strings"

	"github.com/samber/lo"
)

// NewTestLogger returns a new test logger.
func NewTestLogger() (*slog.Logger, *observer) {
	observer := &observer{
		buf: new(bytes.Buffer),
	}
	log := slog.New(slog.NewJSONHandler(observer.buf, nil))
	slog.SetDefault(log)

	return log, observer
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
		lo.Must0(json.Unmarshal([]byte(line), &m))

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
		lo.Must0(json.NewEncoder(buf).Encode(line))
	}

	return &observer{
		buf: buf,
	}
}

// RemoveKeys removes the provided keys from the observed logs.
func (b *observer) RemoveKeys(keys ...string) *observer {
	filtered := make([]map[string]any, 0)
	for _, line := range b.parse() {
		for _, key := range keys {
			delete(line, key)
		}

		filtered = append(filtered, line)
	}

	buf := new(bytes.Buffer)
	for _, line := range filtered {
		lo.Must0(json.NewEncoder(buf).Encode(line))
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
