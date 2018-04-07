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
	GOPATHS_CONFIG_FILE         = "config.toml"
	GOPATHS_GOPATHS_FILE        = "gopaths.toml"
	ERR_NOTIMPL                 = errors.New("NotImplemented")
)

type (
	Config struct {
		SettingPath string
	}

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
	if err := checkGopathsConfig(c.SettingPath); err != nil {
		return nil, errors.Wrap(err, "config dir")
	}
	inf, err := getInfo(c.SettingPath + GOPATHS_GOPATHS_FILE)
	if err != nil {
		return nil, errors.Wrap(err, "repos info")
	}
	return &App{
		GOPATH:    os.Getenv("GOPATH"),
		ReposPath: c.SettingPath,
		Info:      inf,
	}, nil
}

func checkGopathsConfig(fpath string) error {
	if _, err := os.Stat(fpath); err != nil {
		if err := os.MkdirAll(fpath, 0777); err != nil {
			return err
		}
	}
	gopaths_file := fpath + GOPATHS_GOPATHS_FILE
	if _, err := os.Stat(gopaths_file); err != nil {
		if _, err := os.Create(gopaths_file); err != nil {
			return err
		}
	}
	config_file := fpath + GOPATHS_CONFIG_FILE
	if _, err := os.Stat(config_file); err != nil {
		if _, err := os.Create(fpath + GOPATHS_CONFIG_FILE); err != nil {
			return err
		}
	}

	return ERR_NOTIMPL
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
// for manage gopaths config.
func (app *App) Config(config *AppConfigConfig) error {
	return ERR_NOTIMPL
}

// Enable is `gopaths enable`.
// for enable gopaths we set.
func (app *App) Enable() error {
	paths := strings.Join(app.Info.Repos, ":")
	return os.Setenv("GOPATH", app.Info.GOPATH+":"+paths)
}

// Disable is `gopaths disable`.
// for disable gopaths we set and recover original gopath.
func (app *App) Disable(config *AppDisableConfig) error {
	return os.Setenv("GOPATH", app.Info.GOPATH)
}

// Add is `gopaths add`.
// for adding path to gopaths manage.
func (app *App) Add(config AppAddConfig) error {
	rinfo, err := getInfo(app.ReposPath)
	if err != nil {
		return err
	}
	rinfo.Repos = append(rinfo.Repos, config.Path)
	return setInfo(app.ReposPath, rinfo)
}

// Remove is `gopaths remove`
// for removing path to gopaths manage.
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
// for generating bash completion
func (app *App) Complete() error {
	return ERR_NOTIMPL
}

// Clean is `gopaths clean`.
// clean up parameter `gopaths` settings.
// mainly use, enviroment variables.
func (app *App) Clean() error {
	return ERR_NOTIMPL
}
