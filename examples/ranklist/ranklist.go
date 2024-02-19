package ranklist

import (
	"context"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/ext/client/kredis"
	"github.com/redis/go-redis/v9"
)

type config struct {
	Prefix      string
	RedisConfig kredis.Config
}

type impl struct {
	kod.Implements[Component]
	kod.WithConfig[config]

	redis *kredis.Client
}

func (i *impl) Init(ctx context.Context) error {
	i.redis = i.Config().RedisConfig.Build()

	return nil
}

type AddRequest struct {
	Key    string  `validate:"required"`
	Member string  `validate:"required"`
	Score  float64 `validate:"required"`
}

func (i *impl) Add(ctx context.Context, req *AddRequest) error {
	_, err := i.redis.ZAdd(ctx,
		i.Config().Prefix+req.Key,
		redis.Z{Score: req.Score, Member: req.Member}).Result()
	if err != nil {
		i.L(ctx).Error("failed to add", "req", req, "error", err)
		return err
	}

	return nil
}

type RankListRequest struct {
	Key    string `validate:"required"`
	Min    string `validate:"required"`
	Max    string `validate:"required"`
	Offset int64
	Count  int64
}

func (i *impl) RankList(ctx context.Context, req *RankListRequest) ([]string, error) {
	return i.redis.ZRevRangeByScore(ctx, i.Config().Prefix+req.Key, &redis.ZRangeBy{
		Min:    req.Min,
		Max:    req.Max,
		Offset: req.Offset,
		Count:  req.Count,
	}).Result()
}
