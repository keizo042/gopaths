package gopaths

import (
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type (
	App struct {
	}

	Config struct {
	}

	RepoInfo struct {
		Repos []string
	}
)

func getInfo(fpath string) (*RepoInfo, error) {
	var rinfo RepoInfo
	if _, err := toml.DecodeFile(fpath, &rinfo); err != nil {
		return nil, errors.Wrap(err, "toml")
	}
	return &rinfo, nil
}

func NewApp(c *Config) (*App, error) {
	return &App{}, nil
}

func (app *App) Init() error {
	return errors.New("TBD")
}

func (app *App) Config() error {
	return errors.New("TBD")
}

func (app *App) Enable() error {
	return errors.New("TBD")
}

func (app *App) Disable() error {
	return errors.New("TBD")
}

func (app *App) Add() error {
	return errors.New("TBD")
}

func (app *App) Remove() error {
	return errors.New("TBD")
}

func (app *App) Complete() error {
	return errors.New("TBD")
}
