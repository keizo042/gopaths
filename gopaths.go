package gopaths

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"os"
	"strings"
)

var (
	GOPATHS_ENV_ORIGINAL_GOPATH = "GOPATHS_GOPATH"
	ERR_NOTIMPL                 = errors.New("NotImplemented")
)

type (
	App struct {
		GOPATH    string
		ReposPath string
		Info      *RepoInfo
	}

	AppInitConfig struct {
	}

	AppConfigConfig struct {
	}

	AppAddConfig struct {
		Path string
	}

	AppRemoveConfig struct {
		Path string
	}

	AppEnableConfig struct {
	}

	AppDisableConfig struct {
	}

	Config struct {
		reposPath string
	}

	RepoInfo struct {
		Version      int32
		GOPATH       string
		Repos        []string
		DisableRepos []string
	}
)

func getInfo(fpath string) (*RepoInfo, error) {
	var rinfo RepoInfo
	if _, err := toml.DecodeFile(fpath, &rinfo); err != nil {
		return nil, errors.Wrap(err, "toml")
	}
	return &rinfo, nil
}

func setInfo(fpath string, rinfo *RepoInfo) error {
	f, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer f.Close()
	return toml.NewEncoder(f).Encode(rinfo)
}

func NewApp(c *Config) (*App, error) {
	inf, err := getInfo(c.reposPath)
	if err != nil {
		return nil, errors.Wrap(err, "repos info")
	}
	return &App{
		GOPATH: os.Getenv("GOPATH"),
		Info:   inf,
	}, nil
}

// Init is `gopaths init`.
// for initalizing GOPATH which gopaths maintain.
func (app *App) Init() error {
	path := strings.Join(app.Info.Repos, ":")
	if err := os.Setenv(GOPATHS_ENV_ORIGINAL_GOPATH, path); err != nil {
		return err
	}
	if _, err := fmt.Printf("export GOPATH=$GOPATH:%s", path); err != nil {
		return err
	}
	return nil
}

// Config is `gopaths config`
// for config gopaths configuration.
func (app *App) Config() error {
	return ERR_NOTIMPL
}

// Enable is `gopaths enable`.
//
func (app *App) Enable() error {
	paths := strings.Join(app.Info.Repos, ":")
	return os.Setenv("GOPATH", app.Info.GOPATH+":"+paths)
}

// Disable is `gopaths disable`.
func (app *App) Disable(config AppDisableConfig) error {
	return os.Setenv("GOPATH", app.Info.GOPATH)
}

// Add is `gopaths add`.
//
func (app *App) Add(config AppAddConfig) error {
	rinfo, err := getInfo(app.ReposPath)
	if err != nil {
		return err
	}
	rinfo.Repos = append(rinfo.Repos, config.Path)
	return setInfo(app.ReposPath, rinfo)
}

// Remove is `gopaths remove`
//
func (app *App) Remove(config AppRemoveConfig) error {
	var repos []string
	for _, path := range app.Info.Repos {
		if strings.Compare(path, config.Path) != 0 {
			repos = append(repos, path)
		}
	}
	return setInfo(app.ReposPath, app.Info)
}

// Complete is `gopaths complete`.
//
func (app *App) Complete() error {
	return ERR_NOTIMPL
}

// Clean is `gopaths clean`.
func (app *App) Clean() error {
	return ERR_NOTIMPL
}
