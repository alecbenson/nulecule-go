package cmd

import (
	"github.com/alecbenson/nulecule-go/atomicgo/nulecule"
	"github.com/codegangsta/cli"
)

func stopFlagSet() []cli.Flag {
	return []cli.Flag{
	//No stop flags exist at this time
	}
}

func stopFunction(c *cli.Context) {
	//Set up parameters
	var (
		app    string
		target string
	)
	dryRun := c.GlobalBool("dry-run")
	base := nulecule.New(target, app, dryRun)
	base.Stop()
}

//StopCommand returns an initialized CLI stop command
func StopCommand() cli.Command {
	return cli.Command{
		Name:   "stop",
		Usage:  "Undeploys the application and all of its providers",
		Action: stopFunction,
		Flags:  stopFlagSet(),
	}
}
