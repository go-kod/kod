package case1

import (
	"context"
	"errors"
	"fmt"
	"syscall"
	"testing"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/internal/mock"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/goleak"
	"go.uber.org/mock/gomock"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m,
		goleak.IgnoreAnyFunction("github.com/go-kod/kod/interceptor/internal/ratelimit.cpuproc"),
		goleak.IgnoreAnyFunction("gopkg.in/natefinch/lumberjack%2ev2.(*Logger).millRun"),
	)
}

func TestRun(t *testing.T) {
	t.Parallel()

	t.Run("case1", func(t *testing.T) {
		err := kod.Run(context.Background(), Run)
		require.Equal(t, "test1:B", err.Error())
	})
}

func TestImpl(t *testing.T) {
	t.Parallel()
	kod.RunTest(t, func(ctx context.Context, k *test1Component) {
		_, err := k.Foo(ctx, &FooReq{})
		fmt.Println(err)
		require.Equal(t, "test1:B", err.Error())
	})
}

func TestInterface(t *testing.T) {
	t.Parallel()
	kod.RunTest(t, func(ctx context.Context, k Test1Component) {
		ctx = StartTrace(ctx)

		ctx, span := otel.Tracer("").Start(ctx, "Run", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			span.End()
			fmt.Println("!!!!!!")
		}()

		_, err := k.Foo(ctx, &FooReq{Id: 1})
		res, err := k.Foo(ctx, &FooReq{Id: 2})
		fmt.Println(err)
		require.Equal(t, "test1:B", err.Error())
		require.True(t, span.SpanContext().IsValid())
		require.Equal(t, 2, res.Id)
	})
}

func TestInterfacePanic(t *testing.T) {
	t.Parallel()
	kod.RunTest(t, func(ctx context.Context, k Test1Component) {
		_, err := k.Foo(ctx, &FooReq{
			Panic: true,
		})
		require.Contains(t, err.Error(), "panic caught: test panic")
	})
}

func TestInterfacValidate(t *testing.T) {
	t.Parallel()
	kod.RunTest(t, func(ctx context.Context, k Test1Component) {
		_, err := k.Foo(ctx, &FooReq{
			Id: 101,
		})
		require.Contains(t, err.Error(), "validate failed: Key: 'FooReq.Id' Error:Field validation for 'Id' failed on the 'lt' tag")
	})
}

func TestFake(t *testing.T) {
	t.Parallel()
	fakeTest1 := &fakeTest1Component{"B"}
	kod.RunTest(t, func(ctx context.Context, k Test1Component) {
		_, err := k.Foo(ctx, &FooReq{})
		fmt.Println(err)
		require.Equal(t, errors.New("A:B"), err)
	}, kod.WithFakes(kod.Fake[Test1Component](fakeTest1)))
}

func TestFakeWithMock(t *testing.T) {
	t.Parallel()
	fakeTest1 := NewMockTest1Component(gomock.NewController(t))
	fakeTest1.EXPECT().Foo(gomock.Any(), gomock.Any()).Return(&FooRes{}, errors.New("A:B"))
	kod.RunTest(t, func(ctx context.Context, k Test1Component) {
		_, err := k.Foo(ctx, &FooReq{})
		fmt.Println(err)
		require.Equal(t, errors.New("A:B"), err)
	}, kod.WithFakes(kod.Fake[Test1Component](fakeTest1)))
}

func TestConflictFake(t *testing.T) {
	t.Parallel()
	fakeTest1 := &fakeTest1Component{"B"}
	mock.ExpectFailure(t, func(tt testing.TB) {
		kod.RunTest(tt, func(ctx context.Context, k *test1Component) {
			_, err := k.Foo(ctx, &FooReq{})
			fmt.Println(err)
			require.Equal(t, errors.New("A:B"), err)
		}, kod.WithFakes(kod.Fake[Test1Component](fakeTest1)))
	})
}

func TestConfigFile1(t *testing.T) {
	t.Parallel()
	kod.RunTest(t, func(ctx context.Context, k *test1Component) {
		_, err := k.Foo(ctx, &FooReq{})
		fmt.Println(err)
		require.Equal(t, "B", k.Config().A)
		require.Equal(t, "test1:B", err.Error())
	}, kod.WithConfigFile("kod.toml"))
}

func TestConfigFile2(t *testing.T) {
	t.Parallel()
	kod.RunTest(t, func(ctx context.Context, k *test1Component) {
		_, err := k.Foo(ctx, &FooReq{})
		fmt.Println(err)
		require.Equal(t, "test1:B2", err.Error())
	}, kod.WithConfigFile("kod2.toml"))
}

func TestConfigNotFound(t *testing.T) {
	t.Parallel()
	kod.RunTest(t, func(ctx context.Context, k *test1Component) {
		_, err := k.Foo(ctx, &FooReq{})
		fmt.Println(err)
		require.Equal(t, "test1:", err.Error())
	}, kod.WithConfigFile("kod-notfound.toml"))
}

func TestRunKill(t *testing.T) {
	t.Run("case1", func(t *testing.T) {
		err := kod.Run(context.Background(), Run)

		require.Nil(t, syscall.Kill(syscall.Getpid(), syscall.SIGINT))

		require.Equal(t, "test1:B", err.Error())
	})
}

func TestPanicKod(t *testing.T) {
	kod.RunTest(t, func(ctx context.Context, k *test1Component) {
		require.Panics(t, func() {
			kod := kod.FromContext(context.Background())
			kod.Config()
		})
	})
}

func BenchmarkCase1(b *testing.B) {
	b.Run("case1", func(b *testing.B) {
		kod.RunTest(b, func(ctx context.Context, k *test1Component) {
			for i := 0; i < b.N; i++ {
				_, err := k.Foo(ctx, &FooReq{})
				require.Equal(b, "B", k.Config().A)
				require.Equal(b, "test1:B", err.Error())
			}
		})
	})
}
