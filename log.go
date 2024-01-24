package kod

import (
	"context"
	"io"
	"log/slog"
	"os"

	"github.com/go-kod/kod/internal/hooks"
	"github.com/go-kod/kod/internal/otelslog"
	"github.com/go-kod/kod/internal/paths"
	"github.com/samber/lo"
	"gopkg.in/natefinch/lumberjack.v2"
)

func (k *Kod) initLog() {

	k.logLevelVar = new(slog.LevelVar)
	lo.Must0(k.logLevelVar.UnmarshalText([]byte(k.config.Log.Level)))

	// Default to stdout.
	var writer io.Writer = os.Stdout
	// If a log file is specified, use it.
	if k.config.Log.File != "" {
		logger := &lumberjack.Logger{
			Filename:   k.config.Log.File,
			MaxSize:    500, // megabytes
			MaxBackups: 7,
			MaxAge:     28, //days
			Compress:   false,
		}
		k.hooker.Add(hooks.HookFunc{
			Name: PkgPath,
			Fn: func(ctx context.Context) error {
				return logger.Close()
			},
		})
		writer = logger
	}

	jsonHandler := slog.NewJSONHandler(
		writer, &slog.HandlerOptions{
			AddSource: true,
			Level:     k.logLevelVar,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				// Remove the directory from the source's filename.
				if a.Key == slog.SourceKey {
					source := a.Value.Any().(*slog.Source)
					source.File = paths.CustomBase(source.File, 2)
					source.Function = paths.CustomBase(source.Function, 1)
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
}
