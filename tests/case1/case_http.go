package case1

import (
	"net/http"

	"github.com/go-kod/kod"
)

type httpControllerImpl struct {
	kod.Implements[HTTPController]
}

// Foo is a http handler
func (t *httpControllerImpl) Foo(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello, World!"))
}
