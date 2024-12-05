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

type SemVer struct {
	Major int
	Minor int
	Patch int
}

func (v SemVer) String() string {
	return fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
}

// SelfVersion returns the version of the running tool binary.
func SelfVersion() (string, error) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		// Should never happen.
		return "", fmt.Errorf("tool binary must be built from a module")
	}
	return info.Main.Version, nil
}
