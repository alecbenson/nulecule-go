package logging

import (
	"github.com/Sirupsen/logrus"
)

func SetLogLevel(level int) {
	if level >= 0 && level <= 5 {
		logLevel := getLevel(level)
		logrus.SetLevel(logLevel)
	}
}

func getLevel(level int) logrus.Level {
	switch level {
	case 5:
		return logrus.DebugLevel
	case 4:
		return logrus.InfoLevel
	case 3:
		return logrus.WarnLevel
	case 2:
		return logrus.ErrorLevel
	case 1:
		return logrus.FatalLevel
	case 0:
		return logrus.PanicLevel
	default:
		return logrus.ErrorLevel
	}
}
