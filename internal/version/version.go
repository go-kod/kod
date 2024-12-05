package version

import (
	"fmt"
	"runtime/debug"
)

const (
	// CodeGenMajor is the major version of the generated code.
	CodeGenMajor = 0
	// CodeGenMinor is the minor version of the generated code.
	CodeGenMinor = 1
	// codeGenPatch is the patch version of the generated code.
	codeGenPatch = 0
)

// CodeGenSemVersion is the version of the generated code.
var CodeGenSemVersion = SemVer{Major: CodeGenMajor, Minor: CodeGenMinor, Patch: codeGenPatch}

// SemVer represents a semantic version.
type SemVer struct {
	Major int
	Minor int
	Patch int
}

// String returns the string representation of the semantic version.
func (v SemVer) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// SelfVersion returns the version of the running tool binary.
func SelfVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		// Should never happen.
		panic("tool binary must be built from a module")
	}
	return info.Main.Version
}
