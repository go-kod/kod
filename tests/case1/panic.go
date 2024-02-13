package case1

import (
	"net/http"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor/krecovery"
)

type panicCase struct {
	kod.Implements[panicCaseInterface]
}

func (t *panicCase) TestPanic(r *http.Request) {
	panic("panic")
}

func (t *panicCase) Interceptors() []kod.Interceptor {
	return []kod.Interceptor{
		krecovery.Interceptor(),
	}
}
