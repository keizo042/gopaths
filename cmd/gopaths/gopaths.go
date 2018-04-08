package main

// gopaths

import (
	"fmt"
	"os"

	"github.com/keizo042/gopaths"
	"github.com/pkg/errors"
	"github.com/urfave/cli"
)

var (
	defaultGoPathsSettings = ".config/gopaths/"
)

type (
	Cli struct {
		*cli.App
	}
)

var (
	globalFlags = []cli.Flag{}
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
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "show",
				},
				cli.BoolFlag{
					Name: "verbose",
				},
				cli.StringFlag{
					Name: "set",
				},
			},
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
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "verbose",
				},
			},
		},
		{
			Name:   "rm",
			Action: ActionRemove,
			Usage:  "remove repo that be maintained by gopaths",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name: "force",
				},
				cli.BoolFlag{
					Name: "all",
				},
				cli.BoolFlag{
					Name: "verbose",
				},
			},
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
	cfg := &gopaths.Config{
		SettingPath: ctx.String("configpath"),
	}
	if cfg.SettingPath == "" {
		homedir := os.Getenv("HOME")
		cfg.SettingPath = homedir + "/" + defaultGoPathsSettings
	}
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
	return errors.Wrap(app.Init(), "init")
}
func ActionConfig(ctx *cli.Context) error {
	app, err := newApp(ctx)
	if err != nil {
		return errors.Wrap(err, "initalize app")
	}
	return errors.Wrap(app.Config(&gopaths.AppConfigConfig{
		Args:    ctx.Args(),
		Show:    ctx.Bool("show"),
		Verbose: ctx.Bool("verbose"),
	}), "config")
}

func ActionEnable(ctx *cli.Context) error {
	app, err := newApp(ctx)
	if err != nil {
		return err
	}
	return errors.Wrap(app.Enable(), "disable")
}

func ActionDisable(ctx *cli.Context) error {
	app, err := newApp(ctx)
	if err != nil {
		return err
	}
	return errors.Wrap(app.Disable(&gopaths.AppDisableConfig{}), "disable")
}

func ActionAdd(ctx *cli.Context) error {
	app, err := newApp(ctx)
	if err != nil {
		return err
	}
	return errors.Wrap(app.Add(&gopaths.AppAddConfig{
		Paths: ctx.Args(),
	}), "add")
}

func ActionRemove(ctx *cli.Context) error {
	app, err := newApp(ctx)
	if err != nil {
		return err
	}
	return errors.Wrap(app.Remove(&gopaths.AppRemoveConfig{
		Paths:   ctx.Args(),
		Verbose: ctx.Bool("verbose"),
		All:     ctx.Bool("all"),
	}), "remove")
}

func ActionResotre(ctx *cli.Context) error {
	app, err := newApp(ctx)
	if err != nil {
		return errors.Wrap(err, "initalize app")
	}
	return errors.Wrap(app.Restore(), "restore")
}

func ActionComplete(ctx *cli.Context) error {
	app, err := newApp(ctx)
	if err != nil {
		return err
	}
	return errors.Wrap(app.Complete(), "complete")
}

func main() {
	c := NewCli()
	if err := c.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", gopaths.APP_NAME, err.Error())
	}
}
