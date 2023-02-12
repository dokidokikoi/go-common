package config

import "fmt"

type RabbitMqConfig struct {
	Host     string
	Port     int
	Username string `validate:"required"`
	Password string
}

func (r RabbitMqConfig) Dns() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d",
		r.Username,
		r.Password,
		r.Host,
		r.Port)
}

const (
	rabbitMqKey string = "rabbitmq"
)

func GetRabbitMqInfo() *RabbitMqConfig {
	rabbitMqConfig := &RabbitMqConfig{
		Port: 5672, Host: "127.0.0.1", Username: "harukaze", Password: "123456",
	}
	conf := GetSpecConfig(rabbitMqKey)
	if conf != nil {
		conf.Unmarshal(rabbitMqConfig)
	}

	return rabbitMqConfig
}
