package case1

import "github.com/gin-gonic/gin"

// Foo2 ...
func (t *test1ControllerImpl) Foo2(cccccc *gin.Context) {
	_, _ = t.test1Component.Get().Foo(cccccc, &FooReq{})
}
