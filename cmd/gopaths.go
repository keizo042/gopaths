package main

// gopaths

import (
	"fmt"
	"os"

	"github.com/keizo042/gopaths"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

const (
	defaultConfigPath = "./.config/gopaths/gopaths.toml"
)

type (
	Cli struct {
		*cli.App
	}
)

var (
	globalFlags = []cli.Flag{}
	configFlags = []cli.Flag{}
	initFlags   = []cli.Flag{}
	commands    = []cli.Command{
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
			Name: "add",
		},
		{
			Name: "remove",
		},
		{
			Name: "complete",
		},
	}
)

func NewCli() *Cli {
	app := cli.NewApp()
	app.Name = gopaths.APP_NAME
	app.Version = gopaths.APP_VERSION_TEXT
	app.Usage = "mutiple gopath manager"

	return &Cli{
		App: app}

}

func NewConfig(ctx *cli.Context) (*gopaths.Config, error) {
	cfg := &gopaths.Config{}
	return cfg, nil
}

func ActionConfig(ctx *cli.Context) error {
	cfg, err := NewConfig(ctx)
	if err != nil {
		return errors.Wrap(err, "config")
	}
	app, err := gopaths.NewApp(cfg)
	if err != nil {
		return errors.Wrap(err, "app")
	}
	return errors.Wrap(app.Config(), "runtime")
}

func ActionEnable(ctx *cli.Context) error {
	cfg, err := NewConfig(ctx)
	if err != nil {
		return errors.Wrap(err, "config")
	}
	app, err := gopaths.NewApp(cfg)
	if err != nil {
		return errors.Wrap(err, "app")
	}
	return errors.Wrap(app.Config(), "runtime")
}

func ActionDisable(ctx *cli.Context) error {
}

func ActionUpdate(ctx *cli.Context) error {
}

func ActionAdd(ctx *cli.Context) error {
}

func ActionRemove(ctx *cli.Context) error {
}

func ActionComplete(ctx *cli.Context) error {
}

func main() {
	c := NewCli()
	if err := c.Run(os.Args); err != nil {
		fmt.Printf("%s: %s", gopaths.APP_NAME, err.Error())
	}
}
