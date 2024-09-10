package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var Global *config

func Load(cfgFile string) error {
	viper.SetConfigFile(cfgFile)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return errors.Wrap(err, "failed to read config file")
	}

	if err := viper.Unmarshal(&Global); err != nil {
		return errors.Wrap(err, "failed to unmarshal config")
	}

	return nil
}

type config struct {
	Gen struct {
		Model struct {
			DSN     string                       `mapstructure:"dsn"`
			Schema  string                       `mapstructure:"schema"`
			Source  string                       `mapstructure:"source"` // postgres, mysql, cockroachdb, mariadb or sqlite
			Path    string                       `mapstructure:"path"`
			Types   map[string]map[string]string `mapstructure:"types"`
			Ignores []string                     `mapstructure:"ignores"` // ignore tables
		} `mapstructure:"model"`
	} `mapstructure:"gen"`
}
