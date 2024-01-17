package case1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-kod/kod"
)

func TestHttpHandler(t *testing.T) {
	t.Parallel()
	kod.RunTest(t, func(ctx context.Context, k HTTPController) {
		record := httptest.NewRecorder()

		r, _ := http.NewRequest(http.MethodGet, "/hello/gin", nil)
		// if ctx is not passed, this will panic
		k.Foo(record, r)
	})
}
