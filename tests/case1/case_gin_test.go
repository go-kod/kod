package case1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-kod/kod"
	"github.com/stretchr/testify/assert"
)

func Test_testGinControllerImpl_Hello(t *testing.T) {
	kod.RunTest(t, func(ctx context.Context, controller testGinController) {
		rec := httptest.NewRecorder()
		c := gin.CreateTestContextOnly(rec, gin.New())
		c.Request, _ = http.NewRequest(http.MethodGet, "/hello/gin", nil)
		controller.Hello(c)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}
