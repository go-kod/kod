package case1

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-kod/kod"
	"github.com/stretchr/testify/assert"
)

func TestHttpHandler(t *testing.T) {
	t.Parallel()
	kod.RunTest(t, func(ctx context.Context, k HTTPController) {
		record := httptest.NewRecorder()

		r, _ := http.NewRequest(http.MethodGet, "/hello/gin", nil)
		ctx = StartTrace(ctx)
		r = r.WithContext(ctx)

		k.Foo(record, r)

		assert.Equal(t, http.StatusOK, record.Code)
		assert.Equal(t, "Hello, World!", record.Body.String())
	})
}
