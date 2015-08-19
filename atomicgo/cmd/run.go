package cmd

import (
	"github.com/alecbenson/nulecule-go/atomicgo/nulecule"
	"github.com/codegangsta/cli"
)

//RunCommand returns an initialized CLI stop command
func RunCommand() cli.Command {
	return cli.Command{
		Name:   "run",
		Usage:  "Deploys a nulecule application and all of its providers",
		Action: runFunction,
		Flags:  runFlagSet(),
	}
}

func runFlagSet() []cli.Flag {
	return []cli.Flag{
		cli.BoolFlag{
			Name:  "ask",
			Usage: "If true, the user will be prompted to provide the values of all parameters",
		},
		cli.StringFlag{
			Name:  "write",
			Usage: "A file which will contain anwsers provided in interactive mode",
		},
	}
}

func runFunction(c *cli.Context) {
	//Set up parameters
	var (
		app    string
		target string
	)
	ask := c.Bool("ask")
	answersFile := c.String("write")
	dryRun := c.GlobalBool("dry-run")

	if len(c.Args()) >= 1 {
		target = c.Args()[0]
	}
	//Start run sequence
	base := nulecule.New(target, app, dryRun)
	base.Run(ask, answersFile)
}
