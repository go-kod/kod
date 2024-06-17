package internal

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraph(t *testing.T) {
	t.Run("no filepath", func(t *testing.T) {
		assert.Equal(t, "please input the binary filepath", execute("callgraph"))
	})

	t.Run("unknown format", func(t *testing.T) {
		assert.Panics(t, func() {
			execute("callgraph callgraph.go")
		})
	})

	for _, test := range []struct{ os, arch string }{
		{"linux", "amd64"},
		{"windows", "amd64"},
		{"darwin", "arm64"},
		{"darwin", "amd64"},
	} {
		t.Run(test.os+"_"+test.arch, func(t *testing.T) {
			cmd := exec.Command("go", "build", "-o", "graphcase", "../../../tests/graphcase")
			cmd.Env = append(os.Environ(), "GOOS="+test.os, "GOARCH="+test.arch)
			assert.Nil(t, cmd.Run())

			execute("callgraph graphcase")
			assert.FileExists(t, "my-graph.dot")

			data, err := os.ReadFile("my-graph.dot")
			assert.Nil(t, err)

			assert.Contains(t, string(data), "github.com/go-kod/kod/Main")
			os.Remove("my-graph.dot")
			os.Remove("graphcase")
		})
	}

	t.Run("json format", func(t *testing.T) {
		cmd := exec.Command("go", "build", "-o", "graphcase", "../../../tests/graphcase")
		assert.Nil(t, cmd.Run())

		data := execute("callgraph graphcase --t json")

		assert.Contains(t, string(data), "github.com/go-kod/kod/Main")
		os.Remove("graphcase")
	})
}
