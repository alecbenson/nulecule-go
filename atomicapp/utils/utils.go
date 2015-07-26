package utils

import (
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Sirupsen/logrus"
	"gopkg.in/yaml.v1"
)

const (
	FILE_PREFIX   = "file://"
	DOCKER_PREFIX = "docker://"
)

//PathExists verifies that the given path points to a valid file or directory
func PathExists(path string) bool {
	if path == "" {
		return false
	}

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			logrus.Errorf("Failed to read file: %v", err)
			return false
		}
	}
	return true
}

//PathIsFile returns true if the path points to a vald file
func PathIsFile(path string) bool {
	if path == "" {
		return false
	}

	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode().IsRegular()
}

//PathIsDirectory returns true if the path points to a valid directory
func PathIsDirectory(path string) bool {
	if path == "" {
		return false
	}

	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.IsDir()
}

//SanitizePath Strips the protocol from the provided path
func SanitizePath(url string) string {
	if strings.HasPrefix(url, FILE_PREFIX) {
		return strings.TrimPrefix(url, FILE_PREFIX)
	}
	return url
}

//ParseYAMLFile will read in a path to a yaml file and parse to a map
func ParseYAMLFile(filepath string, result interface{}) error {

	//Check for an instantiated map
	if result == nil {
		return errors.New("Uninstantiated map proved")
	}

	//Check for valid path
	if !PathExists(filepath) {
		return errors.New("File does not exist")
	}

	//Read the file
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	//Attempt to parse
	err = yaml.Unmarshal(file, result)
	if err != nil {
		return err
	}
	return nil
}
