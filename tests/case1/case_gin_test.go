package case1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"

	"github.com/go-kod/kod"
)

func Test_testGinControllerImpl_Hello(t *testing.T) {
	kod.RunTest(t, func(ctx context.Context, controller testGinController) {
		rec := httptest.NewRecorder()
		c := gin.CreateTestContextOnly(rec, gin.New())
		c.Request, _ = http.NewRequest(http.MethodGet, "/hello/gin", nil)
		controller.Hello(c)
		require.Equal(t, http.StatusOK, rec.Code)
	})
}
