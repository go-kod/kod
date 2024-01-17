package case1

import (
	"net/http"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor/kaccesslog"
	"github.com/go-kod/kod/interceptor/krecovery"
	"github.com/go-kod/kod/interceptor/ktimeout"
)

type httpControllerImpl struct {
	kod.Implements[HTTPController]
}

// Foo is a http handler
func (t *httpControllerImpl) Foo(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Deadline()
	if !ok {
		panic("no deadline")
	}
}

func (t *httpControllerImpl) Interceptors() []kod.Interceptor {
	return []kod.Interceptor{
		krecovery.Interceptor(),
		kaccesslog.Interceptor(),
		ktimeout.Interceptor(),
	}
}
