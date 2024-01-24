package registry

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/dominikbraun/graph"
	"github.com/go-kod/kod/internal/callgraph"
	"github.com/go-kod/kod/internal/reflects"
)

// CheckCircularDependency checks that there are no circular dependencies
// between registered components.
func CheckCircularDependency(reg []*Registration) error {
	g := graph.New(graph.StringHash, graph.Directed(), graph.PreventCycles())

	for _, reg := range reg {
		if err := g.AddVertex(reg.Name); err != nil {
			return fmt.Errorf("components [%s], error %s", reg.Name, err)
		}
	}

	var errs []error
	for _, reg := range reg {
		edges := callgraph.ParseEdges([]byte(reg.Refs))
		for _, edge := range edges {
			err := g.AddEdge(edge[0], edge[1])
			if err != nil {
				switch err {
				case graph.ErrEdgeAlreadyExists, graph.ErrEdgeCreatesCycle:
					err = fmt.Errorf("components [%s] and [%s] have cycle Ref", edge[0], edge[1])
				}
				errs = append(errs, err)
			}
		}
	}

	return errors.Join(errs...)
}

// ValidateRegistrations checks that all registered component interfaces are
// implemented by a registered component implementation struct.
func ValidateRegistrations(reg []*Registration) error {
	// Gather the set of registered interfaces.
	intfs := map[reflect.Type]struct{}{}
	for _, reg := range reg {
		intfs[reg.Iface] = struct{}{}
	}

	// Check that for every kod.Ref[T] field in a component implementation
	// struct, T is a registered interface.
	var errs []error
	for _, reg := range reg {
		for i := 0; i < reg.Impl.NumField(); i++ {
			f := reg.Impl.Field(i)
			switch {
			case f.Type.Implements(reflects.TypeFor[interface{ isRef() }]()):
				// f is a kod.Ref[T].
				v := f.Type.Field(0) // a Ref[T]'s value field
				if _, ok := intfs[v.Type]; !ok {
					// T is not a registered component interface.
					err := fmt.Errorf(
						"component implementation struct %v has field %v, but component %v was not registered; maybe you forgot to run 'kod generate'",
						reg.Impl, f.Type, v.Type,
					)
					errs = append(errs, err)
				}
			}
		}
	}
	return errors.Join(errs...)
}
