package hlog

import (
	"io"
	"os"
	"reflect"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Logger interface {
	Trace(msg string, v interface{})
	Debug(msg string, v interface{})
	Info(msg string, v interface{})
	Warn(msg string, v interface{})
	Error(msg string, v interface{})
	Fatal(msg string, v interface{})
	WithError(err error) Logger
}

type loggerImpl struct {
	logrusLogger *logrus.Logger
}

func (l *loggerImpl) Trace(msg string, v interface{}) {
	l.logrusLogger.WithFields(convertToLogrusFields(v)).Trace(msg)
}

func (l *loggerImpl) Debug(msg string, v interface{}) {
	l.logrusLogger.WithFields(convertToLogrusFields(v)).Debug(msg)
}

func (l *loggerImpl) Info(msg string, v interface{}) {
	l.logrusLogger.WithFields(convertToLogrusFields(v)).Info(msg)
}

func (l *loggerImpl) Warn(msg string, v interface{}) {
	l.logrusLogger.WithFields(convertToLogrusFields(v)).Warn(msg)
}

func (l *loggerImpl) Error(msg string, v interface{}) {
	l.logrusLogger.WithFields(convertToLogrusFields(v)).Error(msg)
}

func (l *loggerImpl) Fatal(msg string, v interface{}) {
	l.logrusLogger.WithFields(convertToLogrusFields(v)).Fatal(msg)
}

func (l *loggerImpl) WithError(err error) Logger {
	l.logrusLogger.WithError(err)
	return l
}

func convertToLogrusFields(v interface{}) logrus.Fields {
	fields := make(logrus.Fields)
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map:
		iter := rv.MapRange()
		for iter.Next() {
			fields[iter.Key().String()] = iter.Value().String()
		}
	case reflect.Struct:
		for i := 0; i < rv.NumField(); i++ {
			name := rv.Type().Field(i).Name
			value := rv.FieldByName(name)
			fields[name] = value.String()
		}
	}
	return fields
}

func NewLogger(config Config) Logger {
	var writer io.Writer
	if len(config.OutputFilePath) == 0 {
		writer = os.Stdout
	} else {
		lumberjackWriter := &lumberjack.Logger{
			Filename:   config.OutputFilePath,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}
		if config.Stdout {
			writer = io.MultiWriter(lumberjackWriter, os.Stdout)
		} else {
			writer = lumberjackWriter
		}
	}

	logrusLogger := logrus.New()

	logrusLogger.SetOutput(writer)
	logrusLogger.SetLevel(getLogrusLogLevel(config.LogLevel))

	return &loggerImpl{
		logrusLogger: logrusLogger,
	}
}

func getLogrusLogLevel(logLevel LogLevel) logrus.Level {
	var res logrus.Level
	switch logLevel {
	case Fatal:
		res = logrus.FatalLevel
	case Error:
		res = logrus.ErrorLevel
	case Warn:
		res = logrus.WarnLevel
	case Debug:
		res = logrus.DebugLevel
	case Trace:
		res = logrus.TraceLevel
	default:
		res = logrus.InfoLevel
	}
	return res
}
