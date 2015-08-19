package container

import (
	"fmt"
	"github.com/alecbenson/nulecule-go/atomicgo/constants"
	"github.com/alecbenson/nulecule-go/atomicgo/utils"
	"os/exec"
)

const (
	CONTAINTER_PLATFORM = "docker"
)

//Remove removes the container with the specified name
func Remove(containerName string) error {
	dockerRm := exec.Command(CONTAINTER_PLATFORM, "rm", containerName)
	if _, err := utils.CheckCommandOutput(dockerRm, true); err != nil {
		return err
	}
	return nil
}

//Create creates a container with the given name from the image provided
func Create(containerName, image string) error {
	dockerCreate := exec.Command(CONTAINTER_PLATFORM, "create", "--name", containerName, image, "nop")
	if _, err := utils.CheckCommandOutput(dockerCreate, false); err != nil {
		return err
	}
	return nil
}

//Copy copies the copyFrom directory from container with name containerName
//to the copyTo directory
func Copy(containerName, copyFrom, copyTo string) error {
	cpField := fmt.Sprintf("%s:%s", containerName, constants.APP_ENT_PATH)
	dockerCp := exec.Command(CONTAINTER_PLATFORM, "cp", cpField, copyTo)
	if _, err := utils.CheckCommandOutput(dockerCp, false); err != nil {
		return err
	}
	return nil
}

//Pull pulls the container with the given name from docker
func Pull(image string) error {
	dockerPull := exec.Command(CONTAINTER_PLATFORM, "pull", image)
	if _, err := utils.CheckCommandOutput(dockerPull, false); err != nil {
		return err
	}
	return nil
}
