package kpyroscope

import (
	"context"
	"fmt"
	"os"

	"github.com/go-kod/kod"
	pyroscope "github.com/grafana/pyroscope-go"
)

type Config struct {
	ServerAddress string
}

func (c Config) Build(ctx context.Context) (*pyroscope.Profiler, error) {
	k := kod.FromContext(ctx)

	p, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: k.Config().Name,
		// Logger:          pyroscope.StandardLogger,
		ServerAddress: c.ServerAddress,
		Tags:          map[string]string{"hostname": os.Getenv("HOSTNAME")},
		ProfileTypes: []pyroscope.ProfileType{
			// these profile types are enabled by default:
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			// these profile types are optional:
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to start pyroscope: %w", err)
	}

	return p, nil
}
