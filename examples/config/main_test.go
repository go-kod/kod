package main

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-kod/kod"
)

// TestApp tests the app component.
func TestApp(t *testing.T) {
	t.Parallel()

	kod.RunTest(t, func(ctx context.Context, comp *app) {
		assert.Equal(t, "globalConfig", comp.Config().Name)
	}, kod.WithConfigFile("./kod.toml"))
}
