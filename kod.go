package kod

import (
	"context"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/creasty/defaults"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"

	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/internal/hooks"
	"github.com/go-kod/kod/internal/registry"
	"github.com/go-kod/kod/internal/signals"
)

const (
	PkgPath = "github.com/go-kod/kod"
)

// Implements[T any] provides a common structure for components,
// with logging/tracing/metrics capabilities and a reference to the component's interface.
type Implements[T any] struct {
	name string
	//nolint
	component_interface_type T
}

// L returns the associated logger.
func (i *Implements[T]) L(ctx context.Context) *slog.Logger {
	s := trace.SpanContextFromContext(ctx)
	if s.IsValid() {
		return slog.Default().With(
			slog.String("component", i.name),
			slog.String("span_id", s.SpanID().String()),
			slog.String("trace_id", s.TraceID().String()),
		)
	}

	return slog.Default().With(slog.String("component", i.name))
}

// Tracer return the associated tracer.
func (i *Implements[T]) Tracer(opts ...trace.TracerOption) trace.Tracer {
	return otel.Tracer(i.name, opts...)
}

// Meter return the associated meter.
func (i *Implements[T]) Meter(opts ...metric.MeterOption) metric.Meter {
	return otel.Meter(i.name, opts...)
}

// setName sets the name for the component.
// nolint
func (i *Implements[T]) setName(name string) {
	i.name = name
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
	value  T
	once   sync.Once
	getter componentGetter
}

// Get returns the held reference value.
func (r *Ref[T]) Get() T {
	r.init()
	return r.value
}

// isRef is a marker method to identify a Ref type.
// nolint
func (r Ref[T]) isRef() {}

// setRef sets the reference value.
// nolint
func (r *Ref[T]) setRef(lazyInit bool, getter componentGetter) {
	r.getter = getter
	if !lazyInit {
		r.init()
	}
}

// init initializes the reference value.
func (r *Ref[T]) init() {
	r.once.Do(func() { r.value = lo.Must(r.getter()).(T) })
}

// componentGetter is a function type for getting a reference value.
type componentGetter func() (any, error)

// LazyInit is a marker type for lazy initialization of components.
type LazyInit struct{}

// isLazyInit is a marker method to identify a LazyInit type.
// nolint
func (r LazyInit) isLazyInit() {}

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
// (i.e. T is a struct that embeds kod.Implements[kod.Main]).
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

// WithGlobalConfig[T any] is a struct to hold global configuration of type T.
// The struct is expected to be a field of a component struct.
// The configuration is loaded from a file, and is accessible via the Config() method.
type WithGlobalConfig[T any] struct {
	config T
}

// Config returns a pointer to the config.
func (wc *WithGlobalConfig[T]) Config() *T {
	return &wc.config
}

// getGlobalConfig returns the config.
func (wc *WithGlobalConfig[T]) getGlobalConfig() any {
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

// MustRun is a helper function to run the application with the provided main component and options.
// It panics if an error occurs during the execution.
func MustRun[T any, P PointerToMain[T]](ctx context.Context, run func(context.Context, *T) error, opts ...func(*options)) {
	lo.Must0(Run[T, P](ctx, run, opts...))
}

// Run initializes and runs the application with the provided main component and options.
func Run[T any, _ PointerToMain[T]](ctx context.Context, run func(context.Context, *T) error, opts ...func(*options)) error {
	// Create a new Kod instance.
	kod, err := newKod(ctx, opts...)
	if err != nil {
		return err
	}

	// create a new context with kod
	ctx = newContext(ctx, kod)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// get the main component implementation
	main, err := kod.getImpl(ctx, reflect.TypeFor[T]())
	if err != nil {
		return err
	}

	// wait for shutdown signal
	stop := make(chan struct{}, 2)
	sig := make(chan os.Signal, 2)
	signals.Shutdown(ctx, sig, func(_ bool) {
		cancel()
		stop <- struct{}{}
	})

	// run the main component
	go func() {
		err = run(ctx, main.(*T))
		stop <- struct{}{}
	}()

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

	cfg *koanf.Koanf

	hooker *hooks.Hooker

	regs                []*Registration
	registryByName      map[string]*Registration
	registryByInterface map[reflect.Type]*Registration
	registryByImpl      map[reflect.Type]*Registration

	components         map[string]any
	lazyInitComponents map[reflect.Type]bool
	opts               *options
}

// options defines the configuration options for Kod.
type options struct {
	configFilename string
	fakes          map[reflect.Type]any
	registrations  []*Registration
	interceptors   []interceptor.Interceptor
}

// newKod creates a new instance of Kod with the provided registrations and options.
func newKod(_ context.Context, opts ...func(*options)) (*Kod, error) {
	opt := &options{}
	for _, o := range opts {
		o(opt)
	}

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
		opts:                opt,
	}

	kod.register(opt.registrations)

	err := kod.parseConfig(opt.configFilename)
	if err != nil {
		return nil, err
	}

	kod.lazyInitComponents, err = processRegistrations(kod.regs)
	if err != nil {
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

// unmarshalConfig sets the configuration for the Kod instance.
func (k *Kod) unmarshalConfig(key string, out interface{}) error {
	err := defaults.Set(out)
	if err != nil {
		return fmt.Errorf("set defaults: %w", err)
	}

	return k.cfg.Unmarshal(key, out)
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

	c := koanf.New(".")
	err := c.Load(env.Provider("KOD_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(s), "_", ".", -1)
	}), nil)
	if err != nil {
		return fmt.Errorf("load env config: %w", err)
	}

	// get ext
	ext := filepath.Ext(filename)
	switch ext {
	case ".toml":
		err = c.Load(file.Provider(filename), toml.Parser())
	case ".yaml":
		err = c.Load(file.Provider(filename), yaml.Parser())
	case ".json":
		err = c.Load(file.Provider(filename), json.Parser())
	default:
		return fmt.Errorf("read config file: Unsupported Config Type %q", ext)
	}

	if err != nil {
		switch err.(type) {
		case *fs.PathError:
			if noConfigProvided {
				fmt.Fprintln(os.Stderr, "failed to load config file, use default config")
			} else {
				return fmt.Errorf("read config file: %w", err)
			}
		default:
			return fmt.Errorf("read config file: %w", err)
		}
	}

	k.cfg = c
	err = c.Unmarshal("kod", &k.config)
	if err != nil {
		return fmt.Errorf("unmarshal config: %w", err)
	}

	return nil
}
