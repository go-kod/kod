package kod

import (
	"context"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/samber/lo"
	"github.com/spf13/viper"
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
func (i *Implements[T]) L() *slog.Logger {
	return i.log
}

// setLogger sets the logger for the component.
func (i *Implements[T]) setLogger(log *slog.Logger) {
	i.log = log
}

// implements is a marker method to assert implementation of an interface.
// nolint
func (Implements[T]) implements(T) {}

// Ref[T any] is a reference holder to a value of type T.
type Ref[T any] struct {
	value T
}

// Get returns the held reference value.
func (r Ref[T]) Get() T { return r.value }

// isRef is a marker method to identify a Ref type.
func (r Ref[T]) isRef() {}

// setRef sets the reference value.
func (r *Ref[T]) setRef(val any) { r.value = val.(T) }

// Main is the interface that should be implemented by an application's main component.
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
type WithConfig[T any] struct {
	config T
}

// Config returns a pointer to the config.
func (wc *WithConfig[T]) Config() *T {
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

// Run initializes and runs the application with the provided main component and options.
func Run[T any, _ PointerToMain[T]](ctx context.Context, app func(context.Context, *T) error, opts ...func(*options)) error {
	opt := &options{}
	for _, o := range opts {
		o(opt)
	}

	kod, err := newKod(getRegs(), *opt)
	if err != nil {
		return err
	}

	ctx = newContext(ctx, kod)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	main, err := kod.getImpl(ctx, rtype[T]())
	if err != nil {
		return err
	}

	stop := make(chan struct{}, 2)
	shutdown(ctx, func(grace bool) {
		kod.log.Info("shutdown ...")
		cancel()
		stop <- struct{}{}
	})

	go func() {
		err = app(ctx, main.(*T))
		stop <- struct{}{}
	}()

	<-stop
	kod.close(ctx)

	return err
}

// logConfig defines the configuration for logging.
type logConfig struct {
	Level string
}

// kodConfig defines the overall configuration for the Kod application.
type kodConfig struct {
	Name    string
	Env     string
	Version string

	Log             logConfig
	ShutdownTimeout time.Duration
}

// Kod represents the core structure of the application, holding configuration and component registrations.
type Kod struct {
	mu *sync.Mutex

	config kodConfig

	log   *slog.Logger
	viper *viper.Viper

	deferMux sync.Mutex
	defers   []deferFn

	regs            []*Registration
	registryByName  map[string]*Registration
	registryByIface map[reflect.Type]*Registration
	registryByImpl  map[reflect.Type]*Registration

	components map[string]any
	opts       options
}

// options defines the configuration options for Kod.
type options struct {
	configFilename string
	fakes          map[reflect.Type]any
}

// newKod creates a new instance of Kod with the provided registrations and options.
func newKod(regs []*Registration, opts options) (*Kod, error) {

	kod := &Kod{
		mu: &sync.Mutex{},
		config: kodConfig{
			Name:            filepath.Base(lo.Must(os.Executable())),
			Env:             "local",
			Log:             logConfig{Level: "info"},
			ShutdownTimeout: 5 * time.Second,
		},
		regs:            make([]*Registration, 0),
		registryByName:  make(map[string]*Registration),
		registryByIface: make(map[reflect.Type]*Registration),
		registryByImpl:  make(map[reflect.Type]*Registration),
		components:      make(map[string]any),
		opts:            opts,
	}

	kod.register(regs)

	if err := kod.parseConfig(opts.configFilename); err != nil {
		return nil, err
	}

	if err := kod.validateRegistrations(); err != nil {
		return nil, err
	}

	if err := kod.checkCircularDependency(); err != nil {
		return nil, err
	}

	kod.initLog()

	return kod, nil
}

// register adds the given implementations to the Kod instance.
func (k *Kod) register(impl []*Registration) {
	k.regs = append(k.regs, impl...)

	for _, v := range impl {
		k.registryByName[v.Name] = v
		k.registryByIface[v.Iface] = v
		k.registryByImpl[v.Impl] = v
	}
}

// close gracefully shuts down the Kod application.
func (k *Kod) close(ctx context.Context) {
	ctx, timeoutCancel := context.WithTimeout(context.WithoutCancel(ctx), k.config.ShutdownTimeout)
	defer timeoutCancel()

	err := k.runDefer(ctx)
	if err != nil {
		k.log.Error("runDefer failed", "error", err)
	}
}

// Config returns the current configuration of the Kod instance.
func (k *Kod) Config() kodConfig {
	return k.config
}
