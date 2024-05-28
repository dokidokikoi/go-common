package config

type AppConfig struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	AppName   string `mapstructure:"app_name"`
	SslEnable bool   `mapstructure:"ssl_enable"`
	SslKey    string `mapstructure:"ssl_key"`
	SslCrt    string `mapstructure:"ssl_crt"`
}
