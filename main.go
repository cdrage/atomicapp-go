package main

import (
	"github.com/cdrage/atomicapp-go/cmd"
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
