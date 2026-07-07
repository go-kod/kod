package kod

import (
	"context"
	"log/slog"
	"os"
	"reflect"
	"sync"
	"time"

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

// componentInterfaceType returns the component interface type.
// nolint
func (Implements[T]) componentInterfaceType() reflect.Type { return reflect.TypeFor[T]() }

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

// refType returns the referenced component interface type.
// nolint
func (r Ref[T]) refType() reflect.Type { return reflect.TypeFor[T]() }

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

// InstanceOf[T any] is an interface for asserting implementation of an interface T.
type InstanceOf[T any] interface {
	implements(T)
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
		opts.interceptor = interceptor.Chain(interceptors)
	}
}

// WithShutdownTimeout sets how long shutdown hooks may run.
func WithShutdownTimeout(timeout time.Duration) func(*options) {
	return func(opts *options) {
		opts.shutdownTimeout = timeout
	}
}

// MustRun is a helper function to run the application with the provided main component and options.
// It panics if an error occurs during the execution.
func MustRun[T InstanceOf[Main]](ctx context.Context, run func(context.Context, T) error, opts ...func(*options)) {
	lo.Must0(Run(ctx, run, opts...))
}

// Run initializes and runs the application with the provided main component and options.
func Run[T InstanceOf[Main]](ctx context.Context, run func(context.Context, T) error, opts ...func(*options)) error {
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
	main, err := kod.Get[T](ctx)
	if err != nil {
		return err
	}

	// wait for shutdown signal
	sig := make(chan os.Signal, 2)
	signals.Shutdown(ctx, sig, func(_ bool) {
		cancel()
	})

	// wait for stop signal
	err = run(ctx, main)

	ctx, timeoutCancel := context.WithTimeout(
		context.WithoutCancel(ctx), kod.shutdownTimeout)
	defer timeoutCancel()

	// run hook functions
	kod.hooker.Do(ctx)

	return err
}

// Kod represents the core structure of the application, holding configuration and component registrations.
type Kod struct {
	mu *sync.Mutex

	shutdownTimeout time.Duration
	hooker          *hooks.Hooker

	regs                []*Registration
	registryByName      map[string]*Registration
	registryByInterface map[reflect.Type]*Registration
	registryByImpl      map[reflect.Type]*Registration

	components         map[string]any
	impls              map[string]any
	lazyInitComponents map[reflect.Type]bool
	opts               *options
}

// options defines the configuration options for Kod.
type options struct {
	fakes           map[reflect.Type]any
	registrations   []*Registration
	interceptor     interceptor.Interceptor
	shutdownTimeout time.Duration
}

// newKod creates a new instance of Kod with the provided registrations and options.
func newKod(_ context.Context, opts ...func(*options)) (*Kod, error) {
	opt := &options{}
	for _, o := range opts {
		o(opt)
	}

	shutdownTimeout := opt.shutdownTimeout
	if shutdownTimeout <= 0 {
		shutdownTimeout = 5 * time.Second
	}

	kod := &Kod{
		mu:                  &sync.Mutex{},
		shutdownTimeout:     shutdownTimeout,
		hooker:              hooks.New(),
		regs:                registry.All(),
		registryByName:      make(map[string]*Registration),
		registryByInterface: make(map[reflect.Type]*Registration),
		registryByImpl:      make(map[reflect.Type]*Registration),
		components:          make(map[string]any),
		impls:               make(map[string]any),
		opts:                opt,
	}

	kod.register(opt.registrations)

	var err error
	kod.lazyInitComponents, err = processRegistrations(kod.regs)
	if err != nil {
		return nil, err
	}

	if err := checkCircularDependency(kod.regs); err != nil {
		return nil, err
	}

	return kod, nil
}

// SetDefaultInterceptor sets the default interceptor for the Kod instance.
func (k *Kod) SetInterceptors(interceptors ...interceptor.Interceptor) {
	k.opts.interceptor = interceptor.Chain(interceptors)
}

// Defer adds a hook function to the Kod instance.
func (k *Kod) Defer(name string, fn func(context.Context) error) {
	k.hooker.Add(hooks.HookFunc{Name: name, Fn: fn})
}

// register adds the given implementations to the Kod instance.
func (k *Kod) register(regs []*Registration) {
	if len(regs) > 0 {
		k.regs = regs
	}

	for _, v := range k.regs {
		if v == nil {
			continue
		}
		k.registryByName[v.Name] = v
		k.registryByInterface[v.Interface] = v
		k.registryByImpl[v.Impl] = v
	}
}
