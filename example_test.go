package kod_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/examples/helloworld"
	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/interceptor/kmetric"
	"github.com/go-kod/kod/interceptor/krecovery"
	"github.com/go-kod/kod/interceptor/ktrace"
	"github.com/go-kod/kod/internal/kslog"
	"go.uber.org/mock/gomock"
)

func Example_mainComponent() {
	kod.Run(context.Background(), func(ctx context.Context, app *helloworld.App) error {
		fmt.Println("Hello, World!")
		return nil
	})
	// Output:
	// helloWorld init
	// Hello, World!
	// helloWorld shutdown
}

func Example_componentRefAndCall() {
	kod.Run(context.Background(), func(ctx context.Context, app *helloworld.App) error {
		app.HelloWorld.Get().SayHello(ctx)
		return nil
	})
	// Output:
	// helloWorld init
	// Hello, World!
	// helloWorld shutdown
}

func Example_componentLazyInit() {
	kod.Run(context.Background(), func(ctx context.Context, app *helloworld.App) error {
		app.HelloLazy.Get().SayHello(ctx)
		app.HelloWorld.Get().SayHello(ctx)
		return nil
	})
	// Output:
	// helloWorld init
	// lazyHelloBob init
	// Hello, Bob!
	// Hello, World!
	// lazyHelloBob shutdown
	// helloWorld shutdown
}

func Example_componentMock() {
	mock := helloworld.NewMockHelloWorld(gomock.NewController(nil))
	mock.EXPECT().SayHello(gomock.Any()).Return()

	kod.Run(context.Background(), func(ctx context.Context, app *helloworld.App) error {
		app.HelloWorld.Get().SayHello(ctx)
		fmt.Println("Nothing printed from mock")
		return nil
	}, kod.WithFakes(kod.Fake[helloworld.HelloWorld](mock)))
	// Output:
	// Nothing printed from mock
}

func Example_config() {
	kod.Run(context.Background(), func(ctx context.Context, app *helloworld.App) error {
		fmt.Println(app.Config().Name)
		app.HelloWorld.Get().SayHello(ctx)
		return nil
	}, kod.WithConfigFile("./examples/helloworld/config.toml"))
	// Output:
	// helloWorld init
	// globalConfig
	// Hello, World!config
	// helloWorld shutdown
}

func Example_log() {
	logger, observer := kslog.NewTestLogger()

	kod.RunTest(&testing.T{}, func(ctx context.Context, app *helloworld.App) {
		app.L(ctx).Debug("Hello, World!")
		app.L(ctx).Info("Hello, World!")
		app.L(ctx).Warn("Hello, World!")
		app.L(ctx).Error("Hello, World!")
		app.HelloWorld.Get().SayHello(ctx)
	}, kod.WithLogger(logger))

	fmt.Println(observer)
	// Output:
	// helloWorld init
	// Hello, World!
	// helloWorld shutdown
	// {"level":"INFO","msg":"Hello, World!","component":"github.com/go-kod/kod/Main"}
	// {"level":"WARN","msg":"Hello, World!","component":"github.com/go-kod/kod/Main"}
	// {"level":"ERROR","msg":"Hello, World!","component":"github.com/go-kod/kod/Main"}
	// {"level":"INFO","msg":"Hello, World!","component":"github.com/go-kod/kod/examples/helloworld/HelloWorld"}
}

func Example_interceptor() {
	interceptor := interceptor.Interceptor(func(ctx context.Context, info interceptor.CallInfo, req, res []interface{}, next interceptor.HandleFunc) error {
		fmt.Println("Before call")
		err := next(ctx, info, req, res)
		fmt.Println("After call")
		return err
	})

	kod.Run(context.Background(), func(ctx context.Context, app *helloworld.App) error {
		app.HelloWorld.Get().SayHello(ctx)
		return nil
	}, kod.WithInterceptors(interceptor))
	// Output:
	// helloWorld init
	// Before call
	// Hello, World!
	// After call
	// helloWorld shutdown
}

func Example_interceptorBuiltin() {
	kod.Run(context.Background(), func(ctx context.Context, app *helloworld.App) error {
		app.HelloWorld.Get().SayHello(ctx)
		return nil
	}, kod.WithInterceptors(krecovery.Interceptor(), ktrace.Interceptor(), kmetric.Interceptor()))
	// Output:
	// helloWorld init
	// Hello, World!
	// helloWorld shutdown
}

func Example_test() {
	kod.RunTest(&testing.T{}, func(ctx context.Context, app *helloworld.App) {
		app.HelloWorld.Get().SayHello(ctx)
	})
	// Output:
	// helloWorld init
	// Hello, World!
	// helloWorld shutdown
}

func Example_testWithMockComponent() {
	mock := helloworld.NewMockHelloWorld(gomock.NewController(nil))
	mock.EXPECT().SayHello(gomock.Any()).Return()

	kod.RunTest(&testing.T{}, func(ctx context.Context, app *helloworld.App) {
		app.HelloWorld.Get().SayHello(ctx)
		fmt.Println("Nothing printed from mock")
	}, kod.WithFakes(kod.Fake[helloworld.HelloWorld](mock)))
	// Output:
	// Nothing printed from mock
}

func Example_testWithConfig() {
	kod.RunTest(&testing.T{}, func(ctx context.Context, app *helloworld.App) {
		fmt.Println(app.Config().Name)
		app.HelloWorld.Get().SayHello(ctx)
	}, kod.WithConfigFile("./examples/helloworld/config.toml"))
	// Output:
	// helloWorld init
	// globalConfig
	// Hello, World!config
	// helloWorld shutdown
}

// Example_testWithLogObserver demonstrates how to test log output.
func Example_testWithLogObserver() {
	kod.RunTest(&testing.T{}, func(ctx context.Context, app *helloworld.App) {
		app.L(ctx).Debug("Hello, World!")
		app.L(ctx).Info("Hello, World!")
		app.L(ctx).Warn("Hello, World!")
		app.L(ctx).Error("Hello, World!")
	})

	// Output:
	// helloWorld init
	// helloWorld shutdown
}
