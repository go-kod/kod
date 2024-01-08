<div align="center">

[![Build and Test](https://github.com/go-kod/kod/actions/workflows/go.yml/badge.svg)](https://github.com/go-kod/kod/actions/workflows/go.yml)
[![codecov](https://codecov.io/github/go-kod/kod/graph/badge.svg?token=OOZNKQJL38)](https://codecov.io/github/go-kod/kod)

[**English**](./README.md) •
[**简体中文**](./README_CN.md)

</div>

# Kod


Kod  stands for **Killer Of Dependency**, a generics based dependency injection toolkit for Go.

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

kod.Run(...) initializes and runs the Kod application. In particular, kod.Run finds the main component, creates it, and passes it to a supplied function. In this example,app is the main component since it contains a kod.Implements[kod.Main] field.

```bash
go mod tidy
kod generate .
go run .
Hello
```

## FUNDAMENTALS

### Components

Components are Kod's core abstraction. Concretely, a component is represented as a Go interface and corresponding implementation of that interface. Consider the following Adder component for example:

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

Adder defines the component's interface, and adder defines the component's implementation. The two are linked with the embedded kod.Implements[Adder] field. You can call kod.Ref[Adder].Get() to get a caller to the Adder component.

#### Implementation

A component implementation must be a struct that looks like:

```go
type foo struct{
    kod.Implements[Foo]
    // ...
}
```

It must be a struct.
It must embed a kod.Implements[T] field where T is the component interface it implements.
If a component implementation implements an Init(context.Context) error method, it will be called when an instance of the component is created.

```go
func (f *foo) Init(context.Context) error {
    // ...
}

func (f *foo) Stop(context.Context) error {
    // ...
}
```

#### Interceptors
Kod has built-in common interceptors, and components can implement the following methods to inject these interceptors into component methods:

```go
func (f *foo) Interceptors() []kod.Interceptor {
    return []kod.Interceptor{
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

### Config

Kod uses config files, written in TOML, to configure how applications are run. A minimal config file, for example, simply lists the application name:

```toml
[kod]
name = "hello"
```

A config file may also contain component-specific configuration sections, which allow you to configure the components in your application. For example, consider the following Greeter component.

```go
type Greeter interface {
    Greet(context.Context, string) (string, error)
}

type greeter struct {
    kod.Implements[Greeter]
}

func (g *greeter) Greet(_ context.Context, name string) (string, error) {
    return fmt.Sprintf("Hello, %s!", name), nil
}
```

Rather than hard-coding the greeting "Hello", we can provide a greeting in a config file. First, we define a options struct.

```go
type greeterOptions struct {
    Greeting string
}
```

Next, we associate the options struct with the greeter implementation by embedding the kod.WithConfig[T] struct.

```go
type greeter struct {
    kod.Implements[Greeter]
    kod.WithConfig[greeterOptions]
}
```

Now, we can add a Greeter section to the config file. The section is keyed by the full path-prefixed name of the component.

```toml
["example.com/mypkg/Greeter"]
Greeting = "Bonjour"
```

When the Greeter component is created, Kod will automatically parse the Greeter section of the config file into a greeterOptions struct. You can access the populated struct via the Config method of the embedded WithConfig struct. For example:

```go
func (g *greeter) Greet(_ context.Context, name string) (string, error) {
    greeting := g.Config().Greeting
    if greeting == "" {
        greeting = "Hello"
    }
    return fmt.Sprintf("%s, %s!", greeting, name), nil
}
```

You can use toml struct tags to specify the name that should be used for a field in a config file. For example, we can change the greeterOptions struct to the following.

```go
type greeterOptions struct {
    Greeting string `toml:"my_custom_name"`
}
```

### Testing

#### Unit Test

Kod includes a Test function that you can use to test your Kod applications. For example, create an adder_test.go file with the following contents.

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

Run go test to run the test. kod.Test will create a sub-test and within it will create an Adder component and pass it to the supplied function. If you want to test the implementation of a component, rather than its interface, specify a pointer to the implementing struct as an argument. For example, if the adderImpl struct implemented the Adder interface, we could write the following:

```go
kod.RunTest(t, func(ctx context.Context, adder *adderImpl) {
    // Test adder...
})
```

#### Benchmark

You can also use kod.Test to benchmark your application. For example, create an adder_benchmark.go file with the following contents.

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

You can replace a component implementation with a fake implementation in a test using kod.Fake. Here's an example where we replace the real implementation of a Clock component with a fake implementation that always returns a fixed time.

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

#### Config


You can also provide the contents of a config file to a runner by setting the Runner.Config field:

```go
func TestArithmetic(t *testing.T) {
    kod.RunTest(t, func(ctx context.Context, adder Adder) {
        // ...
    }, kod.WithConfigFile("testdata/config.toml"))
}
```

### Logging

Kod provides a logging API, kod.L. Kod also integrates the logs into the environment where your application is deployed.

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
    logger := a.L()
    logger.DebugContext(ctx, "A debug log.")
    logger.InfoContext(ctx, "An info log.")
    logger.ErrorContext(ctx, "An error log.", fmt.Errorf("an error"))
    return x + y, nil
}
```

### Opentelemetry

Kod relies on OpenTelemetry to collect trace and metrics from your application.

### Acknowledge

This project was heavily inspired by [ServiceWeaver](https://github.com/ServiceWeaver/kod).