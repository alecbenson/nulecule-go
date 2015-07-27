package cli

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/alecbenson/nulecule-go/atomicapp/constants"
	"github.com/alecbenson/nulecule-go/atomicapp/nulecule"
	"github.com/alecbenson/nulecule-go/atomicapp/utils"
)

const (
	CONTAINTER_PLATFORM = "docker"
)

func installFlagSet() *flag.FlagSet {
	installFlagSet := flag.NewFlagSet("install", flag.PanicOnError)
	installFlagSet.Bool("no-deps", false, "Skip pulling dependencies of the app")
	installFlagSet.Bool("update", false, "Re-pull images and overwrite existing files")
	installFlagSet.String("destination", "", "Destination directory for install")
	return installFlagSet
}

//installFunction is the function that begins the install process
func installFunction() func(cmd *Command, args []string) {
	return func(cmd *Command, args []string) {
		//Set up parameters
		if len(args) == 0 {
			logrus.Fatalf("Please provide a repository to install from.")
		}
		target := args[0]

		//Start install sequence
		base := nulecule.New(target)
		if err := pullContainer(target); err != nil {
			return
		}
		//Copy files from the container to a temporary directory
		tmpFiles, err := copyFromContainer(target)
		if err != nil {
			return
		}
		//Move container files to the install directory
		wd, _ := os.Getwd()
		err = utils.CopyDirectory(tmpFiles, wd)
		if err != nil {
			logrus.Errorf("Error copying directory: %s", err)
			return
		}
		base.CheckSpecVersion()
	}
}

//InstallCommand returns an initialized CLI stop command
func InstallCommand() *Command {
	return &Command{
		Name:    "install",
		Func:    installFunction(),
		FlagSet: installFlagSet(),
	}
}

//pullContainer pulls the container with the given name from docker
func pullContainer(image string) error {
	dockerPull := exec.Command(CONTAINTER_PLATFORM, "pull", image)
	if _, err := utils.CheckCommandOutput(dockerPull, false); err != nil {
		return err
	}
	return nil
}

//copyFromContainer creates the container with the given image name and copies the applicationEntity directory
//to a temporary directory. Once this is done, the container is then cleaned up
//returns the path to the copied contents
func copyFromContainer(image string) (string, error) {
	//Create a new container from this image and give it a unique name
	baseName := utils.GetBaseImageName(image)
	containerName, err := utils.GenerateUniqueName(baseName)
	if err != nil {
		return "", err
	}
	dockerCreate := exec.Command(CONTAINTER_PLATFORM, "create", "--name", containerName, image, "nop")
	if _, err := utils.CheckCommandOutput(dockerCreate, false); err != nil {
		return "", err
	}

	//Generate a unique name for the temporary directory
	tmpDirName, err := utils.GenerateUniqueName("nulecule")
	if err != nil {
		return "", err
	}
	tmpDirPath := filepath.Join(os.TempDir(), tmpDirName)

	//Copy the contents to the temporary application-entity directory (usually /tmp/nulecule-xxxxxx)
	cpField := fmt.Sprintf("%s:%s", containerName, constants.APP_ENT_PATH)
	dockerCp := exec.Command(CONTAINTER_PLATFORM, "cp", cpField, tmpDirPath)
	if _, err := utils.CheckCommandOutput(dockerCp, false); err != nil {
		return "", err
	}
	logrus.Debugf("Copied container application entity data to %s", tmpDirPath)

	//Clean up the container, we don't need it anymore since we have its contents copied
	dockerRm := exec.Command(CONTAINTER_PLATFORM, "rm", containerName)
	if _, err := utils.CheckCommandOutput(dockerRm, true); err != nil {
		return "", err
	}
	return tmpDirPath, nil
}
