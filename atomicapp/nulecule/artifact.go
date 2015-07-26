package nulecule

import (
	"errors"
	"strings"
	"bytes"
	"io/ioutil"
	"path/filepath"
	"fmt"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/alecbenson/nulecule-go/atomicapp/utils"
	"github.com/alecbenson/nulecule-go/atomicapp/constants"
)

//IsExternal returns true if the component is an external resource, and false if the component is
//A local resource
func IsExternal(component Component) bool {
	if len(component.Artifacts) == 0 {
		return false
	}
	if component.Source == "" {
		return false
	}
	return true
}

//GetSourceImage fetches the sanitized source path of the image
func GetSourceImage(component Component) (string, error) {
	source := component.Source

	if !IsExternal(component) {
		logrus.Errorf("Cannot get external source of local component\n")
		return "", errors.New("Cannot get source of local component")
	}

	if strings.HasPrefix(source, utils.DOCKER_PREFIX) {
		return strings.TrimPrefix(source, utils.DOCKER_PREFIX), nil
	}

	logrus.Errorf("Could not get source image from component source: %v\n", component)
	return "", errors.New("Could not get source image")
}

//ApplyTemplate reads the file located at artifactPath and replaces all parameters
//With those specified in the nulecule parameters objects
//artifactPath is the path provided by the artifact
//TargetPath is the base directory to install the workdir directory into
func ApplyTemplate(artifactPath string, targetPath string, params []Param) ([]byte, error) {
	//Read the file
	data, err := ioutil.ReadFile(artifactPath)
	if err != nil {
		return data, err
	}
	//Replaces every instance of $param with
	//the value provided in the nulecule file
	for _, param := range params {
		if param.Name == "" || param.Default == "" {
			continue
		}
		name := bytes.Trim([]byte(param.Name), "$")
		logrus.Debugf("Param name is %s", name)
		value := bytes.Trim([]byte(param.Default), "$")
		key := []byte(fmt.Sprintf("$%s", name))
		data = bytes.Replace(data, key, value, -1)
	}
	name := filepath.Base(artifactPath)
	SaveArtifact(data, targetPath, name)
	return data, nil
}

//Saves a templated artifact to the .workdir directory.
//If .workdir does not exist, it is created.
//data - a []byte of the templated file
//name - the name of the file to write to
func SaveArtifact(data []byte, targetPath, name string) error {
	workdir := filepath.Join(targetPath, constants.WORKDIR)
	//If the .workdir directory does not exist in targetPath, make it.
	if !utils.PathExists(workdir){
		logrus.Infof("Making workdir in %s", targetPath)
		err := os.MkdirAll(workdir, 0700)
		if err != nil {
			logrus.Fatalf("Failed to make work directory in %s", targetPath)
			return errors.New("Failed to make work directory")
		}
	}

	//Create the file to write the template to
	fullPath := filepath.Join(workdir, name)
	templateFile, err := os.Create(fullPath)
	if err != nil {
		logrus.Fatalf("Unable to create template file: %s", err)
		return errors.New("Failed to create template file")
	}

	//Write the data to the template file
	defer templateFile.Close()
	out, err := templateFile.Write(data)
	if err != nil {
		logrus.Errorf("Failed to write to template file at %s", fullPath)
		return errors.New("Failed to write to template file")
	}
	logrus.Debugf("Wrote %d bytes to the template file:\n %s", out, data)
	return nil
}
