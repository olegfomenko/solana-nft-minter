package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
)

const configFileEnvKey = "CONFIG"

func ConfigureService() {
	viper.SetConfigFile(os.Getenv(configFileEnvKey))
	err := viper.ReadInConfig()

	if err != nil {
		panic(errors.Wrap(err, "error reading config"))
	}
}
