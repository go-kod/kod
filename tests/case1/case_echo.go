package case1

import (
	"context"
	"errors"

	"github.com/avast/retry-go/v4"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor/kretry"
	"github.com/labstack/echo/v4"
)

type testEchoControllerImpl struct {
	kod.Implements[testEchoController]

	retry int
}

// Hello is a method of testEchoControllerImpl
func (t *testEchoControllerImpl) Hello(c echo.Context) error {
	return c.String(200, "Hello, World!")
}

// Error is a method of testEchoControllerImpl
func (t *testEchoControllerImpl) Error(c echo.Context) error {
	t.retry++

	return errors.New("!!!")
}

func (t *testEchoControllerImpl) Interceptors() []kod.Interceptor {
	return []kod.Interceptor{
		func(ctx context.Context, info kod.CallInfo, req, reply []any, invoker kod.HandleFunc) (err error) {
			return invoker(ctx, info, req, reply)
		},
		kretry.Interceptor(retry.Attempts(2)),
	}
}
