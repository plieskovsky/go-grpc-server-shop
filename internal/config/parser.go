package config

import (
	"github.com/plieskovsky/go-grpc-server-shop/internal/server"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var defaultCfg = Configuration{Server: Servers{Grpc: server.DefaultConfig}}

// MustParse must parse and validate viper config.
func MustParse(cfgFile string) Configuration {
	cfg, err := Parse(cfgFile)
	if err != nil {
		panic(err)
	}
	return cfg
}

// Parse parses and validates viper config.
func Parse(cfgFile string) (Configuration, error) {
	viper.SetConfigFile(cfgFile)

	cfg := defaultCfg

	// Load config from file
	if err := viper.ReadInConfig(); err != nil {
		return cfg, errors.Wrap(err, "failed to read configuration")
	}

	// Expand env variables to loaded config
	expandEnvVariables()

	// Deserialize config to struct
	if err := viper.Unmarshal(&cfg); err != nil {
		return cfg, errors.Wrap(err, "failed to deserialize config")
	}

	return cfg, nil
}

func expandEnvVariables() {
	// Need to expand in this way due to https://github.com/spf13/viper/issues/315
	for _, k := range viper.AllKeys() {
		switch value := viper.Get(k).(type) {
		case string:
			viper.Set(k, os.ExpandEnv(value))
		}
	}
}
