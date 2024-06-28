package kgin

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/go-kod/kod/ext/internal/knet"
	"github.com/go-kod/kod/ext/registry"
)

type Config struct {
	Address string
}

func (c Config) Build() *Server {
	s := gin.Default()
	s.Use(otelgin.Middleware(""))

	return &Server{Engine: s, http: &http.Server{Handler: s}, c: c}
}

func (s *Server) WithRegistry(r registry.Registry) *Server {
	s.registry = r
	return s
}

type (
	Engine      = gin.Engine
	Context     = gin.Context
	H           = gin.H
	RouterGroup = gin.RouterGroup

	Server struct {
		http *http.Server
		c    Config
		*gin.Engine
		registry registry.Registry
		lis      net.Listener
	}
)

func (s *Server) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.c.Address)
	if err != nil {
		return err
	}

	s.lis = lis

	if s.registry != nil {
		time.AfterFunc(time.Second, func() {
			err := s.registry.Register(ctx, registry.ServiceInfo{
				Scheme:   s.Scheme(),
				Addr:     lo.Must(knet.ExtractAddress(s.c.Address, s.lis)),
				Metadata: nil,
			})
			if err != nil {
				panic(err)
			}
		})
	}

	slog.Info("gin server started on: " + s.lis.Addr().String())
	return s.http.Serve(lis)
}

func (s *Server) GracefulStop(ctx context.Context) error {
	if s.registry != nil {
		err := s.registry.UnRegister(ctx, registry.ServiceInfo{
			Scheme:   s.Scheme(),
			Addr:     lo.Must(knet.ExtractAddress(s.c.Address, s.lis)),
			Metadata: nil,
		})
		if err != nil {
			return err
		}
	}

	slog.Info("grpc server stopped")
	return s.http.Shutdown(ctx)
}

func (s *Server) Scheme() string {
	return "http"
}
