package internal

import "testing"

func TestStruct2Interface(t *testing.T) {
	execute(t, "struct2interface ../../../tests/case1 -v")
}
