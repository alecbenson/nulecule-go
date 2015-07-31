package cli

import (
	"flag"
	"github.com/alecbenson/nulecule-go/atomicapp/nulecule"
)

func stopFlagSet() *flag.FlagSet {
	stopFlagSet := flag.NewFlagSet("stop", flag.PanicOnError)
	return stopFlagSet
}

func stopFunction() func(cmd *Command, args []string) {
	return func(cmd *Command, args []string) {
		//Set up parameters
		var (
			app    string
			target string
		)
		if len(args) > 0 {
			target = args[0]
		}
		base := nulecule.New(target, app)
		base.Stop()
	}
}

//StopCommand returns an initialized CLI stop command
func StopCommand() *Command {
	return &Command{
		Name:    "stop",
		Func:    stopFunction(),
		FlagSet: stopFlagSet(),
	}
}
