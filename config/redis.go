package config

import (
	"database/sql/driver"
	"encoding/json"
)

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

func (s RedisConfig) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *RedisConfig) Scan(src interface{}) (err error) {
	var c RedisConfig
	err = json.Unmarshal(src.([]byte), &c)
	if err != nil {
		return
	}
	*s = c
	return nil
}
