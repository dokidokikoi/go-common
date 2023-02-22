package config

import (
	"database/sql/driver"
	"encoding/json"
)

type MongoConfig struct {
	Uri      string `mapstructure:"uri"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Option   *MongoOption
}

type MongoOption struct {
	MaxPoolSize int
	MinPoolSize int
}

func (s MongoConfig) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *MongoConfig) Scan(src interface{}) (err error) {
	var c MongoConfig
	err = json.Unmarshal(src.([]byte), &c)
	if err != nil {
		return
	}
	*s = c
	return nil
}
