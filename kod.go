package kod

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/internal/hooks"
	"github.com/go-kod/kod/internal/kslog"
	"github.com/go-kod/kod/internal/reflects"
	"github.com/go-kod/kod/internal/registry"
	"github.com/go-kod/kod/internal/signals"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/exporters/autoexport"
	"go.opentelemetry.io/contrib/instrumentation/host"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

const (
	PkgPath = "github.com/go-kod/kod"
)

// Implements[T any] provides a common structure for components,
// with logging capabilities and a reference to the component's interface.
type Implements[T any] struct {
	log *slog.Logger
	//nolint
	component_interface_type T
}

// L returns the associated logger.
func (i *Implements[T]) L(ctx context.Context) *slog.Logger {
	return kslog.LogWithContext(ctx, i.log)
}

// setLogger sets the logger for the component.
// nolint
func (i *Implements[T]) setLogger(log *slog.Logger) {
	i.log = log
}

// implements is a marker method to assert implementation of an interface.
// nolint
func (Implements[T]) implements(T) {}

// Ref[T any] is a reference holder to a value of type T.
// The reference is expected to be a field of a component struct.
// The value is set by the framework, and is accessible via the Get() method.
//
// Example:
//
//	type app struct {
//		kod.Implements[kod.Main]
//		component kod.Ref[example.Component]
//	}
//
//	func main() {
//		kod.Run(context.Background(), func(ctx context.Context, main *app) error {
//			component := main.component.Get()
//			// ...
//		})
//	}
type Ref[T any] struct {
	value T
}

// Get returns the held reference value.
func (r Ref[T]) Get() T { return r.value }

// isRef is a marker method to identify a Ref type.
// nolint
func (r Ref[T]) isRef() {}

// setRef sets the reference value.
// nolint
func (r *Ref[T]) setRef(val any) { r.value = val.(T) }

// Main is the interface that should be implemented by an application's main component.
// The main component is the entry point of the application,
// and is expected to be a struct that embeds Implements[Main].
//
// Example:
//
//	type app struct {
//		kod.Implements[kod.Main]
//	}
//
//	func main() {
//		kod.Run(context.Background(), func(ctx context.Context, main *app) error {
//			fmt.Println("Hello, World!")
//			return nil
//		})
//	}
type Main interface{}

// PointerToMain is a type constraint that asserts *T is an instance of Main
// (i.e. T is a struct that embeds weaver.Implements[weaver.Main]).
type PointerToMain[T any] interface {
	*T
	InstanceOf[Main]
}

// InstanceOf[T any] is an interface for asserting implementation of an interface T.
type InstanceOf[T any] interface {
	implements(T)
}

// WithConfig[T any] is a struct to hold configuration of type T.
// The struct is expected to be a field of a component struct.
// The configuration is loaded from a file, and is accessible via the Config() method.
//
// Example:
//
//	type app struct {
//		kod.Implements[kod.Main]
//		kod.WithConfig[appConfig]
//	}
//
//	type appConfig struct {
//		Host string
//		Port int
//	}
//
//	func main() {
//		kod.Run(context.Background(), func(ctx context.Context, main *app) error {
//			fmt.Println("config:", main.Config())
//		})
//	}
type WithConfig[T any] struct {
	config T
}

// Config returns a pointer to the config.
func (wc *WithConfig[T]) Config() *T {
	return &wc.config
}

// getConfig returns the config.
func (wc *WithConfig[T]) getConfig() any {
	return &wc.config
}

// WithConfigFile is an option setter for specifying a configuration file.
func WithConfigFile(filename string) func(*options) {
	return func(opts *options) {
		opts.configFilename = filename
	}
}

// WithFakes is an option setter for specifying fake components for testing.
func WithFakes(fakes ...fakeComponent) func(*options) {
	return func(opts *options) {
		opts.fakes = lo.SliceToMap(fakes, func(f fakeComponent) (reflect.Type, any) { return f.intf, f.impl })
	}
}

// WithLogWrapper is an option setter for specifying a logger.
func WithLogWrapper(h func(slog.Handler) slog.Handler) func(*options) {
	return func(opts *options) {
		opts.logWrapper = h
	}
}

// WithRegistrations is an option setter for specifying component registrations.
func WithRegistrations(regs ...*Registration) func(*options) {
	return func(opts *options) {
		opts.registrations = regs
	}
}

// WithInterceptors is an option setter for specifying interceptors.
func WithInterceptors(interceptors ...interceptor.Interceptor) func(*options) {
	return func(opts *options) {
		opts.interceptors = interceptors
	}
}

// Run initializes and runs the application with the provided main component and options.
func Run[T any, _ PointerToMain[T]](ctx context.Context, run func(context.Context, *T) error, opts ...func(*options)) error {
	opt := &options{}
	for _, o := range opts {
		o(opt)
	}

	// Create a new Kod instance.
	kod, err := newKod(ctx, *opt)
	if err != nil {
		return err
	}

	ctx, span := otel.Tracer(PkgPath).Start(ctx, "kod.Run")

	// create a new context with kod
	ctx = newContext(ctx, kod)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// get the main component implementation
	main, err := kod.getImpl(ctx, reflects.TypeFor[T]())
	if err != nil {
		return err
	}

	// wait for shutdown signal
	stop := make(chan struct{}, 2)
	signals.Shutdown(ctx, func(grace bool) {
		kod.log.InfoContext(ctx, "kod.Shutdown ...")
		cancel()
		stop <- struct{}{}
	})

	// run the main component
	go func() {
		err = run(ctx, main.(*T))
		stop <- struct{}{}
	}()

	span.End()
	// wait for stop signal
	<-stop

	ctx, timeoutCancel := context.WithTimeout(
		context.WithoutCancel(ctx), kod.config.ShutdownTimeout)
	defer timeoutCancel()

	// run hook functions
	kod.hooker.Do(ctx)

	return err
}

// kodConfig defines the overall configuration for the Kod application.
type kodConfig struct {
	Name    string
	Env     string
	Version string

	ShutdownTimeout time.Duration
}

// Kod represents the core structure of the application, holding configuration and component registrations.
type Kod struct {
	mu *sync.Mutex

	config kodConfig

	viper *viper.Viper
	log   *slog.Logger

	hooker *hooks.Hooker

	regs                []*Registration
	registryByName      map[string]*Registration
	registryByInterface map[reflect.Type]*Registration
	registryByImpl      map[reflect.Type]*Registration

	components map[string]any
	opts       options
}

// options defines the configuration options for Kod.
type options struct {
	configFilename string
	fakes          map[reflect.Type]any
	logWrapper     func(slog.Handler) slog.Handler
	registrations  []*Registration
	interceptors   []interceptor.Interceptor
}

// newKod creates a new instance of Kod with the provided registrations and options.
func newKod(ctx context.Context, opts options) (*Kod, error) {
	kod := &Kod{
		mu: &sync.Mutex{},
		config: kodConfig{
			Name:            filepath.Base(lo.Must(os.Executable())),
			Env:             "local",
			ShutdownTimeout: 5 * time.Second,
		},
		hooker:              hooks.New(),
		regs:                registry.All(),
		registryByName:      make(map[string]*Registration),
		registryByInterface: make(map[reflect.Type]*Registration),
		registryByImpl:      make(map[reflect.Type]*Registration),
		components:          make(map[string]any),
		opts:                opts,
	}

	kod.initOpenTelemetry(ctx)

	kod.register(opts.registrations)

	if err := kod.parseConfig(opts.configFilename); err != nil {
		return nil, err
	}

	if err := validateRegistrations(kod.regs); err != nil {
		return nil, err
	}

	if err := checkCircularDependency(kod.regs); err != nil {
		return nil, err
	}

	return kod, nil
}

// Config returns the current configuration of the Kod instance.
func (k *Kod) Config() kodConfig {
	return k.config
}

// L() returns the logger of the Kod instance.
func (k *Kod) L(ctx context.Context) *slog.Logger {
	return kslog.LogWithContext(ctx, k.log)
}

// register adds the given implementations to the Kod instance.
func (k *Kod) register(regs []*Registration) {
	if len(regs) > 0 {
		k.regs = regs
	}

	for _, v := range k.regs {
		k.registryByName[v.Name] = v
		k.registryByInterface[v.Interface] = v
		k.registryByImpl[v.Impl] = v
	}
}

// parseConfig parses the config file.
func (k *Kod) parseConfig(filename string) error {
	noConfigProvided := false
	if filename == "" {
		filename = os.Getenv("KOD_CONFIG")
		if filename == "" {
			noConfigProvided = true
			filename = "kod.toml"
		}
	}

	vip := viper.New()

	vip.SetConfigFile(filename)
	vip.AddConfigPath(".")
	err := vip.ReadInConfig()
	if err != nil {
		switch err.(type) {
		case viper.ConfigFileNotFoundError, *fs.PathError:
			if !noConfigProvided {
				fmt.Fprintln(os.Stderr, "failed to load config file, use default config")
			}
		default:
			return fmt.Errorf("read config file: %w", err)
		}
	}

	k.viper = vip

	return vip.UnmarshalKey("kod", &k.config)
}

func (k *Kod) initOpenTelemetry(ctx context.Context) {
	if os.Getenv("OTEL_SDK_DISABLED") == "true" {
		k.log = slog.Default()
		return
	}

	lo.Must0(host.Start())
	lo.Must0(runtime.Start())

	res := lo.Must(resource.New(ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(k.config.Name),
			semconv.ServiceVersionKey.String(k.config.Version),
			semconv.DeploymentEnvironmentKey.String(k.config.Env),
		)),
	)

	metricReader := lo.Must(autoexport.NewMetricReader(ctx))
	metricProvider := metric.NewMeterProvider(
		metric.WithReader(metricReader),
		metric.WithResource(res),
	)

	otel.SetMeterProvider(metricProvider)

	spanExporter := lo.Must(autoexport.NewSpanExporter(ctx))
	spanProvider := trace.NewTracerProvider(
		trace.WithBatcher(spanExporter),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(spanProvider)

	logExporter := lo.Must(getLogAutoExporter(ctx))
	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(
			log.NewBatchProcessor(logExporter),
		),
		log.WithResource(res),
	)

	global.SetLoggerProvider(loggerProvider)

	k.hooker.Add(hooks.HookFunc{
		Name: "OpenTelemetry",
		Fn: func(ctx context.Context) error {
			_ = metricProvider.Shutdown(ctx)
			_ = spanProvider.Shutdown(ctx)
			_ = loggerProvider.Shutdown(ctx)

			return nil
		},
	})

	var handler slog.Handler
	if k.opts.logWrapper != nil {
		handler = k.opts.logWrapper(otelslog.NewHandler(k.config.Name,
			otelslog.WithSchemaURL(PkgPath),
			otelslog.WithVersion(k.config.Version),
		))
	} else {
		handler = otelslog.NewHandler(k.config.Name,
			otelslog.WithSchemaURL(PkgPath),
			otelslog.WithVersion(k.config.Version),
		)
	}

	k.log = slog.New(handler)
}

func getLogAutoExporter(ctx context.Context) (log.Exporter, error) {
	var (
		err      error
		exporter log.Exporter
	)

	logsExporter := os.Getenv("OTEL_LOGS_EXPORTER")
	if logsExporter == "" {
		logsExporter = "otlp"
	}

	switch logsExporter {
	case "otlp":

		proto := os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL")
		if proto == "" {
			proto = "http/protobuf"
		}

		switch proto {
		case "grpc":
			exporter, err = otlploggrpc.New(ctx)
		case "http/protobuf":
			opts := []otlploghttp.Option{}
			if os.Getenv("OTEL_EXPORTER_OTLP_INSECURE") == "true" {
				opts = append(opts, otlploghttp.WithInsecure())
			}

			exporter, err = otlploghttp.New(ctx, opts...)
		default:
			return nil, fmt.Errorf("unsupported OTLP exporter protocol: %s", proto)
		}
	case "console":
		exporter, err = stdoutlog.New()
	}

	return exporter, err
}
