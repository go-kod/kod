package redis

import (
	"context"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/ext/client/kredis"
	goredis "github.com/redis/go-redis/v9"
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
