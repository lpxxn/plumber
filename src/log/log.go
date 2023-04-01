package log

import (
	"context"
	"os"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var plumberLog *logrus.Logger

func init() {
	plumberLog = logrus.New()
	plumberLog.Formatter = &logrus.JSONFormatter{}
	plumberLog.SetReportCaller(true)
	plumberLog.Out = &lumberjack.Logger{
		Filename:   "logs/plumber.plumberLog",
		MaxSize:    10, // megabytes
		MaxBackups: 5,
		MaxAge:     30, //days
	}
	logLevel := os.Getenv("LOG_LEVEL")
	switch logLevel {
	case "debug":
		plumberLog.SetLevel(logrus.DebugLevel)
	case "info":
		plumberLog.SetLevel(logrus.InfoLevel)
	case "warn":
		plumberLog.SetLevel(logrus.WarnLevel)
	case "error":
		plumberLog.SetLevel(logrus.ErrorLevel)
	case "fatal":
		plumberLog.SetLevel(logrus.FatalLevel)
	default:
		plumberLog.SetLevel(logrus.InfoLevel)
	}
	plumberLog.Debugf("plumberLog init success")
}

type ILog interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Panicf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

func Info(args ...interface{}) {
	plumberLog.Info(args...)
}
func Infof(format string, args ...interface{}) {
	plumberLog.Infof(format, args...)
}

func Error(args ...interface{}) {
	plumberLog.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	plumberLog.Errorf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	plumberLog.Debugf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	plumberLog.Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	plumberLog.Fatalf(format, args...)
}

func For(ctx context.Context) ILog {
	return plumberLog
}
