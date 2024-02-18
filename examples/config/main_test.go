package main

import (
	"context"
	"testing"

	"github.com/go-kod/kod"
	"github.com/stretchr/testify/assert"
)

// TestApp tests the app component.
func TestApp(t *testing.T) {
	t.Parallel()

	kod.RunTest(t, func(ctx context.Context, comp *app) {
		assert.Equal(t, "config", comp.Config().Name)
	}, kod.WithConfigFile("./kod.toml"))
}
