package hlog

import (
	"errors"
	"reflect"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
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
