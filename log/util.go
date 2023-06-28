package log

import "github.com/sirupsen/logrus"

var enableInfo = true
var enableError = true
var enableDebug = true

func Info(format string, args ...any) {
	if enableInfo {
		logrus.Infof(format, args...)
	}
}

func Error(format string, args ...any) {
	if enableError {
		logrus.Errorf(format, args...)
	}
}

func Debug(format string, args ...any) {
	if enableDebug {
		logrus.Debugf(format, args...)
	}
}
