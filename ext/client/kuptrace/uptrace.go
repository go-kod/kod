package kuptrace

import (
	"context"
	"fmt"

	"dario.cat/mergo"
	"github.com/go-kod/kod"
	"github.com/uptrace/uptrace-go/uptrace"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/contrib/samplers/probability/consistent"
)

type Config struct {
	DSN      string
	Debug    bool
	Fraction float64
}

type Client struct {
}

func (c Config) Build(ctx context.Context) (*Client, error) {
	err := mergo.Merge(&c, Config{
		Fraction: 1,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to merge config: %w", err)
	}

	err = host.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start host instrumentation: %w", err)
	}
	err = runtime.Start()
	if err != nil {
		return nil, fmt.Errorf("failed to start runtime instrumentation: %w", err)
	}

	kodConfig := kod.FromContext(ctx).Config()

	sampler := consistent.ParentProbabilityBased(
		consistent.ProbabilityBased(c.Fraction),
	)

	opts := []uptrace.Option{
		uptrace.WithDSN(c.DSN),
		uptrace.WithServiceName(kodConfig.Name),
		uptrace.WithServiceVersion(kodConfig.Version),
		uptrace.WithDeploymentEnvironment(kodConfig.Env),
		uptrace.WithTraceSampler(sampler),
	}

	if c.Debug {
		opts = append(opts, uptrace.WithPrettyPrintSpanExporter())
	}

	// Configure OpenTelemetry with sensible defaults.
	uptrace.ConfigureOpentelemetry(opts...)

	return &Client{}, nil
}

func (c *Client) Stop(ctx context.Context) error {
	return uptrace.ForceFlush(ctx)
}
