package hlog

import (
	"errors"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewConfigFromYamlFile(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		expected *Config
		err      error
	}{
		{
			name:     "success to load yaml",
			filePath: "./testdata/validconf.yaml",
			expected: &Config{
				LogLevel:       Debug,
				OutputFilePath: "log/dev.log",
				MaxAge:         14,
				MaxBackups:     3,
				MaxSize:        500,
				Compress:       true,
				Stdout:         true,
			},
			err: nil,
		},
		{
			name:     "file path is empty",
			filePath: "",
			expected: nil,
			err:      filePathEmptyErr,
		},
		{
			name:     "file is not exist",
			filePath: "notexist.yaml",
			expected: nil,
			err:      fileReadingErr,
		},
		{
			name:     "file reading err",
			filePath: "testdata/invalidconf.yaml",
			expected: nil,
			err:      fileReadingErr,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := NewConfigFromYamlFile(test.filePath)
			if test.err == nil {
				assert.True(t, reflect.DeepEqual(&actual, &test.expected), "config data is not expected. expect=%v, actual=%v\n", test.expected, actual)
				assert.Nil(t, err, "error should be nil, but got %s\n", err)
			} else {
				assert.Nil(t, actual, "config data should be nil, but got %v\n", actual)
				assert.True(t, errors.Is(err, test.err), "error is not expected, expect=%s, actual=%s\n", test.expected, actual)
			}
		})
	}
}

func Test_LogLevel_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name      string
		unmarshal func(interface{}) error
		expect    LogLevel
		isErr     bool
	}{
		{
			name: "fatal",
			unmarshal: func(i interface{}) error {
				*i.(*string) = "Fatal"
				return nil
			},
			expect: Fatal,
			isErr:  false,
		},
		{
			name: "error",
			unmarshal: func(i interface{}) error {
				*i.(*string) = "Error"
				return nil
			},
			expect: Error,
			isErr:  false,
		},
		{
			name: "warn",
			unmarshal: func(i interface{}) error {
				*i.(*string) = "Warn"
				return nil
			},
			expect: Warn,
			isErr:  false,
		},
		{
			name: "info",
			unmarshal: func(i interface{}) error {
				*i.(*string) = "info"
				return nil
			},
			expect: Info,
			isErr:  false,
		},
		{
			name: "debug",
			unmarshal: func(i interface{}) error {
				*i.(*string) = "Debug"
				return nil
			},
			expect: Debug,
			isErr:  false,
		},
		{
			name: "trace",
			unmarshal: func(i interface{}) error {
				*i.(*string) = "Trace"
				return nil
			},
			expect: Trace,
			isErr:  false,
		},
		{
			name: "parse error",
			unmarshal: func(i interface{}) error {
				*i.(*string) = "aaa"
				return nil
			},
			expect: Trace,
			isErr:  true,
		},
		{
			name: "unmarshal error",
			unmarshal: func(i interface{}) error {
				return errors.New("unmarshal error")
			},
			expect: Trace,
			isErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var target LogLevel
			err := target.UnmarshalYAML(test.unmarshal)

			if test.isErr {
				assert.NotNil(t, err, "error should not be nil")
			} else {
				assert.Equal(t, test.expect, target, "log level is not expected. expect=%v, want=%v", test.expect, target)
			}
		})
	}
}

func Test_LogLevel_String(t *testing.T) {
	m := make(map[LogLevel]string)

	m[Fatal] = "Fatal"
	m[Error] = "Error"
	m[Warn] = "Warn"
	m[Info] = "Info"
	m[Debug] = "Debug"
	m[Trace] = "Trace"

	for l, s := range m {
		assert.Equal(t, s, l.String(), "loglevel %v String() is not expected, got %v", s, l)
	}
}
