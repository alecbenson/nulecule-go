package nulecule

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/alecbenson/nulecule-go/atomicgo/constants"
	"github.com/alecbenson/nulecule-go/atomicgo/utils"
)

//ArtifactEntry is a source control repository struct used to specify an artifact
type ArtifactEntry struct {
	Path string
	Repo SrcControlRepo
}

type SrcControlRepo struct {
	Inherit []string
	Source  string
	Path    string
	Type    string
	Branch  string
	Tag     string
}

//IsExternal returns true if the component is an external resource, and false if the component is
//A local resource
func IsExternal(component Component) bool {
	if len(component.Artifacts) == 0 {
		return true
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

//GetExternalAppDirectory returns the directory in which ane external app is installed to
func (b *Base) makeExternalAppDirectory(c Component) (string, error) {
	fp := b.GetExternallAppDirectory(c)
	err := os.MkdirAll(fp, 0700)
	if err != nil {
		logrus.Fatalf("Failed to make external app directory in %s", b.Target())
		return "", err
	}
	return fp, nil
}

//GetExternallAppDirectory returns the directory in which an external app is installed to
func (b *Base) GetExternallAppDirectory(c Component) string {
	return filepath.Join(b.Target(), constants.EXTERNAL_APP_DIR, c.Name)
}

//Generates a new work directory
//Returns a path to the directory
func makeWorkDirectory(targetPath string) (string, error) {
	workdir := filepath.Join(targetPath, constants.WORKDIR)
	//If the .workdir directory does not exist in targetPath, make it.
	if !utils.PathExists(workdir) {
		logrus.Debugf("Making workdir in %s", targetPath)
		err := os.MkdirAll(workdir, 0700)
		if err != nil {
			logrus.Fatalf("Failed to make work directory in %s", targetPath)
			return "", errors.New("Failed to make work directory")
		}
	}
	return workdir, nil
}

//SaveArtifact writes a templated artifact to the .workdir directory.
//If .workdir does not exist, it is created.
//data - a []byte of the templated file
//name - the name of the file to write to
func SaveArtifact(data []byte, targetPath, name string) error {
	workdir, err := makeWorkDirectory(targetPath)
	if err != nil {
		return err
	}

	//Create the file to write the template to
	fullPath := filepath.Join(workdir, name)
	templateFile, err := os.Create(fullPath)
	if err != nil {
		logrus.Fatalf("Unable to create template file: %s", err)
		return errors.New("Failed to create template file")
	}

	if utils.WriteToFile(data, templateFile); err != nil {
		return err
	}
	return nil
}
