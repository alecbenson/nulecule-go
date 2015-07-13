package cli

import (
	"flag"
	"fmt"
)

func stopFlagSet() *flag.FlagSet {
	stopFlagSet := flag.NewFlagSet("stop", flag.PanicOnError)
	stopFlagSet.String("APP", "", "Path to the directory where the atomicapp"+
		"is installed or an image containing atomicapp which should be stopped.")
	return stopFlagSet
}

func stopFunction() func(cmd *Command) {
	return func(cmd *Command) {
		fmt.Printf("STOP COMMAND\n")
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
