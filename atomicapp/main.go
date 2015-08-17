package main

import (
	"github.com/alecbenson/nulecule-go/atomicapp/cmd"
	"github.com/codegangsta/cli"
)

func main() {
	commands := []cli.Command{
		cmd.RunCommand(),
		cmd.InstallCommand(),
		cmd.StopCommand(),
	}

	//Initialize the application and run
	cmd.InitApp(commands)
}
