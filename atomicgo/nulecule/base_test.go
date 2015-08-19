package nulecule

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
)

var base = Base{}

func TestSetTargetPath(t *testing.T) {
	var (
		target string
		err    error
	)
	target = "junk/path/doesn't/exist"

	base.setTargetPath(target)
	if err = base.setTargetPath(target); err != nil {
		t.Fatal("Failed to set target path")
	}

	cwd, _ := os.Getwd()
	if base.Target() != cwd {
		logrus.Errorf("Gave target path %s and expected CWD, got %s", target, base.Target())
		t.Fatal("Target path set incorrectly")
	}

	target = "/"
	base.setTargetPath(target)
	if err = base.setTargetPath(target); err != nil {
		t.Fatal("Failed to set target path")
	}

	if base.Target() != target {
		logrus.Errorf("Gave target path %s and expected /, got %s", target, base.Target())
		t.Fatal("Target path set incorrectly")
	}
}
