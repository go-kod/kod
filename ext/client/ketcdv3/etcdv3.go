package ketcdv3

import (
	"context"
	"fmt"
	"time"

	"dario.cat/mergo"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

type Config struct {
	Endpoints []string
	Timeout   time.Duration
}

func (r Config) Build(_ context.Context) (*clientv3.Client, error) {
	err := mergo.Merge(&r, Config{
		Endpoints: []string{"localhost:2379"},
		Timeout:   3 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to merge config: %w", err)
	}

	etcd, err := clientv3.New(clientv3.Config{
		Endpoints:   r.Endpoints,
		DialTimeout: r.Timeout,
		DialOptions: []grpc.DialOption{},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create etcd client: %w", err)
	}

	return etcd, nil
}
