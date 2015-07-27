package cli

import (
	"flag"
	"fmt"
)

func stopFlagSet() *flag.FlagSet {
	stopFlagSet := flag.NewFlagSet("stop", flag.PanicOnError)
	return stopFlagSet
}

func stopFunction() func(cmd *Command, args []string) {
	return func(cmd *Command, args []string) {
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
