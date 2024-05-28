package zaplog

import (
	"os"

	"github.com/dokidokikoi/go-common/config"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap/zapcore"
)

var (
	defaultFilename   = "./log/log.log"
	defaultMaxSize    = 500
	defaultMaxBackups = 3
	defaultMaxAge     = 90
)

func NewStdCore() zapcore.Core {
	consoleWriter := zapcore.Lock(os.Stdout)
	consoleEncoder := zapcore.NewConsoleEncoder(NewConsoleEncoderConfig())
	return zapcore.NewCore(consoleEncoder, consoleWriter, infoPriority)
}

func NewFileCore(conf config.LogConfig) []zapcore.Core {
	var cores []zapcore.Core
	conf = defaultOption(conf)
	for _, c := range conf {
		writer := &lumberjack.Logger{
			Filename:   c.Filename,
			MaxSize:    c.MaxSize,
			MaxBackups: c.MaxBackups,
			MaxAge:     c.MaxAge,
			Compress:   c.Compress,
		}
		cores = append(cores, zapcore.NewCore(
			zapcore.NewJSONEncoder(NewFileEncoderConfig()),
			zapcore.Lock(zapcore.AddSync(writer)),
			getPriority(c.LevelRange),
		))
	}

	return cores

}

func defaultOption(conf config.LogConfig) config.LogConfig {
	if len(conf) < 1 {
		conf = append(conf, config.LogC{
			Filename:   defaultFilename,
			MaxSize:    defaultMaxSize,
			MaxBackups: defaultMaxBackups,
			MaxAge:     defaultMaxAge,
		})
	}
	for i := range conf {
		if conf[i].Filename == "" {
			conf[i].Filename = defaultFilename
		}
		if conf[i].MaxSize <= 0 {
			conf[i].MaxSize = defaultMaxSize
		}
		if conf[i].MaxBackups <= 0 {
			conf[i].MaxBackups = defaultMaxBackups
		}
		if conf[i].MaxAge <= 0 {
			conf[i].MaxAge = defaultMaxAge
		}
	}
	return conf
}
