package cli

import (
	"errors"
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
		argCount := len(args)
		if argCount <= 0 {
			logrus.Fatalf("Please provide a repository to install from.")
		}
		app := args[0]
		flags := cmd.FlagSet
		target := getVal(flags, "destination").(string)
		b := nulecule.New(target, app)
		//Start install sequence
		if err := pullContainer(app); err != nil {
			return
		}
		//Copy files from the container to a temporary directory
		tmpFiles, err := copyFromContainer(app)
		if err != nil {
			return
		}
		//Verify that the install target would not overwrite an existing nulecule app
		if err := checkInstallOverwrite(tmpFiles, b.Target()); err != nil {
			return
		}
		//Move container files to the install directory
		err = utils.CopyDirectory(tmpFiles, b.Target())
		if err != nil {
			logrus.Errorf("Error copying directory: %s", err)
			return
		}
		b.ReadMainFile()
		b.CheckSpecVersion()
		b.CheckAllArtifacts()
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

//checkInstallOverwrite returns a nil error if no nulecule app exists in the target directory.
//If a nulecule app does exist in the target directory, it checks if it has a matching app ID and
//if it does, allows the installation process to continue. If a different nulecule app has been found,
//it warns the user that they are overwriting an existing app and returns an error status
func checkInstallOverwrite(tmpTarget, dstTarget string) error {
	if !utils.PathExists(filepath.Join(dstTarget, constants.MAIN_FILE)) {
		logrus.Debugf("No existing app found in target directory, installation can continue")
		return nil
	}
	//Parse the nulecule file in the target directory
	targetBase := nulecule.New(dstTarget, "")
	targetBase.ReadMainFile()
	targetID := targetBase.MainfileData.ID

	//Parse the nulecule file in the tmp directory
	tmpBase := nulecule.New(tmpTarget, "")
	tmpBase.ReadMainFile()
	tmpID := tmpBase.MainfileData.ID

	if targetID != tmpID {
		logrus.Fatalf("You are trying to overwrite existing app %s with app %s - clear or change current directory.", targetID, tmpID)
		return errors.New("Conflicting app install path")
	}
	return nil
}
