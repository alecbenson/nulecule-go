package cmd

import (
	"github.com/Sirupsen/logrus"
	"github.com/alecbenson/nulecule-go/atomicgo/nulecule"
	"github.com/codegangsta/cli"
)

func installFlagSet() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "no-deps",
			Usage: "Skip pulling dependencies of the app",
		},
		cli.BoolFlag{
			Name:  "update",
			Usage: "Re-pull images and overwrite existing files",
		},
		cli.StringFlag{
			Name:  "destination",
			Usage: "Destination directory for install",
		},
	}
}

//installFunction is the function that begins the install process
func installFunction(c *cli.Context) {
	//Set up parameters
	if len(c.Args()) < 1 {
		cli.ShowCommandHelp(c, "install")
		logrus.Fatalf("Please provide a repository to install from.")
	}
	app := c.Args()[0]
	dryRun := c.GlobalBool("dry-run")
	nodeps := c.Bool("no-deps")

	target := c.String("destination")
	b := nulecule.New(target, app, dryRun)
	b.Nodeps = nodeps
	b.Install()
}

//InstallCommand returns an initialized CLI stop command
func InstallCommand() cli.Command {
	return cli.Command{
		Name:   "install",
		Usage:  "Installs a nulecule project from the provided repository",
		Action: installFunction,
		Flags:  installFlagSet(),
	}
}
