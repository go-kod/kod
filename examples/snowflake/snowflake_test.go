package snowflake

import (
	"context"
	"testing"

	"github.com/go-kod/kod"
	"github.com/stretchr/testify/assert"
)

func TestSnowflake(t *testing.T) {
	kod.RunTest(t, func(ctx context.Context, c Component) {
		nextIDRes, err := c.NextID(ctx, &NextIDRequest{})
		assert.Nil(t, err)
		assert.NotZero(t, nextIDRes.ID)
	})
}
