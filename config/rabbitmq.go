package config

import "fmt"

type RabbitMqConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
}

func (r RabbitMqConfig) Dns() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d",
		r.Username,
		r.Password,
		r.Host,
		r.Port)
}
