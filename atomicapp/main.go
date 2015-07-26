package main

import (
	"github.com/alecbenson/nulecule-go/atomicapp/cli"
)

func main() {
	cli.InitGeneralFlags()

	commands := []*cli.Command{
		cli.RunCommand(),
		cli.InstallCommand(),
		cli.StopCommand(),
	}

	//Iterate through all commands and dispatch appropriately
	cli.ParseCommands(commands)
}
