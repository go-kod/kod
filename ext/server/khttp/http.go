package khttp

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/samber/lo"

	"github.com/go-kod/kod/ext/internal/knet"
	"github.com/go-kod/kod/ext/registry"
)

type Config struct {
	Address string
}

func (c Config) Build() *Server {
	s := &http.Server{
		Addr: c.Address,
	}

	return &Server{Server: s, Config: c}
}

func (s *Server) WithRegistry(r registry.Registry) *Server {
	s.registry = r
	return s
}

type (
	Server struct {
		*http.Server
		Config
		registry registry.Registry
		lis      net.Listener
	}
)

func (s *Server) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}

	s.lis = lis

	if s.registry != nil {
		time.AfterFunc(time.Second, func() {
			err = s.registry.Register(ctx, registry.ServiceInfo{
				Scheme:   s.Scheme(),
				Addr:     lo.Must(knet.ExtractAddress(s.Config.Address, s.lis)),
				Metadata: nil,
			})
			if err != nil {
				panic(err)
			}
		})
	}

	slog.Info("http server started on: " + lis.Addr().String())
	return s.Server.Serve(lis)
}

func (s *Server) GracefulStop(ctx context.Context) error {
	if s.registry != nil {
		err := s.registry.UnRegister(ctx, registry.ServiceInfo{
			Scheme:   s.Scheme(),
			Addr:     lo.Must(knet.ExtractAddress(s.Config.Address, s.lis)),
			Metadata: nil,
		})
		if err != nil {
			return err
		}
	}
	slog.Info("http server stopped")
	return s.Server.Shutdown(ctx)
}

func (s *Server) Scheme() string {
	return "http"
}
