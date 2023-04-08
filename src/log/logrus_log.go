package log

import (
	"io"
	"os"

	"github.com/google/wire"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

func provideLog(logLevel string) Loger {
	logrusLog := logrus.New()
	logrusLog.Formatter = &logrus.TextFormatter{ForceColors: true, FullTimestamp: true}
	logrusLog.SetReportCaller(true)
	jackOut := &lumberjack.Logger{
		Filename:   "logs/plumber.log",
		MaxSize:    10, // megabytes
		MaxBackups: 5,
		MaxAge:     30, //days
	}
	logrusLog.SetOutput(io.MultiWriter(jackOut, os.Stdout))
	switch logLevel {
	case "debug":
		logrusLog.SetLevel(logrus.DebugLevel)
	case "info":
		logrusLog.SetLevel(logrus.InfoLevel)
	case "warn":
		logrusLog.SetLevel(logrus.WarnLevel)
	case "error":
		logrusLog.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrusLog.SetLevel(logrus.FatalLevel)
	case "trace":
		logrusLog.SetLevel(logrus.TraceLevel)
	default:
		logrusLog.SetLevel(logrus.InfoLevel)
	}
	logrusLog.Debugf("log init success")
	return logrusLog
}

func logLevel() string {
	return os.Getenv("LOG_LEVEL")
}

var logSet = wire.NewSet(logLevel, provideLog)
