package config

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

const (
	redisKey string = "redis"
)

func GetRadisInfo() RedisConfig {
	redisConfig := &RedisConfig{
		Port: 6379, Host: "127.0.0.1", DB: 0, Password: "",
	}
	conf := GetSpecConfig(redisKey)
	if conf != nil {
		conf.Unmarshal(redisConfig)
	}

	return *redisConfig
}
