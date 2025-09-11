package config

import (
	"io"
	"strings"

	"github.com/spf13/viper"
)

var config *viper.Viper

func Config() *viper.Viper {
	return config
}

func Parse(defaultConfig io.Reader, obj interface{}, envPrefix string) {
	config = viper.New()
	config.SetConfigType("yml")
	// 环境变量
	config.AutomaticEnv()
	config.SetEnvPrefix(envPrefix)
	replacer := strings.NewReplacer(".", "_")
	config.SetEnvKeyReplacer(replacer)
	err := config.ReadConfig(defaultConfig)
	if err != nil {
		panic(err)
	}

	err = config.Unmarshal(obj)
	if err != nil {
		panic(err)
	}
}
