package hlog

import (
	"errors"
	"io"
	"os"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Test_convertToLogrusFields(t *testing.T) {
	tests := []struct {
		name   string
		arg    interface{}
		expect logrus.Fields
	}{
		{
			name: "map",
			arg: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
			expect: logrus.Fields{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name: "struct",
			arg: struct {
				key1 string
				key2 string
			}{
				key1: "value1",
				key2: "value2",
			},
			expect: logrus.Fields{
				"key1": "value1",
				"key2": "value2",
			},
		},
		{
			name:   "integer",
			arg:    1,
			expect: logrus.Fields{},
		},
		{
			name:   "nil",
			arg:    nil,
			expect: logrus.Fields{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := convertToLogrusFields(test.arg)
			assert.True(t, reflect.DeepEqual(actual, test.expect), "result is not match, expect=%v, actual=%v", test.expect, actual)
		})
	}
}

func Test_Logger(t *testing.T) {
	log, hook := test.NewNullLogger()

	// change log level
	log.SetLevel(logrus.TraceLevel)

	target := loggerImpl{
		logrusLogger: log,
	}

	msg := "test message"

	// Trace
	target.Trace(msg, nil)
	assert.Equal(t, logrus.TraceLevel, hook.LastEntry().Level)
	assert.Equal(t, msg, hook.LastEntry().Message)

	// Debug
	target.Debug(msg, nil)
	assert.Equal(t, logrus.DebugLevel, hook.LastEntry().Level)
	assert.Equal(t, msg, hook.LastEntry().Message)

	// Info
	target.Info(msg, nil)
	assert.Equal(t, logrus.InfoLevel, hook.LastEntry().Level)
	assert.Equal(t, msg, hook.LastEntry().Message)

	// Warn
	target.Warn(msg, nil)
	assert.Equal(t, logrus.WarnLevel, hook.LastEntry().Level)
	assert.Equal(t, msg, hook.LastEntry().Message)

	// Error
	target.WithError(errors.New("test error")).Error(msg, nil)
	assert.Equal(t, logrus.ErrorLevel, hook.LastEntry().Level)
	assert.Equal(t, msg, hook.LastEntry().Message)
}

func Test_getLogrusLogLevel(t *testing.T) {
	tests := []struct {
		name   string
		arg    LogLevel
		expect logrus.Level
	}{
		{
			name:   "fatal",
			arg:    Fatal,
			expect: logrus.FatalLevel,
		},
		{
			name:   "error",
			arg:    Error,
			expect: logrus.ErrorLevel,
		},
		{
			name:   "warn",
			arg:    Warn,
			expect: logrus.WarnLevel,
		},
		{
			name:   "debug",
			arg:    Debug,
			expect: logrus.DebugLevel,
		},
		{
			name:   "trace",
			arg:    Trace,
			expect: logrus.TraceLevel,
		},
		{
			name:   "info",
			arg:    Info,
			expect: logrus.InfoLevel,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := getLogrusLogLevel(test.arg)
			assert.Equal(t, test.expect, actual, "expect=%v, actual=%v", test.expect, actual)
		})
	}
}

func Test_NewLogger(t *testing.T) {
	tests := []struct {
		name   string
		arg    Config
		expect Logger
	}{
		{
			name: "only stdout",
			arg: Config{
				OutputFilePath: "",
				LogLevel:       Debug,
			},
			expect: newTestLoggerImpl(os.Stdout, logrus.DebugLevel),
		},
		{
			name: "only file",
			arg: Config{
				OutputFilePath: "test.log",
				LogLevel:       Debug,
			},
			expect: newTestLoggerImpl(
				&lumberjack.Logger{Filename: "test.log"},
				logrus.DebugLevel,
			),
		},
		{
			name: "file and stdout",
			arg: Config{
				OutputFilePath: "test.log",
				LogLevel:       Debug,
				Stdout:         true,
			},
			expect: newTestLoggerImpl(
				io.MultiWriter(&lumberjack.Logger{Filename: "test.log"}, os.Stdout),
				logrus.DebugLevel,
			),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert.NotNil(t, NewLogger(test.arg))
		})
	}
}

func newTestLoggerImpl(out io.Writer, level logrus.Level) Logger {
	logrusLogger := logrus.New()
	logrusLogger.SetOutput(out)
	logrusLogger.SetLevel(level)
	return &loggerImpl{
		logrusLogger: logrusLogger,
	}
}
