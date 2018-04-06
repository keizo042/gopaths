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
	defaultConfigPath = "./.config/gopaths/config.toml"
	defaultReposPath  = "./.config/gopaths/repos.toml"
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
			Name:   "init",
			Action: ActionInit,
			Usage:  "manage bash config",
		},
		{
			Name:   "config",
			Action: ActionConfig,
			Usage:  "configuration",
		},
		{
			Name:   "enable",
			Action: ActionEnable,
			Usage:  "enable gopaths's GOPATH",
		},
		{
			Name:   "disable",
			Action: ActionDisable,
			Usage:  "dsiable gopaths's GOPATH",
		},
		{
			Name:   "add",
			Action: ActionAdd,
			Usage:  "add repo that be maintained by gopaths",
		},
		{
			Name:   "remove",
			Action: ActionRemove,
			Usage:  "remove repo that be maintained by gopaths",
		},
		{
			Name:   "complete",
			Action: ActionComplete,
			Usage:  "bash completion command",
		},
	}
)

func NewCli() *Cli {
	app := cli.NewApp()
	app.Name = gopaths.APP_NAME
	app.Version = gopaths.APP_VERSION_TEXT
	app.Commands = commands
	app.Flags = globalFlags
	app.Usage = "mutiple gopath manager"
	return &Cli{
		App: app}

}

func NewConfig(ctx *cli.Context) (*gopaths.Config, error) {
	cfg := &gopaths.Config{}
	return cfg, nil
}

func newApp(ctx *cli.Context) (*gopaths.App, error) {
	cfg, err := NewConfig(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "config")
	}
	app, err := gopaths.NewApp(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "app")
	}
	return app, nil
}

func ActionInit(ctx *cli.Context) error {
	app, err := newApp(ctx)
	if err != nil {
		return err
	}
	return app.Init()
}
func ActionConfig(ctx *cli.Context) error {
	app, err := newApp(ctx)
	if err != nil {
		return err
	}
	return app.Config()
}

func ActionEnable(ctx *cli.Context) error {
	app, err := newApp(ctx)
	if err != nil {
		return err
	}
	return app.Enable()
}

func ActionDisable(ctx *cli.Context) error {
	app, err := newApp(ctx)
	if err != nil {
		return err
	}
	return app.Disable()
}

func ActionAdd(ctx *cli.Context) error {
	app, err := newApp(ctx)
	if err != nil {
		return err
	}
	return app.Add()
}

func ActionRemove(ctx *cli.Context) error {
	app, err := newApp(ctx)
	if err != nil {
		return err
	}
	return app.Remove()
}

func ActionComplete(ctx *cli.Context) error {
	app, err := newApp(ctx)
	if err != nil {
		return err
	}
	return app.Complete()
}

func main() {
	c := NewCli()
	if err := c.Run(os.Args); err != nil {
		fmt.Printf("%s: %s", gopaths.APP_NAME, err.Error())
	}
}
