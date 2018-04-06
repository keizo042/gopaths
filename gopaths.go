package gopaths

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"os"
	"strings"
)

type (
	App struct {
		GOPath string
		Info   *RepoInfo
	}

	Config struct {
		configPath string
		reposPath  string
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
	inf, err := getInfo(c.reposPath)
	if err != nil {
		return nil, errors.Wrap(err, "repos info")
	}
	return &App{
		GOPath: os.Getenv("GOPATH"),
		Info:   inf,
	}, nil
}

func (app *App) Init() error {
	path := strings.Join(app.Info.Repos, ":")
	return fmt.Printf("export GOPATH=$GOPATH:%s", path)
}

func (app *App) Config() error {
	return nil
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
