package internal

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	execute(t, "generate github.com/go-kod/kod/tests/graphcase/...")
}

func TestGenerateWithStruct2Interface(t *testing.T) {
	execute(t, "generate -s github.com/go-kod/kod/tests/graphcase/...")
}
