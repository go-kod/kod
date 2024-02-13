package case1

import (
	"net/http"
	"time"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor/kaccesslog"
	"github.com/go-kod/kod/interceptor/kcircuitbreaker"
	"github.com/go-kod/kod/interceptor/kmetric"
	"github.com/go-kod/kod/interceptor/kratelimit"
	"github.com/go-kod/kod/interceptor/krecovery"
	"github.com/go-kod/kod/interceptor/ktimeout"
	"github.com/go-kod/kod/interceptor/ktrace"
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
	w.Write([]byte("Hello, World!"))
}

func (t *httpControllerImpl) Interceptors() []kod.Interceptor {
	return []kod.Interceptor{
		krecovery.Interceptor(),
		kaccesslog.Interceptor(),
		ktimeout.Interceptor(ktimeout.WithTimeout(time.Second)),
		kmetric.Interceptor(),
		ktrace.Interceptor(),
		kcircuitbreaker.Interceptor(),
		kratelimit.Interceptor(),
	}
}
