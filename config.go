package kod

import (
	"fmt"
	"io/fs"
	"os"
	"reflect"
	"strings"

	"github.com/spf13/viper"
)

// config calls the WithConfig.config method on the provided value and returns
// the result. If the provided value doesn't have a WithConfig.config method,
// config returns nil.
//
// config panics if the provided value is not a pointer to a struct.
func config(v reflect.Value) any {
	if v.Kind() != reflect.Pointer || v.Elem().Kind() != reflect.Struct {
		panic(fmt.Errorf("invalid non pointer to struct value: %v", v))
	}
	s := v.Elem()
	t := s.Type()
	for i := 0; i < t.NumField(); i++ {
		// Check that f is an embedded field of type kod.WithConfig[T].
		f := t.Field(i)
		if !f.Anonymous ||
			f.Type.PkgPath() != PkgPath ||
			!strings.HasPrefix(f.Type.Name(), "WithConfig[") {
			continue
		}

		// Call the Config method to get a *T.
		config := s.Field(i).Addr().MethodByName("Config")
		return config.Call(nil)[0].Interface()
	}
	return nil
}

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
