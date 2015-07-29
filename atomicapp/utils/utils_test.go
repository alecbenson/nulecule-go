package utils

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func TestPathExists(t *testing.T) {
	cwd, _ := os.Getwd()
	if !PathExists(cwd) {
		t.Fatal("PathExists returned false for the current working directory, expected true")
	}
}

func TestPathIsDirectory(t *testing.T) {
	if !PathIsDirectory("/") {
		t.Fatal("PathIsDirectory returned false for the root directory, expected true")
	}
}

func TestPathIsFile(t *testing.T) {
	if !PathIsFile("/etc/hosts") {
		t.Fatal("PathIsFile returned false for hosts file directory, expected true")
	}
}

func TestCommandOutput(t *testing.T) {
	expected := "testing123"
	command := exec.Command("echo", expected)
	out, err := CheckCommandOutput(command, true)
	if err != nil {
		t.Fatalf("Got error running command: %s", err)
	}

	if !strings.Contains(string(out), expected) {
		t.Fatalf("Expected testCommandOutput to contain %+v, got %+v", expected, string(out))
	}
}

func TestGenerateUniqueName(t *testing.T) {
	pattern := "^(testing-[a-f-0-9]+)$"
	name, err := GenerateUniqueName("testing")
	if err != nil {
		t.Fatalf("Got an error generating a unique name: %s", err)
	}
	match, err := regexp.MatchString(pattern, name)
	if err != nil {
		t.Fatalf("Error matching regex pattern: %s", err)
	}

	if !match {
		t.Fatalf("Expected to generate a name matching %s. Output name was %s", pattern, name)
	}
}

func TestSanitizePath(t *testing.T) {
	name := "file://aleciswritingtests"
	expected := "aleciswritingtests"

	if SanitizePath(name) != expected {
		t.Fatalf("Expected file to be sanitized to %s, got %s instead", expected, name)
	}
}

func TestGetBaseImageName(t *testing.T) {
	name := "docker.io/fedora:latest"
	expected := "fedora"
	if GetBaseImageName(name) != expected {
		t.Fatalf("Expected base image name to be %s, got %s instead", expected, name)
	}
}
