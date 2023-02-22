package config

import "fmt"

type EsConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Cert     string `mapstructure:"cert"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func (s EsConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (s EsConfig) Url() string {
	return fmt.Sprintf("https://%s", s.Address())
}
