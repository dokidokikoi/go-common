package config

import (
	"database/sql/driver"
	"encoding/json"
)

type PGConfig struct {
	Host     string
	Port     int
	Database string `validate:"required"`
	Username string `validate:"required"`
	Password string
}

func (s PGConfig) Value() (driver.Value, error) {
	return json.Marshal(s)
}

func (s *PGConfig) Scan(src interface{}) (err error) {
	var c PGConfig
	err = json.Unmarshal(src.([]byte), &c)
	if err != nil {
		return
	}
	*s = c
	return nil
}
