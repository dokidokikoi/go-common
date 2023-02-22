package config

import "fmt"

type PGConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Database string `mapstructure:"database"`
	Username string `mapstructure:"username"`
	TimeZone string `mapstructure:"timezone"`
	Password string `mapstructure:"password"`
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
