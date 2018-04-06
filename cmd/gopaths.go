package main

// gopaths

import (
	"fmt"
	"os"

	"github.com/keizo042/gopaths"
	"github.com/urfave/cli"
)

type (
	Cli struct {
		*cli.App
	}
)

var (
	flags    = []cli.Flag{}
	commands = []cli.Command{
		{
			Name:   "config",
			Action: ActionConfig,
		},
		{
			Name:   "enable",
			Action: ActionEnable,
		},
		{
			Name:   "disable",
			Action: ActionDisable,
		},
		{
			Name: "update",
		},
		{
			Name: "init",
		},
		{
			Name: "delete",
		},
		{
			Name: "bash-completation",
		},
	}
)

func NewCli() *Cli {
	app := cli.NewApp()

	return &Cli
App:
	app{}

}

func NewConfig(ctx *cli.Context) (*gopaths.Config, error) {
}
func main() {
	c := NewCli()
	if err := c.Run(os.Args); err != nil {
		fmt.Printf("%s: %s", gopaths.APP_NAME, err.Error())
	}
}

func ActionConfig(ctx *cli.Context) error {
}

func ActionEnable(ctx *cli.Context) error {
}

func ActionDisable(ctx *cli.Context) error {
}

func ActionUpdate(ctx *cli.Context) error {
}

func ActionInit(ctx *cli.Context) error {
}

func ActionDelete(ctx *cli.Context) error {
}

func ActionBashCompletation(ctx *cli.Context) error {
}
