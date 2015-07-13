package cli

import (
	"flag"
	"fmt"
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
		fmt.Printf("INSTALL COMMAND\n")
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
