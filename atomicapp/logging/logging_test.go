package logging

import (
	"testing"

	"github.com/Sirupsen/logrus"
)

func TestLoggingLevel(t *testing.T) {
	expectedLevel := logrus.ErrorLevel
	if level := getLevel(2); level != expectedLevel {
		logrus.Fatalf("Passed in level 4, expected logrus.ErrorLevel, got %v", level)
	}
}
