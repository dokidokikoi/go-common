package zaplog

import (
	"github.com/dokidokikoi/go-common/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var z *zap.Logger

func init() {
	z = NewDefaultLogger()
}

// NewDefaultLogger 获取默认logger
// 提供两种输出文件以及标准输出
func NewDefaultLogger() *zap.Logger {
	core := zapcore.NewTee(
		append(NewFileCore(config.LogConfig{}), NewStdCore())...,
	)
	return zap.New(core)
}

func NewLogger(conf config.LogConfig) *zap.Logger {
	core := zapcore.NewTee(
		append(NewFileCore(conf), NewStdCore())...,
	)
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.WarnLevel))
}

func SetLogger(conf config.LogConfig) {
	z = NewLogger(conf)
}

func Sugar() *zap.SugaredLogger {
	return z.Sugar()
}

func L() *zap.Logger {
	return z
}
