package case1

import (
	"context"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/interceptor/krecovery"
)

type panicCase struct {
	kod.Implements[panicCaseInterface]
}

func (t *panicCase) TestPanic(ctx context.Context) {
	panic("panic")
}

func (t *panicCase) Interceptors() []interceptor.Interceptor {
	return []interceptor.Interceptor{
		krecovery.Interceptor(),
	}
}

type panicNoRecvoeryCase struct {
	kod.Implements[panicNoRecvoeryCaseInterface]
}

func (t *panicNoRecvoeryCase) TestPanic(ctx context.Context) {
	panic("panic")
}
