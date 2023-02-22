package config

import (
	"database/sql/driver"
	"encoding/json"
)

type KafkaConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

func (s KafkaConfig) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *KafkaConfig) Scan(src interface{}) (err error) {
	var c KafkaConfig
	err = json.Unmarshal(src.([]byte), &c)
	if err != nil {
		return
	}
	*s = c
	return nil
}
