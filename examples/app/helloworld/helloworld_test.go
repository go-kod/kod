package helloworld

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/go-kod/kod"
)

// TestAppMockHelloWorld tests the app component with a mock helloWorld.
func TestAppMockHelloWorld(t *testing.T) {
	t.Parallel()

	mockHelloWorld := NewMockHelloWorld(gomock.NewController(t))
	mockHelloWorld.EXPECT().SayHello(context.TODO()).Return("Hello, Mocker!").Times(1)

	kod.RunTest(t, func(ctx context.Context, comp HelloWorld) {
		assert.Equal(t, "Hello, Mocker!", comp.SayHello(context.TODO()))
	}, kod.WithFakes(kod.Fake[HelloWorld](mockHelloWorld)))
}

// BenchmarkAppHelloWorld benchmarks the app component.
func BenchmarkAppHelloWorld(b *testing.B) {
	kod.RunTest(b, func(ctx context.Context, comp HelloWorld) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			comp.SayHello(context.TODO())
		}
	})
}
