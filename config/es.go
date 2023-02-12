package config

import "fmt"

type EsConfig struct {
	Host     string
	Port     int
	Cert     string `validate:"required"`
	Username string `validate:"required"`
	Password string
}

func (s EsConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (s EsConfig) Url() string {
	return fmt.Sprintf("https://%s", s.Address())
}

const (
	esKey string = "elasticsearch"
)

func GetEsInfo() *EsConfig {
	esConfig := &EsConfig{
		Port: 9200, Host: "127.0.0.1", Cert: "http_ca.crt", Username: "elastic", Password: "",
	}
	conf := GetSpecConfig(esKey)
	if conf != nil {
		conf.Unmarshal(esConfig)
	}

	return esConfig
}
