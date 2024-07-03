package case1

import (
	"context"

	"github.com/gin-gonic/gin"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
)

type testGinControllerImpl struct {
	kod.Implements[testGinController]
}

// Hello is a method of testGinControllerImpl
func (t *testGinControllerImpl) Hello(c *gin.Context) {
	c.String(200, "Hello, World!")
}

func (t *testGinControllerImpl) Interceptors() []interceptor.Interceptor {
	return []interceptor.Interceptor{
		func(ctx context.Context, info interceptor.CallInfo, req, reply []any, invoker interceptor.HandleFunc) (err error) {
			return invoker(ctx, info, req, reply)
		},
	}
}
