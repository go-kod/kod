package case1

import (
	"context"
	"errors"

	"github.com/avast/retry-go/v4"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/interceptor/kretry"
)

type interceptorRetry struct {
	kod.Implements[InterceptorRetry]
}

func (t *interceptorRetry) Init(ctx context.Context) error {
	t.L(ctx).Info("interceptorRetry init...")

	return nil
}

func (t *interceptorRetry) TestError(ctx context.Context) error {
	return errors.New("test error")
}

func (t *interceptorRetry) TestNormal(ctx context.Context) error {
	return nil
}

func (t *interceptorRetry) Shutdown(ctx context.Context) error {
	t.L(ctx).Info("interceptorRetry shutdown...")

	return nil
}

func (t *interceptorRetry) Interceptors() []interceptor.Interceptor {
	return []interceptor.Interceptor{
		kretry.Interceptor(
			retry.Attempts(2),
		),
	}
}
