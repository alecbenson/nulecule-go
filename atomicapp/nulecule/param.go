package nulecule

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/Sirupsen/logrus"
)

//Param represents the Component parameters
type Param struct {
	Name        string
	Description string
	Constraints []Constraint
	Default     string
	Hidden      bool
}

//Constraint is a struct representing a constaint for a parameter object
type Constraint struct {
	AllowedPattern string
	Description    string
}

//CheckConstraints returns true if the parameter has matching, valid constraints. False otherwise.
func CheckConstraints(param Param) (bool, error) {
	value := param.Default
	for _, constraint := range param.Constraints {
		pattern := constraint.AllowedPattern
		valid, err := regexp.MatchString(pattern, value)
		if err != nil {
			logrus.Errorf("Error checking constraint for parameter %s: %s", value, err)
			return false, err
		}

		if !valid {
			return false, nil
		}
	}
	return true, nil
}

//ApplyTemplate reads the file located at artifactPath and replaces all parameters
//with those specified in the nulecule parameters objects
//artifactPath is the path provided by the artifact. The data from this file
//is loaded, and all variables ($var_name) get replaced with their correct values
//TargetPath is the base directory to install the workdir directory into
func ApplyTemplate(artifactPath string, targetPath string, params []Param) ([]byte, error) {
	data, err := ioutil.ReadFile(artifactPath)
	if err != nil {
		return data, err
	}
	//Replaces every instance of $param with
	//the value provided in the nulecule file
	for _, param := range params {
		name := bytes.Trim([]byte(param.Name), "$")
		value := bytes.Trim([]byte(param.Default), "$")
		key := []byte(param.Name)

		//We don't want to replace to/from empty values
		if name == nil || value == nil {
			continue
		}

		//Ensure that each parameter has valid values
		valid, err := CheckConstraints(param)
		if err != nil {
			logrus.Errorf("Error checking constraint %s: %s", name, err)
		}
		if !valid {
			logrus.Fatalf("Parameter '%s' has violated a regex constraint", param.Name)
		}
		data = bytes.Replace(data, key, value, -1)
	}
	name := filepath.Base(artifactPath)
	SaveArtifact(data, targetPath, name)
	return data, nil
}
