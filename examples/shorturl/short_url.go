package shorturl

import (
	"context"
	"time"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/ext/client/kredis"
	"github.com/go-kod/kod/interceptor/kvalidate"
	"github.com/google/uuid"
)

type config struct {
	Prefix      string
	RedisConfig kredis.Config
}

type impl struct {
	kod.Implements[Component]
	kod.WithConfig[config]

	redis *kredis.Client
	uuid  uuid.UUID
}

func (i *impl) Init(ctx context.Context) error {
	i.redis = i.Config().RedisConfig.Build()
	i.uuid = uuid.New()

	if i.Config().Prefix == "" {
		i.Config().Prefix = "shorturl:"
	}

	return nil
}

type GenerateRequest struct {
	URL      string `validate:"required"`
	Duration time.Duration
}

type GenerateResponse struct {
	Short string
}

// Generate generates a short url.
func (i *impl) Generate(ctx context.Context, req *GenerateRequest) (*GenerateResponse, error) {

	short := i.uuid.String()
	key := i.Config().Prefix + short
	_, err := i.redis.SetEx(ctx, key, req.URL, req.Duration).Result()
	if err != nil {
		i.L(ctx).Error("failed to set key", "key", key, "req", req, "error", err)
		return nil, err
	}
	return &GenerateResponse{
		Short: short,
	}, nil
}

type GetRequest struct {
	Short string `validate:"required"`
}

type GetResponse struct {
	URL string
}

// Get gets the original url from short url.
func (i *impl) Get(ctx context.Context, req *GetRequest) (*GetResponse, error) {
	key := i.Config().Prefix + req.Short

	result, err := i.redis.Get(ctx, key).Result()
	if err != nil {
		i.L(ctx).Error("failed to get key", "req", req, "error", err)
		return nil, err
	}
	return &GetResponse{
		URL: result,
	}, nil
}

func (i *impl) Interceptors() []kod.Interceptor {
	return []kod.Interceptor{
		kvalidate.Interceptor(),
	}
}
