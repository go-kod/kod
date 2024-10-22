// Code generated by "kod generate". DO NOT EDIT.
//go:build !ignoreKodGen

package case5

import (
	"context"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"reflect"
)

// Full method names for components.
const ()

func init() {
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/Main",
		Interface: reflect.TypeOf((*kod.Main)(nil)).Elem(),
		Impl:      reflect.TypeOf(refStructImpl{}),
		Refs:      `⟦b915993d:KoDeDgE:github.com/go-kod/kod/Main→github.com/go-kod/kod/tests/case5/testRefStruct1⟧`,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			return main_local_stub{
				impl:        info.Impl.(kod.Main),
				interceptor: info.Interceptor,
				name:        info.Name,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/case5/TestRefStruct1",
		Interface: reflect.TypeOf((*TestRefStruct1)(nil)).Elem(),
		Impl:      reflect.TypeOf(testRefStruct1{}),
		Refs:      ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			return testRefStruct1_local_stub{
				impl:        info.Impl.(TestRefStruct1),
				interceptor: info.Interceptor,
				name:        info.Name,
			}
		},
	})
}

// kod.InstanceOf checks.
var _ kod.InstanceOf[kod.Main] = (*refStructImpl)(nil)
var _ kod.InstanceOf[TestRefStruct1] = (*testRefStruct1)(nil)

// Local stub implementations.

// main_local_stub is a local stub implementation of [kod.Main].
type main_local_stub struct {
	impl        kod.Main
	name        string
	interceptor interceptor.Interceptor
}

// Check that [main_local_stub] implements the [kod.Main] interface.
var _ kod.Main = (*main_local_stub)(nil)

// testRefStruct1_local_stub is a local stub implementation of [TestRefStruct1].
type testRefStruct1_local_stub struct {
	impl        TestRefStruct1
	name        string
	interceptor interceptor.Interceptor
}

// Check that [testRefStruct1_local_stub] implements the [TestRefStruct1] interface.
var _ TestRefStruct1 = (*testRefStruct1_local_stub)(nil)

