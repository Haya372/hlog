package hlog

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type LogLevel int

const (
	Fatal LogLevel = iota
	Error
	Warn
	Info
	Debug
	Trace
)

type Config struct {
	LogLevel       LogLevel `yaml:"logLevel"`
	OutputFilePath string   `yaml:"outputFilePath"`
	MaxAge         int      `yaml:"maxAge"`
	MaxBackups     int      `yaml:"maxBackups"`
	MaxSize        int      `yaml:"maxSize"`
	Compress       bool
	Stdout         bool
}

var (
	filePathEmptyErr          = errors.New("configuration file path is empty")
	fileReadingErr            = errors.New("could not read configuration file")
	unrecognizedLogLevelError = errors.New("unrecognized log level")
)

func NewConfigFromYamlFile(filePath string) (*Config, error) {
	if len(filePath) == 0 {
		return nil, filePathEmptyErr
	}

	file, err := os.Open(filePath)

	if err != nil {
		return nil, errors.Join(fileReadingErr, err)
	}

	decoder := yaml.NewDecoder(file)

	conf := new(Config)

	err = decoder.Decode(conf)
	if err != nil {
		return nil, errors.Join(fileReadingErr, err)
	}

	return conf, nil
}

func (l *LogLevel) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	e := unmarshal(&s)

	if e != nil {
		return errors.Join(fileReadingErr, e)
	}

	switch strings.ToLower(s) {
	case "fatal":
		*l = Fatal
	case "error":
		*l = Error
	case "warn":
		*l = Warn
	case "info":
		*l = Info
	case "debug":
		*l = Debug
	case "trace":
		*l = Trace
	default:
		return fmt.Errorf("log level %s is not supported", s)
	}

	return nil
}

func (l LogLevel) String() string {
	strs := []string{
		"Fatal", "Error", "Warn", "Info", "Debug", "Trace",
	}
	return strs[l]
}
