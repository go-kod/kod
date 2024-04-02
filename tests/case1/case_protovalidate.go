package case1

import (
	"context"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/interceptor/kprotovalidate"
	"github.com/go-kod/kod/tests/proto/examplev1"
)

type protoValidateComponent struct {
	kod.Implements[ProtoValidateComponent]
}

func (ins *protoValidateComponent) Interceptors() []interceptor.Interceptor {
	return []interceptor.Interceptor{
		kprotovalidate.Interceptor(),
	}
}

func (ins *protoValidateComponent) Validate(ctx context.Context, req *examplev1.Person) error {
	return nil
}
