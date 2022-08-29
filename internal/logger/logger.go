package log

import (
	"fmt"

	"go.uber.org/zap"
)

type JaegerLogger struct{}

var Logger *zap.Logger
var JgLogger JaegerLogger

/*
если сделать что то вроде
	func Debug (message string, fields ...zap.Field) {
		Logger.Debug(message, fields...)
	}
	в логгах всегда будет программ logger/logger.go
*/

var Debug func(message string, fields ...zap.Field)
var Info func(message string, fields ...zap.Field)
var Error func(message string, fields ...zap.Field)
var Panic func(message string, fields ...zap.Field)
var Fatal func(message string, fields ...zap.Field)
var Sync func() error

func init() {
	var err error
	Logger, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	Debug = Logger.Debug
	Info = Logger.Info
	Error = Logger.Error
	Panic = Logger.Panic
	Fatal = Logger.Fatal
	Sync = Logger.Sync
}

func (_ *JaegerLogger) Error(msg string) {
	defer Logger.Sync()
	Logger.Error(msg)
}
func (_ *JaegerLogger) Infof(msg string, args ...interface{}) {
	defer Logger.Sync()
	Logger.Info(fmt.Sprintf(msg, args...))
}
