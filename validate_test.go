package kod

import (
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/go-kod/kod/internal/reflects"
	"github.com/go-kod/kod/internal/registry"
)

func TestValidateUnregisteredRef(t *testing.T) {
	type foo interface{}
	type fooImpl struct{ Ref[io.Reader] }
	regs := []*registry.Registration{
		{
			Name:  "foo",
			Iface: reflects.TypeFor[foo](),
			Impl:  reflects.TypeFor[fooImpl](),
		},
	}
	err := validateRegistrations(regs)
	if err == nil {
		t.Fatal("unexpected validateRegistrations success")
	}
	const want = "component io.Reader was not registered"
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("validateRegistrations: got %q, want %q", err, want)
	}
}

// TestValidateNoRegistrations tests that validateRegistrations succeeds on an
// empty set of registrations.
func TestValidateNoRegistrations(t *testing.T) {
	if err := validateRegistrations(nil); err != nil {
		t.Fatal(err)
	}
}

func TestMultipleRegistrations(t *testing.T) {
	type foo interface{}
	type fooImpl struct{ Ref[io.Reader] }
	regs := []*Registration{
		{
			Name:  "github.com/go-kod/kod/Main",
			Iface: reflect.TypeOf((*Main)(nil)).Elem(),
			Impl:  reflect.TypeOf(fooImpl{}),
			Refs:  `⟦48699770:KoDeDgE:github.com/go-kod/kod/Main→github.com/go-kod/kod/tests/graphcase/test1Controller⟧`,
		},
		{
			Name:  "github.com/go-kod/kod/Main",
			Iface: reflect.TypeOf((*foo)(nil)).Elem(),
			Impl:  reflect.TypeOf(fooImpl{}),
			Refs:  `⟦48699770:KoDeDgE:github.com/go-kod/kod/tests/graphcase/test1Controller→github.com/go-kod/kod/Main⟧`,
		},
	}
	err := checkCircularDependency(regs)
	if err == nil {
		t.Fatal("unexpected checkCircularDependency success")
	}
	const want = "components [github.com/go-kod/kod/Main], error vertex already exists"
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("checkCircularDependency: got %q, want %q", err, want)
	}
}

func TestCycleRegistrations(t *testing.T) {
	type test1Controller interface{}
	type test1ControllerImpl struct{ Ref[io.Reader] }
	type mainImpl struct{ Ref[test1Controller] }
	regs := []*Registration{
		{
			Name:  "github.com/go-kod/kod/Main",
			Iface: reflect.TypeOf((*Main)(nil)).Elem(),
			Impl:  reflect.TypeOf(mainImpl{}),
			Refs:  `⟦48699770:KoDeDgE:github.com/go-kod/kod/Main→github.com/go-kod/kod/test1Controller⟧`,
		},
		{
			Name:  "github.com/go-kod/kod/test1Controller",
			Iface: reflect.TypeOf((*test1Controller)(nil)).Elem(),
			Impl:  reflect.TypeOf(test1ControllerImpl{}),
			Refs:  `⟦b8422d0e:KoDeDgE:github.com/go-kod/kod/test1Controller→github.com/go-kod/kod/Main⟧`,
		},
	}
	err := checkCircularDependency(regs)
	if err == nil {
		t.Fatal("unexpected checkCircularDependency success")
	}
	const want = "components [github.com/go-kod/kod/test1Controller] and [github.com/go-kod/kod/Main] have cycle Ref"
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("checkCircularDependency: got %q, want %q", err, want)
	}
}
