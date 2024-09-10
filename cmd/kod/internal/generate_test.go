package internal

import (
	"context"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	execute(t, "generate github.com/go-kod/kod/tests/graphcase/...")
}

func TestGenerateWithStruct2Interface(t *testing.T) {
	execute(t, "generate -s github.com/go-kod/kod/tests/graphcase/...")
}

func TestStartWatch(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	startWatcher(ctx, generate, []string{"github.com/go-kod/kod/tests/graphcase/..."})
}
