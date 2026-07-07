package kod

import (
	"context"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/go-kod/kod/internal/registry"
)

func TestFill(t *testing.T) {
	t.Run("case 2", func(t *testing.T) {
		assert.NotNil(t, fillRefs(nil, nil, nil))
	})

	t.Run("case 3", func(t *testing.T) {
		i := 0
		assert.NotNil(t, fillRefs(&i, nil, nil))
	})
}

func TestRefType(t *testing.T) {
	var ref Ref[io.Reader]
	assert.Equal(t, reflect.TypeFor[io.Reader](), ref.refType())
}

func TestValidateUnregisteredRef(t *testing.T) {
	type foo interface{}
	type fooImpl struct{ Ref[io.Reader] }
	regs := []*registry.Registration{
		{
			Name:      "foo",
			Interface: reflect.TypeFor[foo](),
			Impl:      reflect.TypeFor[fooImpl](),
		},
	}
	_, err := processRegistrations(regs)
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
	if _, err := processRegistrations(nil); err != nil {
		t.Fatal(err)
	}
}

func TestValidateInvalidRegistration(t *testing.T) {
	regs := []*Registration{
		nil,
		{Name: "missing-interface", Impl: reflect.TypeFor[struct{}]()},
		{Name: "missing-impl", Interface: reflect.TypeFor[interface{}]()},
	}

	var err error
	assert.NotPanics(t, func() {
		_, err = processRegistrations(regs)
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "component registration is nil")
	assert.Contains(t, err.Error(), `component registration "missing-interface" interface type is nil`)
	assert.Contains(t, err.Error(), `component registration "missing-impl" implementation type is nil`)
}

func TestNewKodInvalidRegistration(t *testing.T) {
	var err error
	assert.NotPanics(t, func() {
		_, err = newKod(context.Background(), WithRegistrations(nil))
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "component registration is nil")
}

func TestValidateImplementationMismatch(t *testing.T) {
	type foo interface{ Foo() }
	type fooImpl struct{}
	regs := []*Registration{
		{
			Name:      "foo",
			Interface: reflect.TypeFor[foo](),
			Impl:      reflect.TypeFor[fooImpl](),
		},
	}
	_, err := processRegistrations(regs)
	if err == nil {
		t.Fatal("unexpected validateRegistrations success")
	}
	const want = "does not implement interface"
	if !strings.Contains(err.Error(), want) {
		t.Fatalf("validateRegistrations: got %q, want %q", err, want)
	}
}

func TestMultipleRegistrations(t *testing.T) {
	type foo interface{}
	type fooImpl struct{ Ref[io.Reader] }
	regs := []*Registration{
		{
			Name:      "github.com/go-kod/kod/Main",
			Interface: reflect.TypeFor[Main](),
			Impl:      reflect.TypeFor[fooImpl](),
			Refs:      `⟦48699770:KoDeDgE:github.com/go-kod/kod/Main→github.com/go-kod/kod/tests/graphcase/test1Controller⟧`,
		},
		{
			Name:      "github.com/go-kod/kod/Main",
			Interface: reflect.TypeFor[foo](),
			Impl:      reflect.TypeFor[fooImpl](),
			Refs:      `⟦48699770:KoDeDgE:github.com/go-kod/kod/tests/graphcase/test1Controller→github.com/go-kod/kod/Main⟧`,
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
			Name:      "github.com/go-kod/kod/Main",
			Interface: reflect.TypeFor[Main](),
			Impl:      reflect.TypeFor[mainImpl](),
			Refs:      `⟦48699770:KoDeDgE:github.com/go-kod/kod/Main→github.com/go-kod/kod/test1Controller⟧`,
		},
		{
			Name:      "github.com/go-kod/kod/test1Controller",
			Interface: reflect.TypeFor[test1Controller](),
			Impl:      reflect.TypeFor[test1ControllerImpl](),
			Refs:      `⟦b8422d0e:KoDeDgE:github.com/go-kod/kod/test1Controller→github.com/go-kod/kod/Main⟧`,
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

func TestGetImpl(t *testing.T) {
	k, err := newKod(context.Background())
	require.NoError(t, err)

	_, err = k.getImpl(context.Background(), reflect.TypeOf(struct{}{}))
	assert.Error(t, err) // Should fail for unregistered type
}

func TestGetIntf(t *testing.T) {
	k, err := newKod(context.Background())
	require.NoError(t, err)

	_, err = k.getIntf(context.Background(), reflect.TypeFor[interface{}]())
	assert.Error(t, err) // Should fail for unregistered interface
}

func TestGetInterfaceWithoutLocalStub(t *testing.T) {
	type mainImpl struct {
		Implements[Main]
	}
	k, err := newKod(context.Background(), WithRegistrations(&Registration{
		Name:      "main",
		Interface: reflect.TypeFor[Main](),
		Impl:      reflect.TypeFor[mainImpl](),
	}))
	require.NoError(t, err)

	var got Main
	assert.NotPanics(t, func() {
		got, err = k.Get[Main](context.Background())
	})
	require.NoError(t, err)
	assert.IsType(t, &mainImpl{}, got)
}
