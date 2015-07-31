package nulecule

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/Sirupsen/logrus"
	"github.com/alecbenson/nulecule-go/atomicapp/constants"
	"github.com/alecbenson/nulecule-go/atomicapp/container"
	"github.com/alecbenson/nulecule-go/atomicapp/utils"
)

//Install starts the install process
func (b *Base) Install() error {
	app := b.App()
	target := b.Target()

	if err := container.Pull(app); err != nil {
		return err
	}
	tmpFiles, err := copyFromContainer(app)
	if err != nil {
		return err
	}
	if err := checkInstallOverwrite(tmpFiles, target); err != nil {
		return err
	}
	if err := utils.CopyDirectory(tmpFiles, target); err != nil {
		logrus.Errorf("Error copying directory: %s", err)
		return err
	}
	b.ReadMainFile()
	b.CheckSpecVersion()
	b.CheckAllArtifacts()
	b.InstallDependencies()
	return nil
}

//InstallDependencies install all external sources
func (b *Base) InstallDependencies() error {
	components := b.MainfileData.Graph
	for _, c := range components {
		b.UpdateAnswers(c)
		if !IsExternal(c) {
			continue
		}
		logrus.Infof("Installing dependency: %s", c.Name)
		image, err := GetSourceImage(c)
		if err != nil {
			return err
		}
		externalDir, err := b.makeExternalAppDirectory(c)
		if err != nil {
			return err
		}
		externalBase := New(externalDir, image)
		if err := externalBase.Install(); err != nil {
			logrus.Panicf("Failed to install dependency: %s", err)
			return err
		}
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
	if err := container.Create(containerName, image); err != nil {
		return "", err
	}
	//Generate a unique name for the temporary directory
	tmpDirName, err := utils.GenerateUniqueName("nulecule")
	if err != nil {
		return "", err
	}
	tmpDirPath := filepath.Join(os.TempDir(), tmpDirName)
	//Copy the contents to the temporary application-entity directory (usually /tmp/nulecule-xxxxxx)
	if err := container.Copy(containerName, constants.APP_ENT_PATH, tmpDirPath); err != nil {
		return "", err
	}
	logrus.Debugf("Copied container application entity data to %s", tmpDirPath)
	if err := container.Remove(containerName); err != nil {
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
	targetBase := New(dstTarget, "")
	targetBase.ReadMainFile()
	targetID := targetBase.MainfileData.ID

	//Parse the nulecule file in the tmp directory
	tmpBase := New(tmpTarget, "")
	tmpBase.ReadMainFile()
	tmpID := tmpBase.MainfileData.ID

	if targetID != tmpID {
		logrus.Fatalf("You are trying to overwrite existing app %s with app %s - clear or change current directory.", targetID, tmpID)
		return errors.New("Conflicting app install path")
	}
	return nil
}
