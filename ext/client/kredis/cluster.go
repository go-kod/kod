package kredis

import (
	"time"

	"dario.cat/mergo"
	"github.com/redis/go-redis/extra/redisotel/v9"
	redis "github.com/redis/go-redis/v9"
	"github.com/samber/lo"
)

type ClusterClient = redis.ClusterClient

type ClusterConfig struct {
	Addrs        []string
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	Password     string
	DB           int
}

func (c ClusterConfig) Build() *ClusterClient {
	lo.Must0(mergo.Merge(&c, Config{
		DialTimeout:  3 * time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
	}))

	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        c.Addrs,
		Password:     c.Password,
		DialTimeout:  c.DialTimeout,
		ReadTimeout:  c.ReadTimeout,
		WriteTimeout: c.WriteTimeout,
	})

	// Enable tracing instrumentation.
	if err := redisotel.InstrumentTracing(rdb); err != nil {
		panic(err)
	}

	// Enable metrics instrumentation.
	if err := redisotel.InstrumentMetrics(rdb); err != nil {
		panic(err)
	}

	return rdb
}
