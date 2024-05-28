package zaplog

import (
	"testing"

	"github.com/dokidokikoi/go-common/config"
)

func TestMulti(t *testing.T) {
	conf := []config.LogC{
		{
			Filename:   "./log/info.log",
			LevelRange: "[info, error)",
		},
		{
			Filename:   "./log/error.log",
			LevelRange: "[error,]",
		},
		{
			Filename: "./log/all.log",
		},
	}
	l := NewLogger(conf)
	l.Debug("debug")
	l.Info("hello")
	l.Error("error")
	l.Panic("panic")
}
