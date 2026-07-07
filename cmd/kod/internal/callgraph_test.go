package internal

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraph(t *testing.T) {
	t.Run("no filepath", func(t *testing.T) {
		assert.Equal(t, "please input the binary filepath", execute(t, "callgraph"))
	})

	t.Run("unknown format", func(t *testing.T) {
		assert.Panics(t, func() {
			execute(t, "callgraph callgraph.go")
		})
	})

	for _, test := range []struct{ os, arch string }{
		{"linux", "amd64"},
		{"windows", "amd64"},
		{"darwin", "arm64"},
		{"darwin", "amd64"},
	} {
		t.Run(test.os+"_"+test.arch, func(t *testing.T) {
			bin := buildGraphcase(t, test.os, test.arch)
			dot := filepath.Join(t.TempDir(), "my-graph.dot")

			execute(t, fmt.Sprintf("callgraph %s --o %s", bin, dot))
			assert.FileExists(t, dot)

			data, err := os.ReadFile(dot)
			assert.Nil(t, err)

			assert.Contains(t, string(data), "github.com/go-kod/kod/Main")
		})
	}

	t.Run("json format", func(t *testing.T) {
		bin := buildGraphcase(t, "", "")

		data := execute(t, "callgraph "+bin+" --t json")

		assert.Contains(t, string(data), "github.com/go-kod/kod/Main")
	})
}

func buildGraphcase(t *testing.T, goos, goarch string) string {
	t.Helper()

	bin := filepath.Join(t.TempDir(), "graphcase")
	cmd := exec.Command("go", "build", "-o", bin, "./graphcase")
	cmd.Dir = "../../../tests"
	cmd.Env = os.Environ()
	if goos != "" {
		cmd.Env = append(cmd.Env, "GOOS="+goos)
	}
	if goarch != "" {
		cmd.Env = append(cmd.Env, "GOARCH="+goarch)
	}
	if err := cmd.Run(); err != nil {
		t.Fatal(err)
	}
	return bin
}
