// Package files contains file-related utilities.
package internal

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriter(t *testing.T) {
	file := filepath.Join(t.TempDir(), "test")
	w := NewWriter(file)

	_, _ = w.Write([]byte("hello"))

	data, err := os.ReadFile(w.tmpName)
	require.Nil(t, err)
	require.Equal(t, "hello", string(data))

	require.Nil(t, w.Close())

	data, err = os.ReadFile(file)
	require.Nil(t, err)
	require.Equal(t, "hello", string(data))
}

func TestWriter1(t *testing.T) {
	file := filepath.Join(t.TempDir(), "test")
	w := NewWriter(file)

	_, _ = w.Write([]byte("hello"))

	data, err := os.ReadFile(w.tmpName)
	require.Nil(t, err)
	require.Equal(t, "hello", string(data))

	w.Cleanup()
}

func TestWriter2(t *testing.T) {
	file := filepath.Join(t.TempDir(), "test")
	w := NewWriter(file)

	w.err = os.ErrExist

	_, _ = w.Write([]byte("hello"))

	data, err := os.ReadFile(w.tmpName)
	require.Nil(t, err)
	require.Equal(t, "", string(data))

	w.Cleanup()
}

func TestWriter3(t *testing.T) {
	file := filepath.Join(t.TempDir(), "test")
	w := NewWriter(file)

	w.Cleanup()

	_, err := w.Write([]byte("hello"))

	require.NotNil(t, err)
}

func TestWriter4(t *testing.T) {
	file := filepath.Join(t.TempDir(), "test")
	w := NewWriter(file)

	w.err = os.ErrExist

	err := w.Close()

	require.NotNil(t, err)
}

func TestWriter5(t *testing.T) {
	file := filepath.Join(t.TempDir(), "test")
	w := NewWriter(file)

	w.Cleanup()

	err := w.Close()

	require.NotNil(t, err)
}
