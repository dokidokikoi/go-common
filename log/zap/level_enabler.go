package zaplog

import (
	"fmt"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var levelM = map[string]zapcore.Level{
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
}

var (
	highPriority = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	infoPriority = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl > zapcore.DebugLevel
	})
	debugPriority = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})
)

func getPriority(levelRange string) zap.LevelEnablerFunc {
	levelRange = strings.TrimSpace(levelRange)
	arr := strings.Split(levelRange, ",")
	if len(arr) != 2 || len(arr[0]) < 1 || len(arr[1]) < 1 || (arr[0][0] != '[' && arr[0][0] != '(') || (arr[1][len(arr[1])-1] != ']' && arr[1][len(arr[1])-1] != ')') {
		return debugPriority
	}

	llv, lok := levelM[strings.TrimSpace(arr[0][1:])]
	glv, gok := levelM[strings.TrimSpace(arr[1][:len(arr[1])-1])]

	switch fmt.Sprintf("%c%c", arr[0][0], arr[1][len(arr[1])-1]) {
	case "[]":
		return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return (!lok || lvl >= llv) && (!gok || lvl <= glv)
		})
	case "[)":
		return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return (!lok || lvl >= llv) && (!gok || lvl < glv)
		})
	case "(]":
		return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return (!lok || lvl > llv) && (!gok || lvl <= glv)
		})
	case "()":
		return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return (!lok || lvl > llv) && (!gok || lvl < glv)
		})
	default:
		return debugPriority
	}
}
