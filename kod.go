package kod

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"reflect"
	"sync"
	"time"

	"github.com/go-kod/kod/internal/hooks"
	"github.com/go-kod/kod/internal/kslog"
	"github.com/go-kod/kod/internal/paths"
	"github.com/go-kod/kod/internal/reflects"
	"github.com/go-kod/kod/internal/registry"
	"github.com/go-kod/kod/internal/signals"
	"github.com/samber/lo"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
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

// Run initializes and runs the application with the provided main component and options.
func Run[T any, _ PointerToMain[T]](ctx context.Context, run func(context.Context, *T) error, opts ...func(*options)) error {
	opt := &options{}
	for _, o := range opts {
		o(opt)
	}

	// Create a new Kod instance.
	kod, err := newKod(*opt)
	if err != nil {
		return err
	}

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
		kod.log.Info("shutdown ...")
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

// logConfig defines the configuration for logging.
type logConfig struct {
	Level string
	File  string
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

	viper       *viper.Viper
	log         *slog.Logger
	logLevelVar *slog.LevelVar

	hooker *hooks.Hooker

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
	logWrapper     func(slog.Handler) slog.Handler
	registrations  []*Registration
}

// newKod creates a new instance of Kod with the provided registrations and options.
func newKod(opts options) (*Kod, error) {

	kod := &Kod{
		mu: &sync.Mutex{},
		config: kodConfig{
			Name:            filepath.Base(lo.Must(os.Executable())),
			Env:             "local",
			Log:             logConfig{Level: "info"},
			ShutdownTimeout: 5 * time.Second,
		},
		hooker:          hooks.New(),
		regs:            registry.All(),
		registryByName:  make(map[string]*Registration),
		registryByIface: make(map[reflect.Type]*Registration),
		registryByImpl:  make(map[reflect.Type]*Registration),
		components:      make(map[string]any),
		opts:            opts,
	}

	kod.register(opts.registrations)

	if err := kod.parseConfig(opts.configFilename); err != nil {
		return nil, err
	}

	if err := registry.ValidateRegistrations(kod.regs); err != nil {
		return nil, err
	}

	if err := registry.CheckCircularDependency(kod.regs); err != nil {
		return nil, err
	}

	kod.initLog()

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

// LevelVar returns the log level variable of the Kod instance.
func (k *Kod) LevelVar() *slog.LevelVar {
	return k.logLevelVar
}

// register adds the given implementations to the Kod instance.
func (k *Kod) register(regs []*Registration) {
	if len(regs) > 0 {
		k.regs = regs
	}

	for _, v := range k.regs {
		k.registryByName[v.Name] = v
		k.registryByIface[v.Iface] = v
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

func (k *Kod) initLog() {

	k.logLevelVar = new(slog.LevelVar)
	lo.Must0(k.logLevelVar.UnmarshalText([]byte(k.config.Log.Level)))

	// Default to stdout.
	var writer io.Writer = os.Stdout
	// If a log file is specified, use it.
	if k.config.Log.File != "" {
		logger := &lumberjack.Logger{
			Filename:   k.config.Log.File,
			MaxSize:    500, // megabytes
			MaxBackups: 7,
			MaxAge:     28, //days
			Compress:   false,
		}
		k.hooker.Add(hooks.HookFunc{
			Name: PkgPath,
			Fn: func(ctx context.Context) error {
				return logger.Close()
			},
		})
		writer = logger
	}

	jsonHandler := slog.NewJSONHandler(
		writer, &slog.HandlerOptions{
			AddSource: true,
			Level:     k.logLevelVar,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				// Remove the directory from the source's filename.
				if a.Key == slog.SourceKey {
					source := a.Value.Any().(*slog.Source)
					source.File = paths.CustomBase(source.File, 2)
					source.Function = paths.CustomBase(source.Function, 1)
				}

				return a
			},
		},
	)

	var handler slog.Handler
	if k.opts.logWrapper != nil {
		handler = kslog.NewOtelHandler(k.opts.logWrapper(jsonHandler))
	} else {
		handler = kslog.NewOtelHandler(jsonHandler)
	}

	k.log = slog.New(handler)
}
