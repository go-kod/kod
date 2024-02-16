package kgrpc

import (
	"context"
	"time"

	"dario.cat/mergo"
	"github.com/go-kod/kod/ext/registry"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"github.com/samber/lo"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClientConn = grpc.ClientConn

type Config struct {
	Target  string
	Timeout time.Duration

	registry registry.Registry
}

func (c Config) WithRegistry(r registry.Registry) Config {
	c.registry = r
	return c
}

func (c Config) Build(opts ...grpc.DialOption) *ClientConn {
	lo.Must0(mergo.Merge(&c, Config{
		Timeout: 3 * time.Second,
	}))

	ctx := context.Background()

	defaultOpts := []grpc.DialOption{
		grpc.WithNoProxy(),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			retry.UnaryClientInterceptor(),
			timeout.UnaryClientInterceptor(c.Timeout),
		),
		grpc.WithChainStreamInterceptor(
			retry.StreamClientInterceptor(),
		),
	}

	if c.registry != nil {
		builder, err := c.registry.ResolveBuilder(ctx)
		if err != nil {
			panic(err)
		}
		defaultOpts = append(defaultOpts, grpc.WithResolvers(builder))
	}

	defaultOpts = append(defaultOpts, opts...)

	cc, err := grpc.DialContext(ctx, c.Target, defaultOpts...)
	if err != nil {
		panic(err)
	}

	return cc
}
