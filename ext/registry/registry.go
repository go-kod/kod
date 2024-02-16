package registry

import (
	"context"

	"google.golang.org/grpc/resolver"
)

// Registry is the interface that wraps the basic methods of a registry.
type Registry interface {
	// Register registers the provided info to the registry.
	Register(ctx context.Context, info ServiceInfo) error
	// UnRegister unregisters the provided info from the registry.
	UnRegister(ctx context.Context, info ServiceInfo) error
	// ResolveBuilder returns a resolver.Builder that will be used to create a resolver for the provided target.
	ResolveBuilder(ctx context.Context) (resolver.Builder, error)
}

// ServiceInfo is the interface that wraps the basic methods of a service info.
type ServiceInfo struct {
	// Scheme is the scheme of the service.
	Scheme string
	// Addr is the address of the service.
	Addr string
	// Metadata is the metadata of the service.
	Metadata map[string]interface{}
}
