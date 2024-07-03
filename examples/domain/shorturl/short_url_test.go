package shorturl

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/go-kod/kod"
)

func TestShortURL(t *testing.T) {
	kod.RunTest(t, func(ctx context.Context, impl Component) {
		result, err := impl.Generate(ctx, &GenerateRequest{
			URL:      "https://github.com/go-kod/kod",
			Duration: time.Second,
		})
		assert.Nil(t, err)
		assert.NotEmpty(t, result.Short)

		getRes, err := impl.Get(ctx, &GetRequest{
			Short: result.Short,
		})
		assert.Nil(t, err)
		assert.Equal(t, "https://github.com/go-kod/kod", getRes.URL)
	}, kod.WithConfigFile("./kod.toml"))
}

func TestShortURLInvalid(t *testing.T) {
	t.Run("Generate", func(t *testing.T) {
		t.Parallel()

		kod.RunTest(t, func(ctx context.Context, impl Component) {
			result, err := impl.Generate(ctx, &GenerateRequest{
				Duration: time.Second,
			})
			assert.NotNil(t, err)
			assert.Nil(t, result)
		})
	})

	t.Run("Get", func(t *testing.T) {
		t.Parallel()

		kod.RunTest(t, func(ctx context.Context, impl Component) {
			result, err := impl.Get(ctx, &GetRequest{})
			assert.NotNil(t, err)
			assert.Nil(t, result)
		})
	})
}
