package cli

import (
	"flag"

	"github.com/Sirupsen/logrus"
	"github.com/alecbenson/nulecule-go/atomicapp/nulecule"
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
		b.Install()
		b.WriteAnswersSample()
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
