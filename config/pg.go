package config

import "fmt"

type PGConfig struct {
	Host     string
	Port     int
	Database string `validate:"required"`
	Username string `validate:"required"`
	TimeZone string
	Password string
}

func (p PGConfig) Dns() string {
	return fmt.Sprintf(`host=%s user=%s dbname=%s port=%d sslmode=disable TimeZone=%s password=%s`,
		p.Host,
		p.Username,
		p.Database,
		p.Port,
		p.TimeZone,
		p.Password,
	)
}

const (
	pgKey string = "postgresql"
)

func GetPgInfo() *PGConfig {
	pgConfig := &PGConfig{
		Port: 5432, Host: "127.0.0.1", Database: "postgres", TimeZone: "Asia/Shanghai", Password: "postgres",
	}
	conf := GetSpecConfig(pgKey)
	if conf != nil {
		conf.Unmarshal(pgConfig)
	}

	return pgConfig
}
