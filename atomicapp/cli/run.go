package cli

import (
	"flag"
	"github.com/alecbenson/nulecule-go/atomicapp/nulecule"
)

//RunCommand returns an initialized CLI stop command
func RunCommand() *Command {
	return &Command{
		Name:    "run",
		Func:    runFunction(),
		FlagSet: runFlagSet(),
	}
}

func runFlagSet() *flag.FlagSet {
	runFlagSet := flag.NewFlagSet("run", flag.PanicOnError)
	runFlagSet.Bool("ask", false, "Ask for params even if the defaul value is provided")
	runFlagSet.String("write", "", "A file which will contain anwsers provided in interactive mode")
	return runFlagSet
}

func runFunction() func(cmd *Command, args []string) {
	return func(cmd *Command, args []string) {
		//Set up parameters
		var (
			app    string
			target string
		)
		if len(args) > 0 {
			target = args[0]
		}
		flags := cmd.FlagSet
		ask := getVal(flags, "ask").(bool)
		answersFile := getVal(flags, "write").(string)
		//Start run sequence
		base := nulecule.New(target, app)
		base.Run(ask, answersFile)
	}
}
