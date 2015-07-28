package nulecule

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"regexp"

	"github.com/Sirupsen/logrus"
	"github.com/alecbenson/nulecule-go/atomicapp/constants"
)

//Param represents the Component parameters
type Param struct {
	Name        string
	Description string
	Constraints []Constraint
	Default     string
	Hidden      bool
	AskedFor    bool
}

//Constraint is a struct representing a constaint for a parameter object
type Constraint struct {
	AllowedPattern string
	Description    string
}

//CheckConstraints returns true if the parameter has matching, valid constraints. False otherwise.
func checkConstraints(param *Param) (bool, error) {
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
func ApplyTemplate(artifactPath string, targetPath string, c *Component, ask bool) ([]byte, error) {
	data, err := ioutil.ReadFile(artifactPath)
	if err != nil {
		return data, err
	}
	//Replaces every instance of $param with
	//the value provided in the nulecule file
	for index := range c.Params {
		param := &c.Params[index]
		//We run strip on $ just in case one exists.
		name := bytes.Trim([]byte(param.Name), "$")
		key := []byte(fmt.Sprintf("$%s", name))
		value := bytes.Trim([]byte(param.Default), "$")
		if name == nil || value == nil {
			continue
		}

		//If the user wants to be asked for all parameters, do so
		if ask && !param.AskedFor {
			response, err := askForParam(param)
			if err != nil {
				return data, err
			}
			data = bytes.Replace(data, key, response, -1)
			continue
		}

		//Ensure that each parameter has valid values
		valid, err := checkConstraints(param)
		if err != nil {
			logrus.Errorf("Error checking constraint %s: %s", param.Name, err)
		}
		if !valid {
			logrus.Fatalf("Parameter '%s' has violated a regex constraint", param.Name)
		}
		data = bytes.Replace(data, key, value, -1)
	}

	//Find any last unresolved parameters and replace them with user input values
	if data, err = askMissingParams(data); err != nil {
		logrus.Fatalf("Failed to ask for missing parameters: %s", err)
		return data, err
	}
	name := filepath.Base(artifactPath)
	SaveArtifact(data, targetPath, name)
	return data, nil
}

//Ask mising params will search each artifact for
//unresolved parameters and ask the user to provide a value
//It will return a []byte containing artifact data with replacements inserted
func askMissingParams(data []byte) ([]byte, error) {
	paramRegex := regexp.MustCompile("(\\$[a-zA-Z0-9]*)\\w")
	unknownParams := paramRegex.FindAll(data, -1)

	//Iterate through matches...
	for _, param := range unknownParams {
		name := bytes.Trim([]byte(param), "$")
		formattedParam := Param{Name: string(name)}

		//Query user for a valid parameter value...
		value, err := askForParam(&formattedParam)
		if err != nil {
			logrus.Errorf("Failed to ask for parameter %s", name)
		}
		data = bytes.Replace(data, param, value, -1)
	}
	return data, nil
}

//Queries the user for the passed in parameter.
//If the parameter is invalid (it violates a constraint, for example),
//It will prompt the user again. Returns the user's response
func askForParam(param *Param) ([]byte, error) {
	var (
		validated bool
		value     []byte
	)
	//Set the query description
	description := param.Description
	if description == "" {
		description = fmt.Sprintf("No description available. Add this parameter to your %s file", constants.MAIN_FILE)
	}

	//Query until the user gives us valid input
	for !validated {
		fmt.Printf("Enter a value for %s (%s - default: %s): ", param.Name, description, param.Default)
		if _, err := fmt.Scanln(&value); err != nil {
			continue
		}
		//Verify that the entry is valid
		valid, err := checkConstraints(param)
		if err != nil {
			return value, err
		}
		if !valid {
			logrus.Warnf("Invalid input provided for %s. Try again.", param.Name)
		}
		validated = valid && (value != nil)
	}
	//Don't ask the user to provide a value for this param again
	param.AskedFor = true
	return value, nil
}

/*
TODO: When a parameter is asked for, we need to make sure that it doesn't get asked for repeatedly as each artifact gets templated
*/
