package log

import (
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
	Log.Formatter = &logrus.JSONFormatter{}
	Log.SetReportCaller(true)
	Log.Out = &lumberjack.Logger{
		Filename:   "logs/plumber.log",
		MaxSize:    10, // megabytes
		MaxBackups: 5,
		MaxAge:     30, //days
	}
	Log.Infof("log init success")
}
func Info(args ...interface{}) {
	Log.Info(args...)
}
func Infof(format string, args ...interface{}) {
	Log.Infof(format, args...)
}

func Error(args ...interface{}) {
	Log.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	Log.Errorf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	Log.Debugf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	Log.Panicf(format, args...)
}
