package helloworld

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
		assert.Equal(t, "Hello, World!", comp.helloWorld.Get().SayHello())
	})
}

// TestAppMockHelloWorld tests the app component with a mock helloWorld.
func TestAppMockHelloWorld(t *testing.T) {
	t.Parallel()

	mockHelloWorld := NewMockHelloWorld(gomock.NewController(t))
	mockHelloWorld.EXPECT().SayHello().Return("Hello, Mocker!").Times(1)

	kod.RunTest(t, func(ctx context.Context, comp *app) {
		assert.Equal(t, "Hello, Mocker!", comp.helloWorld.Get().SayHello())
	}, kod.WithFakes(kod.Fake[HelloWorld](mockHelloWorld)))
}

// BenchmarkAppHelloWorld benchmarks the app component.
func BenchmarkAppHelloWorld(b *testing.B) {
	kod.RunTest(b, func(ctx context.Context, comp *app) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			comp.helloWorld.Get().SayHello()
		}
	})
}
