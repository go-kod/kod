package kgrpc

import (
	"context"
	"log/slog"
	"net"
	"time"

	"github.com/bufbuild/protovalidate-go"
	"github.com/go-kod/kod/ext/internal/knet"
	"github.com/go-kod/kod/ext/registry"
	protovalidate_middleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/samber/lo"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

var (
	Method                  = grpc.Method
	ChainUnaryInterceptor   = grpc.ChainUnaryInterceptor
	ChainStreamInterceptor  = grpc.ChainStreamInterceptor
	UnaryServerInterceptor  grpc.UnaryServerInterceptor
	StreamServerInterceptor grpc.StreamServerInterceptor
)

type Config struct {
	Address string
}

type Server struct {
	*grpc.Server
	Config
	registry registry.Registry
	lis      net.Listener
}

func (c Config) Build(opts ...grpc.ServerOption) *Server {
	mds := make([]protoreflect.MessageDescriptor, 0)
	protoregistry.GlobalTypes.RangeMessages(func(mt protoreflect.MessageType) bool {
		mds = append(mds, mt.Descriptor())
		return true
	})

	validator, err := protovalidate.New(
		protovalidate.WithDescriptors(mds...),
	)
	if err != nil {
		panic(err)
	}

	defaultOpts := []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(),
			protovalidate_middleware.UnaryServerInterceptor(validator),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(),
			protovalidate_middleware.StreamServerInterceptor(validator),
		),
	}

	defaultOpts = append(defaultOpts, opts...)

	s := grpc.NewServer(defaultOpts...)

	reflection.Register(s)

	return &Server{Server: s, Config: c}
}

func (s *Server) WithRegistry(r registry.Registry) *Server {
	s.registry = r
	return s
}

func (s *Server) Run(ctx context.Context) error {
	lis, err := net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}

	s.lis = lis

	if s.registry != nil {
		time.AfterFunc(time.Second, func() {
			err := s.registry.Register(ctx, registry.ServiceInfo{
				Scheme:   s.Scheme(),
				Addr:     lo.Must(knet.ExtractAddress(s.Config.Address, s.lis)),
				Metadata: nil,
			})
			if err != nil {
				panic(err)
			}
		})
	}

	slog.Info("grpc server started on: " + lis.Addr().String())
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
			panic(err)
		}
	}
	slog.Info("grpc server stopped")
	s.Server.GracefulStop()
	return nil
}

func (s *Server) Scheme() string {
	return "grpc"
}
