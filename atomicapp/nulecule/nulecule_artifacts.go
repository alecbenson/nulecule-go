package nulecule

import (
	"errors"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/alecbenson/nulecule-go/atomicapp/utils"
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
