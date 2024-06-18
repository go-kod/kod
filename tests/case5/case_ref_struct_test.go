package case5

import (
	"context"
	"testing"

	"github.com/go-kod/kod"
	"github.com/stretchr/testify/require"
)

func TestRefStruct(t *testing.T) {
	err := kod.Run(context.Background(), func(ctx context.Context, comp *refStructImpl) error {
		return nil
	})
	require.EqualError(t, err, "component implementation struct case5.refStructImpl has field kod.Ref[github.com/go-kod/kod/tests/case5.testRefStruct1], but field type case5.testRefStruct1 is not an interface")
}
