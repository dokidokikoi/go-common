package config

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

var config *viper.Viper

func Config() *viper.Viper {
	return config
}

func Parse(configFile string, obj interface{}) {
	confFileAbs, err := filepath.Abs(configFile)
	if err != nil {
		panic(err)
	}

	filePathStr, filename := filepath.Split(confFileAbs)
	ext := strings.TrimLeft(path.Ext(filename), ".")
	filename = strings.ReplaceAll(filename, "."+ext, "")

	config = viper.New()
	config.AddConfigPath(filePathStr)
	config.SetConfigName(filename)
	config.SetConfigType(ext)
	err = config.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = config.Unmarshal(obj)
	if err != nil {
		panic(err)
	}
}
