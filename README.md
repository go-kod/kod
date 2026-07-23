<div align="center">

[![Build and Test](https://github.com/go-kod/kod/actions/workflows/go.yml/badge.svg)](https://github.com/go-kod/kod/actions/workflows/go.yml)
[![GitHub Release](https://img.shields.io/github/v/release/go-kod/kod)](https://github.com/go-kod/kod/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-kod/kod)](https://goreportcard.com/report/github.com/go-kod/kod)
[![Code Cov](https://codecov.io/github/go-kod/kod/graph/badge.svg?token=FKCHAE6M2R)](https://codecov.io/github/go-kod/kod)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-kod/kod.svg)](https://pkg.go.dev/github.com/go-kod/kod)
[![GitHub License](https://img.shields.io/github/license/go-kod/kod)](https://github.com/go-kod/kod/blob/main/LICENSE)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

[**English**](./README.md) •
[**简体中文**](./README_CN.md) |
[**Examples**](./example_test.go) •
[**Template**](https://github.com/go-kod/kod-mono) •
[**Extention**](https://github.com/go-kod/kod-ext)

</div>

# Kod

Kod stands for **Killer Of Dependency**, a generics based dependency injection framework for Go.

> Although it seems that most Go enthusiasts dislike dependency injection framework, many companies that widely use Go for their development projects have open-sourced their own dependency injection frameworks. For example, Google has open-sourced Wire, Uber has open-sourced Fx, and Facebook has open-sourced Inject. This is truly a strange phenomenon.

![kod](/assets/kod.excalidraw.png)

## Feature

- **Component Based**: Kod is a component-based framework. Components are the building blocks of a Kod application.
- **Configurable**: Kod can use TOML/YAML/JSON files to configure how applications are run.
- **Testing**: Kod includes a Test function that you can use to test your Kod applications.
- **Logging**: Kod provides a logging API, `kod.L`. Kod also integrates the logs into the environment where your application is deployed.
- **OpenTelemetry**: Kod relies on OpenTelemetry to collect trace and metrics from your application.
- **Hooks**: Kod provides a way to run code when a component start or stop.
- **Interceptors**: Kod has built-in common interceptors, and components can implement the following methods to inject these interceptors into component methods.
- **Interface Generation**: Kod provides a way to generate interface from structure.
- **Code Generation**: Kod provides a way to generate kod related codes for your kod application.

## Installation

```bash
go install github.com/go-kod/kod/cmd/kod@latest
```

If the installation was successful, you should be able to run `kod -h`:

```bash
A powerful tool for writing kod applications.

Usage:
  kod [flags]
  kod [command]

Available Commands:
  callgraph        generate kod callgraph for your kod application.
  completion       Generate the autocompletion script for the specified shell
  generate         generate kod related codes for your kod application.
  help             Help about any command
  struct2interface generate interface from struct for your kod application.

Flags:
  -h, --help      help for kod
  -t, --toggle    Help message for toggle
  -v, --version   Help message for toggle

Use "kod [command] --help" for more information about a command.
```

## Step by Step Tutorial

In this section, we show you how to write Kod applications. To install Kod and follow along, refer to the Installation section. The full source code presented in this tutorial can be found here.

### Components

Kod's core abstraction is the component. A component is like an actor, and a Kod application is implemented as a set of components. Concretely, a component is represented with a regular Go interface, and components interact with each other by calling the methods defined by these interfaces.

In this section, we'll define a simple hello component that just prints a string and returns. First, run `go mod init hello` to create a go module.

```bash
mkdir hello/
cd hello/
go mod init hello
```

Then, create a file called `main.go` with the following contents:

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/go-kod/kod"
)

func main() {
    if err := kod.Run(context.Background(), serve); err != nil {
        log.Fatal(err)
    }
}

// app is the main component of the application. kod.Run creates
// it and passes it to serve.
type app struct{
    kod.Implements[kod.Main]
}

// serve is called by kod.Run and contains the body of the application.
func serve(context.Context, *app) error {
    fmt.Println("Hello")
    return nil
}
```

`kod.Run(...)` initializes and runs the Kod application. In particular, `kod.Run` finds the main component, creates it, and passes it to a supplied function. In this example, `app` is the main component since it contains a `kod.Implements[kod.Main]` field.

```bash
go mod tidy
kod generate .
go run .
Hello
```

## FUNDAMENTALS

### Components

Components are Kod's core abstraction. Concretely, a component is represented as a Go interface and corresponding implementation of that interface. Consider the following `Adder` component for example:

```go
type Adder interface {
    Add(context.Context, int, int) (int, error)
}

type adder struct {
    kod.Implements[Adder]
}

func (*adder) Add(_ context.Context, x, y int) (int, error) {
    return x + y, nil
}
```

Adder defines the component's interface, and adder defines the component's implementation. The two are linked with the embedded `kod.Implements[Adder]` field. You can call `kod.Ref[Adder].Get()` to get a caller to the Adder component.

#### Implementation

A component implementation must be a struct that looks like:

```go
type foo struct{
    kod.Implements[Foo]
    // ...
}
```

It must be a struct.
It must embed a `kod.Implements[T]` field where T is the component interface it implements.
If a component implementation implements an `Init(context.Context) error` method, it will be called when an instance of the component is created.

```go
func (f *foo) Init(context.Context) error {
    // ...
}

func (f *foo) Shutdown(context.Context) error {
    // ...
}
```

#### Lazy Initialization

Components can be lazily initialized by embedding a `kod.LazyInit` field in the component implementation, 
which will be initialized when the component is first used, instead of when the application starts.

Simple demo below:

```go
type foo struct {
    kod.Implements[Foo]
    kod.LazyInit
}
```

#### Interceptors

Kod has built-in common interceptors, and components can implement the following methods to inject these interceptors into component methods:

```go
func (f *foo) Interceptors() []interceptor.Interceptor {
    return []interceptor.Interceptor{
        kmetric.New(),
        ktrace.New(),
    }
}
```

#### Interfaces

Interface can be generated automatically by kod tool.

```go
//go:generate kod struct2interface .
```

### Configuration

Kod does not load application configuration. Treat configuration as an ordinary business component and choose your own source: flags, environment variables, files, secrets managers, or remote config.

```go
type Config interface {
    Greeting() string
}

type config struct {
    kod.Implements[Config]
    greeting string
}

func (c *config) Init(context.Context) error {
    c.greeting = os.Getenv("GREETING")
    if c.greeting == "" {
        c.greeting = "Hello"
    }
    return nil
}

func (c *config) Greeting() string { return c.greeting }

type greeter struct {
    kod.Implements[Greeter]
    config kod.Ref[Config]
}

func (g *greeter) Greet(_ context.Context, name string) (string, error) {
    return fmt.Sprintf("%s, %s!", g.config.Get().Greeting(), name), nil
}
```

### Testing

#### Unit Test

Kod includes a `Test` function that you can use to test your Kod applications. For example, create an `adder_test.go` file with the following contents.

```go
package main

import (
    "context"
    "testing"

    "github.com/go-kod/kod"
)

func TestAdd(t *testing.T) {
     kod.RunTest(t, func(ctx context.Context, adder Adder) {
         got, err := adder.Add(ctx, 1, 2)
         if err != nil {
             t.Fatal(err)
         }
         if want := 3; got != want {
             t.Fatalf("got %q, want %q", got, want)
         }
     })
}
```

Run go test to run the test. `kod.RunTest` will create a sub-test and within it will create an Adder component and pass it to the supplied function. If you want to test the implementation of a component, rather than its interface, specify a pointer to the implementing struct as an argument. For example, if the `adderImpl` struct implemented the Adder interface, we could write the following:

```go
kod.RunTest(t, func(ctx context.Context, adder *adderImpl) {
    // Test adder...
})
```

#### Benchmark

You can also use `kod.RunTest` to benchmark your application. For example, create an adder_benchmark.go file with the following contents.

```go
package main

import (
    "context"
    "testing"

    "github.com/go-kod/kod"
)

func BenchmarkAdd(b *testing.B) {
    kod.RunTest(b, func(ctx context.Context, adder Adder) {
        b.ResetTimer()
        for i := 0; i < b.N; i++ {
            _, err := adder.Add(ctx, 1, 2)
            if err != nil {
                b.Fatal(err)
            }
        }
    })
}
```

#### Fake

You can replace a component implementation with a fake implementation in a test using `kod.Fake`. Here's an example where we replace the real implementation of a Clock component with a fake implementation that always returns a fixed time.

```go
// fakeClock is a fake implementation of the Clock component.
type fakeClock struct {
    now int64
}

// Now implements the Clock component interface. It returns the current time, in
// microseconds, since the unix epoch.
func (f *fakeClock) Now(context.Context) (int64, error) {
    return f.now, nil
}

func TestClock(t *testing.T) {
    t.Run("fake", func(t *testing.T) {
        // Register a fake Clock implementation with the runner.
        fake := kod.Fake[Clock](&fakeClock{100})

        // When a fake is registered for a component, all instances of that
        // component dispatch to the fake.
        kod.RunTest(t, func(ctx context.Context, clock Clock) {
            now, err := clock.UnixMicro(ctx)
            if err != nil {
                t.Fatal(err)
            }
            if now != 100 {
                t.Fatalf("bad time: got %d, want %d", now, 100)
            }

            fake.now = 200
            now, err = clock.UnixMicro(ctx)
            if err != nil {
                t.Fatal(err)
            }
            if now != 200 {
                t.Fatalf("bad time: got %d, want %d", now, 200)
            }
        }, kod.WithFakes(fake))
    })
}
```

### Logging

Kod provides a logging API, `kod.L`. Kod also integrates the logs into the environment where your application is deployed.

Use the Logger method of a component implementation to get a logger scoped to the component. For example:

```go
type Adder interface {
    Add(context.Context, int, int) (int, error)
}

type adder struct {
    kod.Implements[Adder]
}

func (a *adder) Add(ctx context.Context, x, y int) (int, error) {
    // adder embeds kod.Implements[Adder] which provides the L method.
    logger := a.L(ctx)
    logger.DebugContext(ctx, "A debug log.")
    logger.InfoContext(ctx, "An info log.")
    logger.ErrorContext(ctx, "An error log.", fmt.Errorf("an error"))
    return x + y, nil
}
```

### OpenTelemetry

Kod relies on OpenTelemetry to collect trace and metrics from your application.

Supported Environment Variables:

- `OTEL_SDK_DISABLED`: If set to true, disables the OpenTelemetry SDK. Default is false.
- `OTEL_LOGS_EXPORTER`: The logs exporter to use. Supported values are "console" and "otlp", Default is "otlp".
- `OTEL_EXPORTER_OTLP_PROTOCOL`: The protocol to use for the OTLP exporter. Supported values are "grpc" and "http/protobuf", Default is "http/protobuf".
- `OTEL_EXPORTER_OTLP_INSECURE`: If set to true, disables the security features of the OTLP exporter. Default is false.

More information can be found at [OpenTelemetry Website](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/).

### Acknowledge

This project was heavily inspired by [ServiceWeaver](https://github.com/ServiceWeaver/weaver).

## Star History

<a>
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=go-kod/kod&type=Timeline&theme=dark" />
    <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=go-kod/kod&type=Timeline" />
    <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=go-kod/kod&type=Timeline" />
  </picture>
</a>
