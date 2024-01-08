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

	level := lo.If(k.config.Env == "dev", slog.LevelDebug).Else(slog.LevelInfo)

	k.log = slog.New(otelslog.NewOtelHandler(
		slog.NewJSONHandler(
			os.Stdout, &slog.HandlerOptions{
				AddSource: true,
				Level:     level,
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
		),
	))

	slog.SetDefault(k.log)
}

func customBase(fullPath string, levelsToKeep int) string {
	// 使用filepath.Split获取目录和文件名
	dir, file := filepath.Split(fullPath)

	// 将目录拆分成各个部分
	dirParts := strings.Split(dir, string(filepath.Separator))

	// 保留指定的目录层级
	if len(dirParts) > levelsToKeep {
		dirParts = dirParts[len(dirParts)-levelsToKeep:]
	}

	// 使用filepath.Join重新组合目录
	newDir := filepath.Join(dirParts...)

	// 使用filepath.Join组合新目录和文件名
	result := filepath.Join(newDir, file)

	return result
}
