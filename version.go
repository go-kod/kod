package kod

import (
	"github.com/go-kod/kod/internal/version"
)

type (
	// The following types are used to check, at compile time, that every
	// kod_gen.go file uses the codegen API version that is linked into the binary.
	CodeGenVersion[_ any] string
	CodeGenLatestVersion  = CodeGenVersion[[version.CodeGenMajor][version.CodeGenMinor]struct{}]
)
