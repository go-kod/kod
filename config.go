package kod

import (
	"fmt"
	"io/fs"
	"os"

	"github.com/spf13/viper"
)

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
