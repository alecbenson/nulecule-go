package cli

import (
	"flag"
	"github.com/Sirupsen/logrus"
	"github.com/alecbenson/nulecule-go/atomicapp/constants"
	"github.com/alecbenson/nulecule-go/atomicapp/nulecule"
	"path/filepath"
)

func installFlagSet() *flag.FlagSet {
	installFlagSet := flag.NewFlagSet("install", flag.PanicOnError)
	installFlagSet.Bool("no-deps", false, "Skip pulling dependencies of the app")
	installFlagSet.Bool("update", false, "Re-pull images and overwrite existing files")
	installFlagSet.String("destination", "", "Destination directory for install")
	installFlagSet.String("APP", "", "Application to run. This is a container"+
		"image or a path that contains the metadata describing the whole application.")
	return installFlagSet
}

func installFunction() func(cmd *Command) {
	return func(cmd *Command) {
		//Set up parameters
		flags := cmd.FlagSet
		targetPath := getVal(flags, "APP").(string)
		targetFile := filepath.Join(targetPath, constants.MAIN_FILE)
		logrus.Debugf("INSTALL COMMAND: args are %v\n", targetFile)

		//Start install sequence
		base := &nulecule.Base{}
		base.New(targetPath)
		base.ReadMainFile(targetFile)
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

//The install process is as follows:
//The nulecule file is parsed (loadMainFile)
//If not a dry run...
//pull the app from docker registry (or use a specific registry if given)
//create a new instance of the container and copy the application-entity files of the container to a temp dir on the host
//remove the contianer

//verify that the user is not attempting to run install in an existing nulecule directory
//If they are, complain and throw an error

//check spec version and artifacts
//if !nodeps,
//read the mainfile data and call install() on all dependencies
