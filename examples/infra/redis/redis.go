package redis

import (
	"context"

	goredis "github.com/redis/go-redis/v9"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/ext/client/kredis"
)

type impl struct {
	kod.Implements[Component]
	kod.WithConfig[config]

	cc *goredis.Client
}

type config struct {
	RedisConfig kredis.Config
}

func (ins *impl) Init(ctx context.Context) error {
	ins.cc = ins.Config().RedisConfig.Build()
	return nil
}

func (ins *impl) Client() *goredis.Client {
	return ins.cc
}
