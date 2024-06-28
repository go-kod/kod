package case1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"

	"github.com/go-kod/kod"
)

func Test_testEchoControllerImpl_Hello(t *testing.T) {
	kod.RunTest(t, func(ctx context.Context, controller testEchoController) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		require.Nil(t, controller.Hello(c))
		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "Hello, World!", rec.Body.String())
	})
}

func Test_testEchoControllerImpl_Panic(t *testing.T) {
	kod.RunTest(t, func(ctx context.Context, controller testEchoController) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		controller.Error(c)
		require.Equal(t, controller.(testEchoController_local_stub).impl.(*testEchoControllerImpl).retry, 1)
		require.Equal(t, "", rec.Body.String())
	})
}
