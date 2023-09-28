package hlog

import "github.com/sirupsen/logrus"

func SetFuncBeforeExit(f func()) {
	logrus.RegisterExitHandler(f)
}
