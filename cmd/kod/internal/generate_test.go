package internal

import (
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	execute("generate github.com/go-kod/kod/tests/graphcase/...")
}

func TestGenerateWithWatch(t *testing.T) {
	execute("generate -w github.com/go-kod/kod/tests/graphcase/... -t 1s")
	cmd := exec.Command("sh", "-c", "echo test>test.go")
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	assert.Nil(t, cmd.Run())
	time.Sleep(2 * time.Second)
}
