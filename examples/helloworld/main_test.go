package main

import (
	"context"
	"testing"

	"github.com/go-kod/kod"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

// TestApp tests the app component.
func TestApp(t *testing.T) {
	t.Parallel()

	kod.RunTest(t, func(ctx context.Context, comp *app) {
		assert.Equal(t, "Hello, World!", comp.helloworld.Get().SayHello())
	})
}

// TestAppMockHelloworld tests the app component with a mock helloworld.
func TestAppMockHelloworld(t *testing.T) {
	t.Parallel()

	mockHelloworld := NewMockHelloworld(gomock.NewController(t))
	mockHelloworld.EXPECT().SayHello().Return("Hello, Mocker!").Times(1)

	kod.RunTest(t, func(ctx context.Context, comp *app) {
		assert.Equal(t, "Hello, Mocker!", comp.helloworld.Get().SayHello())
	}, kod.WithFakes(kod.Fake[Helloworld](mockHelloworld)))
}
