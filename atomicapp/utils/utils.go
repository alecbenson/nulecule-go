package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
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
		}
		logrus.Errorf("Failed to read file: %v", err)
		return false
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

//CheckCommandOutput runs the specified command and checks the output for errors.
//Returns output of the command and the error, if any.
//The quiet boolean can be used to supress logging of output
func CheckCommandOutput(command *exec.Cmd, quiet bool) ([]byte, error) {
	logrus.Debugf("Running command: %v", command.Args)
	out, err := command.CombinedOutput()
	if err != nil {
		logrus.Errorf("Failed to run command %s: %s", command.Args, command.Stderr)
		return []byte{}, err
	}

	if out != nil && !quiet {
		logrus.Infof(string(out))
	}
	return out, nil
}

//GenerateUniqueName generates a unique name with the given prefix.
//The name is in the format prefix-xxxxxx where 'xxxxxx' is
//a short string of hex encoded characters
func GenerateUniqueName(prefix string) (string, error) {
	//Generate a random hash to append to the container name
	hashSize := 6
	r := make([]byte, hashSize)
	_, err := rand.Read(r)
	if err != nil {
		logrus.Errorf("Failed to generate tempdir name: %s", err)
		return "", nil
	}
	hash := hex.EncodeToString(r)
	tmpDir := fmt.Sprintf("%s-%s", prefix, hash)
	return tmpDir, nil
}

//GetBaseImageName takes in an image name (like projectatomic/helloapache:latest)
//and extracts the base name (helloapache)
func GetBaseImageName(image string) string {
	baseName := filepath.Base(image)
	return strings.Split(baseName, ":")[0]
}
