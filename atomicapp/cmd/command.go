package cmd

import (
	"os"

	"github.com/alecbenson/nulecule-go/atomicapp/constants"
	"github.com/alecbenson/nulecule-go/atomicapp/logging"
	"github.com/codegangsta/cli"
)

//InitApp initializes the application
func InitApp(commands []cli.Command) *cli.App {
	var logLevel int
	app := cli.NewApp()
	app.Name = "atomicgo"
	app.Usage = "A nulecule implementation written in Go"
	app.Version = constants.ATOMICAPPVERSION
	app.Commands = commands
	app.Flags = globalFlagSet()

	//Handle global options such as logging level
	app.Before = func(c *cli.Context) error {
		logLevel = c.GlobalInt("log-level")
		logging.SetLogLevel(logLevel)
		return nil
	}

	app.Run(os.Args)
	return app
}

func globalFlagSet() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name: "dry-run",
			Usage: "Don't actually call provider. The commands that" +
				"should be run will be sent to stdout but not run.",
		},
		cli.IntFlag{
			Name:  "log-level",
			Usage: "Set how much information to display when running commands. 0=least, 5=most.",
			Value: 5,
		},
	}
}
