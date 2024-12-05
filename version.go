package kod

import (
	"github.com/go-kod/kod/internal/version"
)

// The following types are used to check, at compile time, that every
// kod_gen.go file uses the codegen API version that is linked into the binary.
type (
	// CodeGenVersion is the version of the codegen API.
	CodeGenVersion[_ any] string
	// CodeGenLatestVersion is the latest version of the codegen API.
	CodeGenLatestVersion = CodeGenVersion[[version.CodeGenMajor][version.CodeGenMinor]struct{}]
)
