package internal

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	t.Run("generate basic case", func(t *testing.T) {
		err := execute(t, "generate github.com/go-kod/kod/tests/graphcase/... -v")
		require.Empty(t, err)

		// Verify generated files exist
		require.FileExists(t, filepath.Join("../../../tests/graphcase", "kod_gen.go"))
	})

	t.Run("generate with invalid path", func(t *testing.T) {
		err := execute(t, "generate invalid/path")
		require.NotEmpty(t, err)
	})

	t.Run("generate without path", func(t *testing.T) {
		err := execute(t, "generate")
		require.Empty(t, err)
	})
}

func TestGenerateWithStruct2Interface(t *testing.T) {
	t.Run("struct2interface basic case", func(t *testing.T) {
		err := execute(t, "generate -s github.com/go-kod/kod/tests/graphcase/... -v")
		require.Empty(t, err)
	})

	t.Run("struct2interface with invalid path", func(t *testing.T) {
		err := execute(t, "generate -s invalid/path")
		require.NotEmpty(t, err)
	})
}

func TestStartWatch(t *testing.T) {
	t.Run("watch with timeout", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		startWatcher(ctx, generate, []string{"github.com/go-kod/kod/tests/graphcase/..."})
	})

	t.Run("watch without path", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		startWatcher(ctx, generate, nil)
	})

	t.Run("watch with invalid path", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		startWatcher(ctx, generate, []string{"invalid/path"})
	})
}

func TestGenerateOptions(t *testing.T) {
	t.Run("generate with custom warn function", func(t *testing.T) {
		var warnings []error
		opt := Options{
			Warn: func(err error) {
				warnings = append(warnings, err)
			},
		}

		err := Generate(".", []string{"github.com/go-kod/kod/tests/graphcase/..."}, opt)
		require.NoError(t, err)
	})
}
