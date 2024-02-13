package paths

import (
	"path/filepath"
	"strings"
)

// CustomBase returns the last n levels of the path.
func CustomBase(fullPath string, levelsToKeep int) string {
	dir, file := filepath.Split(fullPath)

	dirParts := strings.Split(dir, string(filepath.Separator))

	if len(dirParts) > levelsToKeep {
		dirParts = dirParts[len(dirParts)-levelsToKeep:]
	}

	newDir := filepath.Join(dirParts...)

	result := filepath.Join(newDir, file)

	return result
}
