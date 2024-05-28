package config

type LogConfig = []LogC

type LogC struct {
	Filename   string `mapstructure:"filename"`    // 日志文件名
	MaxSize    int    `mapstructure:"max_size"`    // 日志文件最大大小，以 MB 为单位
	MaxBackups int    `mapstructure:"max_backups"` // 保留日志文件数量
	MaxAge     int    `mapstructure:"max_age"`     // 日志文件最大保留天数
	Compress   bool   `mapstructure:"compress"`    // 是否压缩
	LevelRange string `mapstructure:"level_range"` // 日志级别，[info,error)
}
