package kod

import (
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-kod/kod/internal/otelslog"
	"github.com/samber/lo"
)

func (k *Kod) initLog() {

	k.logLevelVar = new(slog.LevelVar)
	lo.Must0(k.logLevelVar.UnmarshalText([]byte(k.config.Log.Level)))

	jsonHandler := slog.NewJSONHandler(
		os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     k.logLevelVar,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				// Remove the directory from the source's filename.
				if a.Key == slog.SourceKey {
					source := a.Value.Any().(*slog.Source)
					source.File = customBase(source.File, 2)
					source.Function = customBase(source.Function, 1)
				}

				return a
			},
		},
	)

	var handler slog.Handler
	if k.opts.logWrapper != nil {
		handler = otelslog.NewOtelHandler(k.opts.logWrapper(jsonHandler))
	} else {
		handler = otelslog.NewOtelHandler(jsonHandler)
	}

	k.log = slog.New(handler)
	// slog.SetDefault(k.log)
}

// customBase returns the last n levels of the path.
func customBase(fullPath string, levelsToKeep int) string {
	dir, file := filepath.Split(fullPath)

	dirParts := strings.Split(dir, string(filepath.Separator))

	if len(dirParts) > levelsToKeep {
		dirParts = dirParts[len(dirParts)-levelsToKeep:]
	}

	newDir := filepath.Join(dirParts...)

	result := filepath.Join(newDir, file)

	return result
}
