package internal

import (
	"bytes"
	"strings"
	"testing"
)

func TestCmd(t *testing.T) {
	execute(t, "-v")
}

func execute(t *testing.T, args string) string {
	t.Helper()

	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	rootCmd.SetArgs(strings.Split(args, " "))
	err := rootCmd.Execute()
	if err != nil {
		t.Fatal(err)
	}

	return actual.String()
}
