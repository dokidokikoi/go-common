package config

import "fmt"

type ServerConfig struct {
	Host string
	Port int
}

func (s ServerConfig) Address() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

const (
	severKey string = "server"
)

func GetServerInfo() ServerConfig {
	serverConfig := &ServerConfig{
		Host: "127.0.0.1", Port: 8100,
	}
	conf := GetSpecConfig(severKey)
	if conf != nil {
		conf.Unmarshal(serverConfig)
	}

	return *serverConfig
}
