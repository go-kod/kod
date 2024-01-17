package case1

import (
	"context"

	"github.com/go-kod/kod"
	"github.com/labstack/echo/v4"
)

type testEchoControllerImpl struct {
	kod.Implements[testEchoController]
}

// Hello is a method of testEchoControllerImpl
func (t *testEchoControllerImpl) Hello(c echo.Context) error {
	return c.String(200, "Hello, World!")
}

func (t *testEchoControllerImpl) Interceptors() []kod.Interceptor {
	return []kod.Interceptor{
		func(ctx context.Context, info kod.CallInfo, req, reply []any, invoker kod.HandleFunc) (err error) {
			return invoker(ctx, info, req, reply)
		},
	}
}
