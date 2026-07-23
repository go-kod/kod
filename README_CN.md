<div align="center">

[![Build and Test](https://github.com/go-kod/kod/actions/workflows/go.yml/badge.svg)](https://github.com/go-kod/kod/actions/workflows/go.yml)
[![GitHub Release](https://img.shields.io/github/v/release/go-kod/kod)](https://github.com/go-kod/kod/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-kod/kod)](https://goreportcard.com/report/github.com/go-kod/kod)
[![Code Cov](https://codecov.io/github/go-kod/kod/graph/badge.svg?token=FKCHAE6M2R)](https://codecov.io/github/go-kod/kod)
[![Go Reference](https://pkg.go.dev/badge/github.com/go-kod/kod.svg)](https://pkg.go.dev/github.com/go-kod/kod)
[![GitHub License](https://img.shields.io/github/license/go-kod/kod)](https://github.com/go-kod/kod/blob/main/LICENSE)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

[**英文**](./README.md) •
[**简体中文**](./README_CN.md) |
[**实例**](./example_test.go) •
[**模版**](https://github.com/go-kod/kod-mono) •
[**扩展**](https://github.com/go-kod/kod-ext)

</div>

# Kod

Kod代表“依赖杀手”，是基于泛型的Go语言依赖注入工具。

![kod](/assets/kod.excalidraw.png)

## 功能

- **基于组件**: Kod 是一个基于组件的框架。组件是 Kod 应用程序的构建模块。
- **可配置的**: Kod 可以使用 TOML/YAML/JSON 文件来配置应用程序的运行方式。
- **测试**: Kod 包含一个 Test 函数，您可以使用它来测试您的 Kod 应用程序。
- **日志记录**: Kod 提供了一个日志记录 API，kod.L。Kod 还将日志集成到部署您的应用程序的环境- 中。
- **OpenTelemetry**: Kod 依赖于 OpenTelemetry 来收集应用程序的跟踪和指标。
- **钩子**: Kod 提供了一种在组件启动或停止时运行代码的方式。
- **拦截器**: Kod 内置了常见的拦截器，组件可以实现以下方法来将这些拦截器注入到组件方法中。
- **接口生成**: Kod 提供了一种从结构体生成接口的方法。
- **代码生成**: Kod 提供了一种为您的 Kod 应用程序生成与 kod 相关代码的方法。

## 安装

```bash
go install github.com/go-kod/kod/cmd/kod@latest
```

如果安装成功，你应该能够运行 `kod -h`：

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

## 逐步教程

在这一部分，我们将向你展示如何编写 Kod 应用程序。要安装 Kod 并进行学习，请参考安装部分。在本教程中呈现的完整源代码可以在此处找到。

### 组件

Kod 的核心抽象是组件。组件类似于演员，而 Kod 应用程序是作为一组组件实现的。具体而言，组件用常规的 Go 接口表示，组件通过调用这些接口定义的方法来相互交互。

在这一部分，我们将定义一个简单的 hello 组件，它只是打印一个字符串并返回。首先，运行 `go mod init hello` 创建一个 Go 模块。

```bash
mkdir hello/
cd hello/
go mod init hello
```

然后，创建一个名为 `main.go` 的文件，其内容如下：

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

// app 是应用程序的主要组件。kod.Run 创建
// 它并将其传递给 serve 函数。
type app struct{
    kod.Implements[kod.Main]
}

// serve 由 kod.Run 调用，包含应用程序的主体。
func serve(context.Context, *app) error {
    fmt.Println("Hello")
    return nil
}
```

kod.Run(...) 初始化并运行 Kod 应用程序。具体而言，kod.Run 找到主要组件，创建它，并将它传递给提供的函数。在这个例子中，app 是主要组件，因为它包含一个 kod.Implements[kod.Main] 字段。

```bash
go mod tidy
kod struct2interface .
kod generate .
go run .
Hello
```

## 基础知识

### 组件

组件是Kod的核心抽象。具体而言，组件表示为Go接口及其相应的接口实现。例如，考虑以下Adder组件：

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

Adder定义了组件的接口，而adder定义了组件的实现。这两者通过嵌入的kod.Implements[Adder]字段相互关联。您可以调用kod.Ref[Adder].Get()来获取对Adder组件的调用者。

#### 实现🔗

组件的实现必须是一个类似以下结构的结构体：

```go
type foo struct {
    kod.Implements[Foo]
    // ...
}
```

它必须是一个结构体。
它必须嵌入一个kod.Implements[T]字段，其中T是它实现的组件接口。
如果组件实现实现了一个Init(context.Context) error方法，它将在创建组件实例时被调用。

```go
func (f *foo) Init(context.Context) error {
    // ...
}

func (f *foo) Shutdown(context.Context) error {
    // ...
}
```

#### 懒初始化

组件可以通过实现LazyInit接口来实现懒初始化。懒初始化组件的实例将在第一次调用组件方法时创建，而不是在应用程序启动时创建。

简单示例如下：

```go
type foo struct {
    kod.Implements[Foo]
    kod.LazyInit
}
```

#### 拦截器

Kod 内置了常用的拦截器，组件可以实现以下拦截器方法，将这些拦截器注入到组件方法中：

```go
func (f *foo) Interceptors() []interceptor.Interceptor {
    return []interceptor.Interceptor{
        kmetric.New(),
        ktrace.New(),
    }
}
```

#### 接口

接口可以通过kod工具自动生成。

```go
//go:generate kod struct2interface .
```

### 配置

Kod 不负责加载应用配置。配置应作为普通业务组件，由业务自行选择来源：命令行参数、环境变量、文件、密钥系统或远程配置。

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

### 测试

#### 单元测试

Kod包含一个Test函数，您可以使用它来测试Kod应用程序。例如，创建一个包含以下内容的adder_test.go文件。

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

运行go test来运行测试。kod.RunTest将创建一个子测试，并在其中创建一个Adder组件并将其传递给提供的函数。如果要测试组件的实现而不是其接口，请将指向实现结构的指针指定为参数。例如，如果adderImpl结构实现了Adder接口，我们可以编写如下内容：

```go
kod.RunTest(t, func(ctx context.Context, adder *adderImpl) {
    // Test adder...
})
```

#### 基准测试

您还可以使用kod.RunTest来对应用程序进行基准测试。例如，创建一个包含以下内容的adder_benchmark.go文件。

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

#### 伪造

您可以在测试中使用kod.Fake将组件实现替换为伪实现。以下是一个示例，其中我们将Clock组件的真实实现替换为总是返回固定时间的伪实现。

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

### 日志

Kod提供了一个日志API，kod.L。Kod还将日志集成到部署应用程序的环境中。

使用组件实现的Logger方法获取与该组件关联的记录器。例如：

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

### OpenTelemetry 可观测性

Kod 依赖于 OpenTelemetry 来从您的应用程序收集跟踪和指标信息。

支持的环境变量：

- `OTEL_SDK_DISABLED`: 如果设置为 true，则禁用 OpenTelemetry SDK。默认为 false。
- `OTEL_LOGS_EXPORTER`: 要使用的日志导出器。支持的值为 "console" 和 "otlp"，默认为 "otlp"。
- `OTEL_EXPORTER_OTLP_PROTOCOL`: 用于 OTLP 导出器的协议。支持的值为 "grpc" 和 "http/protobuf"，默认为 "http/protobuf"。
- `OTEL_EXPORTER_OTLP_INSECURE`: 如果设置为 true，则禁用 OTLP 导出器的安全功能。默认为 false。

更多信息可以在 [OpenTelemetry官网](https://opentelemetry.io/docs/specs/otel/configuration/sdk-environment-variables/) 找到。

### Acknowledge

这个项目在很大程度上受到了[ServiceWeaver](https://github.com/ServiceWeaver/weaver)的启发。

## Star History

<a>
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/svg?repos=go-kod/kod&type=Timeline&theme=dark" />
    <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/svg?repos=go-kod/kod&type=Timeline" />
    <img alt="Star History Chart" src="https://api.star-history.com/svg?repos=go-kod/kod&type=Timeline" />
  </picture>
</a>
