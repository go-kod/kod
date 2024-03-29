// Code generated by kod struct2interface; DO NOT EDIT.

package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// test1ControllerImpl is a component that implements test1Controller.
type test1Controller interface {
	Foo(cccccc *gin.Context)
}

// httpControllerImpl is a component that implements HTTPController.
type HTTPController interface {
	Foo(w http.ResponseWriter, r http.Request)
}

// serviceImpl is a component that implements testService.
type testService interface {
	Foo(ctx context.Context) error
}

// modelImpl is a component that implements testModel.
type testModel interface {
	Foo(ctx context.Context) error
}

// test1Component is a component that implements Test1Component.
type Test1Component interface {
	Foo(ctx context.Context, req *FooReq) error
}
