package kod

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/samber/lo"
)

type fakeComponent struct {
	intf reflect.Type
	impl any
}

func Fake[T any](impl any) fakeComponent {
	t := rtype[T]()
	if _, ok := impl.(T); !ok {
		panic(fmt.Sprintf("%T does not implement %v", impl, t))
	}
	return fakeComponent{intf: t, impl: impl}
}

type runner struct {
	options options
}

func RunTest[T any](t testing.TB, body func(context.Context, T), opts ...func(*options)) {
	options := &options{}
	for _, o := range opts {
		o(options)
	}

	err := runner{options: *options}.sub(t, body)
	if err != nil {
		t.Logf("RunTest failed: %v", err)
		t.FailNow()
	}
}

func (r runner) sub(t testing.TB, testBody any) error {
	ctx, cancelFn := context.WithCancel(context.Background())
	defer func() {
		// Cancel the context so background activity will stop.
		cancelFn()
	}()

	runner, err := newKod(getRegs(), r.options)
	if err != nil {
		return fmt.Errorf("newKod: %v", err)
	}
	defer runner.close(ctx)

	ctx = newContext(ctx, runner)

	t.Helper()
	body, intfs, err := checkRunFunc(ctx, t, testBody)
	if err != nil {
		return fmt.Errorf("kod.Run argument: %v", err)
	}

	// Assume a component Foo implementing struct foo. We disallow tests
	// like the one below where the user provides a fake and a component
	// implementation pointer for the same component.
	//
	//     runner.Fakes = append(runner.Fakes, kod.Fake[Foo](...))
	//     runner.Test(t, func(t *testing.T, f *foo) {...})
	for _, intf := range intfs {
		if _, ok := r.options.fakes[intf]; ok {
			return fmt.Errorf("Component %v has both fake and component implementation pointer", intf)
		}
	}

	if err := body(ctx, runner); err != nil {
		return fmt.Errorf("kod.Run body: %v", err)
	}

	return nil
}

func checkRunFunc(ctx context.Context, t testing.TB, fn any) (func(context.Context, *Kod) error, []reflect.Type, error) {
	fnType := reflect.TypeOf(fn)
	if fnType == nil || fnType.Kind() != reflect.Func {
		return nil, nil, fmt.Errorf("not a func")
	}
	if fnType.IsVariadic() {
		return nil, nil, fmt.Errorf("must not be variadic")
	}
	n := fnType.NumIn()
	if n < 2 {
		return nil, nil, fmt.Errorf("must have at least two args")
	}
	if fnType.NumOut() > 0 {
		return nil, nil, fmt.Errorf("must have no return outputs")
	}
	if fnType.In(0) != reflect.TypeOf(&ctx).Elem() {
		return nil, nil, fmt.Errorf("function first argument type %v does not match first kod.Run argument %v", fnType.In(0), reflect.TypeOf(&ctx).Elem())
	}
	var intfs []reflect.Type
	for i := 1; i < n; i++ {
		switch fnType.In(i).Kind() {
		case reflect.Interface:
			// Do nothing.
		case reflect.Pointer:
			intf, err := extractComponentInterfaceType(fnType.In(i).Elem())
			if err != nil {
				return nil, nil, err
			}
			intfs = append(intfs, intf)
		default:
			return nil, nil, fmt.Errorf("function argument %d type %v must be a component interface or pointer to component implementation", i, fnType.In(i))
		}
	}

	return func(ctx context.Context, runner *Kod) error {
		args := make([]any, n)
		args[0] = ctx
		for i := 1; i < n; i++ {
			argType := fnType.In(i)
			switch argType.Kind() {
			case reflect.Interface:
				comp, err := runner.getIntf(ctx, argType, "kodtest")
				if err != nil {
					return err
				}
				args[i] = comp
			case reflect.Pointer:
				comp, err := runner.getImpl(ctx, argType.Elem())
				if err != nil {
					return err
				}
				args[i] = comp
			default:
				return fmt.Errorf("argument %v has unexpected type %v", i, argType)
			}
		}

		reflect.ValueOf(fn).Call(lo.Map(args, func(item any, index int) reflect.Value { return reflect.ValueOf(item) }))
		return nil
	}, intfs, nil
}

func extractComponentInterfaceType(t reflect.Type) (reflect.Type, error) {
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("type %v is not a struct", t)
	}
	// See the definition of kod.Implements.
	f, ok := t.FieldByName("component_interface_type")
	if !ok {
		return nil, fmt.Errorf("type %v does not embed kod.Implements", t)
	}
	return f.Type, nil
}
