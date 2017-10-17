package main

import (
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

var (
	homeDirPath, _ = homedir.Dir()
)

const (
	ExitCodeError = 1
)

func main() {
	app := makeApp()
	app.Run(os.Args)
}

func makeApp() *cli.App {
	app := cli.NewApp()

	app.Name = "mino"
	app.Usage = "Make macro easily"
	app.Version = "0.0.1"

	app.Action = ExecMacro

	app.Commands = []cli.Command{
		{
			Name:    "create",
			Aliases: []string{"c"},
			Usage:   "create macro",
			Action:  CreateMacro,
		},
	}

	return app
}
