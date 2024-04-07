package case1

import (
	"context"
	"testing"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/tests/proto/examplev1"
	"github.com/stretchr/testify/require"
)

func TestInterfacProtoValidate(t *testing.T) {
	t.Parallel()
	kod.RunTest(t, func(ctx context.Context, k ProtoValidateComponent) {
		var err error
		err = k.Validate(ctx, &examplev1.Person{})
		require.Contains(t, err.Error(), "proto validate failed: validation error:\n - id: value must be greater than 999 [uint64.gt]")

		err = k.Validate(ctx, &examplev1.Person{
			Id:    1000,
			Email: "hnlq.sysu@gmail.com",
			Name:  "name",
		})
		require.Nil(t, err)
	})
}
