package case1

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/go-kod/kod"
)

type testGinControllerImpl struct {
	kod.Implements[testGinController]
}

// Hello is a method of testGinControllerImpl
func (t *testGinControllerImpl) Hello(c *gin.Context) {
	c.String(200, "Hello, World!")
}

func (t *testGinControllerImpl) Interceptors() []kod.Interceptor {
	return []kod.Interceptor{
		func(ctx context.Context, info kod.CallInfo, req, reply []any, invoker kod.HandleFunc) (err error) {
			return invoker(ctx, info, req, reply)
		},
	}
}
