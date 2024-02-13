package internal

import (
	"bytes"
	"strings"
	"testing"

	"github.com/samber/lo"
)

func TestCmd(t *testing.T) {
	execute("-v")
}

func execute(args string) string {

	actual := new(bytes.Buffer)
	rootCmd.SetOut(actual)
	rootCmd.SetErr(actual)
	rootCmd.SetArgs(strings.Split(args, " "))
	lo.Must0(rootCmd.Execute())

	return actual.String()
}
