package kratelimit

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/interceptor/internal/ratelimit"
	"github.com/go-kod/kod/internal/singleton"
)

var (
	once sync.Once
	pool *singleton.Singleton[*ratelimit.Ratelimit]
)

// Interceptor returns an interceptor do rate limit.
func Interceptor() interceptor.Interceptor {
	once.Do(func() {
		pool = singleton.New[*ratelimit.Ratelimit]()
	})

	return func(ctx context.Context, info interceptor.CallInfo, req, reply []any, invoker interceptor.HandleFunc) error {
		limitor := pool.Get(info.FullMethod, func() *ratelimit.Ratelimit {
			return ratelimit.NewLimiter(ctx, info.FullMethod)
		})

		done, err := limitor.Allow()
		if err != nil {
			return fmt.Errorf("kratelimit: %w", err)
		}

		err = invoker(ctx, info, req, reply)

		done()

		return err
	}
}
