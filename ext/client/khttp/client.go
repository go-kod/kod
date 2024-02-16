package http

import (
	"context"
	"net/http"
	"net/http/httptrace"
	"time"

	"dario.cat/mergo"
	"github.com/samber/lo"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type Client = http.Client

type ClientConfig struct {
	Address string
	Timeout time.Duration
}

func (c ClientConfig) Build() *Client {

	lo.Must0(mergo.Merge(&c, ClientConfig{
		Timeout: 3 * time.Second,
	}))

	defaultTransport := http.DefaultTransport.(*http.Transport).Clone()
	defaultTransport.Proxy = nil

	return &http.Client{
		Transport: otelhttp.NewTransport(defaultTransport,
			otelhttp.WithClientTrace(func(ctx context.Context) *httptrace.ClientTrace {
				return otelhttptrace.NewClientTrace(ctx)
			})),
		Timeout: c.Timeout,
	}
}
