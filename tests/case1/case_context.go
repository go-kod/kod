package case1

import (
	"context"
	"time"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/interceptor/kaccesslog"
	"github.com/go-kod/kod/interceptor/kcircuitbreaker"
	"github.com/go-kod/kod/interceptor/kmetric"
	"github.com/go-kod/kod/interceptor/kratelimit"
	"github.com/go-kod/kod/interceptor/krecovery"
	"github.com/go-kod/kod/interceptor/ktimeout"
	"github.com/go-kod/kod/interceptor/ktrace"
)

type ctxImpl struct {
	kod.Implements[ctxInterface]
}

// Foo is a http handler
func (t *ctxImpl) Foo(ctx context.Context) {
	_, ok := ctx.Deadline()
	if !ok {
		panic("no deadline")
	}
}

func (t *ctxImpl) Interceptors() []interceptor.Interceptor {
	return []interceptor.Interceptor{
		krecovery.Interceptor(),
		kaccesslog.Interceptor(),
		ktimeout.Interceptor(ktimeout.WithTimeout(time.Second)),
		kmetric.Interceptor(),
		ktrace.Interceptor(),
		kcircuitbreaker.Interceptor(),
		kratelimit.Interceptor(),
	}
}
