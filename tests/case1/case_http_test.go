package case1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-kod/kod"
)

func TestHttpHandler(t *testing.T) {
	t.Parallel()
	kod.RunTest(t, func(ctx context.Context, k HTTPController) {
		record := httptest.NewRecorder()

		r, _ := http.NewRequest(http.MethodGet, "/hello/gin", nil)
		r = r.WithContext(ctx)

		k.Foo(record, r)

		require.Equal(t, http.StatusOK, record.Code)
		require.Equal(t, "Hello, World!", record.Body.String())
	})
}
