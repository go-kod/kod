// Code generated by "kod generate"; DO NOT EDIT.

package case3

import (
	"context"
)

// Test1Component is implemented by [test1Component],
// which can be mocked with [NewMockTest1Component].
type Test1Component interface {
	// Foo is implemented by [test1Component.Foo]
	Foo(ctx context.Context, req *FooReq) error
}

// Test2Component is implemented by [test2Component],
// which can be mocked with [NewMockTest2Component].
type Test2Component interface {
	// Foo is implemented by [test2Component.Foo]
	Foo(ctx context.Context, req *FooReq) error
}

// Test3Component is implemented by [test3Component],
// which can be mocked with [NewMockTest3Component].
type Test3Component interface {
	// Foo is implemented by [test3Component.Foo]
	Foo(ctx context.Context, req *FooReq) error
}
