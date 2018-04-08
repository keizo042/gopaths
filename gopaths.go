package gopaths

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
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
		Verbose bool
		File    string
	}

	AppConfigConfig struct {
		GOPATH  string
		Verbose bool
	}

	AppAddConfig struct {
		Paths   []string
		Verbose bool
	}

	AppRemoveConfig struct {
		Paths   []string
		Verbose bool
	}

	AppEnableConfig struct {
		Verbose bool
	}

	AppDisableConfig struct {
		Verbose bool
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
	f, err := os.OpenFile(fpath, os.O_WRONLY, 0777)
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
	fileGOPATH := fpath + GOPATHS_GOPATHS_FILE
	if _, err := os.Stat(fileGOPATH); err != nil {
		if _, err := os.Create(fileGOPATH); err != nil {
			return err
		}
	}
	fileConfig := fpath + GOPATHS_CONFIG_FILE
	if _, err := os.Stat(fileConfig); err != nil {
		if _, err := os.Create(fileConfig); err != nil {
			return err
		}
	}
	return nil
}

func abs(fpath string) (string, error) {
	if strings.HasPrefix(fpath, ".") {
		dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return "", err
		}
		if fpath == "." {
			return dir, nil
		} else {
		}
		distDir, err := filepath.Abs(dir + fpath[1:])
		if err != nil {
			return "", err
		}
		return distDir, nil
	}
	dir, err := filepath.Abs(fpath)
	if err != nil {
		return "", err
	}
	return dir, nil
}

func (app *App) BuildGOPATH() (string, error) {
	var GOPATH string
	if len(app.Info.Repos) == 0 {
		GOPATH = app.GOPATH
	} else {
		path := strings.Join(app.Info.Repos, ":")
		if err := os.Setenv(GOPATHS_ENV_ORIGINAL_GOPATH, path); err != nil {
			return "", err
		}
		if app.GOPATH == "" {
			GOPATH = path
		} else {
			GOPATH = fmt.Sprintf("%s:%s", app.GOPATH, path)
		}
	}
	return GOPATH, nil
}

// Init is `gopaths init`.
// for initalizing GOPATH which gopaths maintain.
func (app *App) Init() error {
	gopath, err := app.BuildGOPATH()
	if err != nil {
		return err
	}
	if _, err := fmt.Printf("export GOPATH=%s", gopath); err != nil {
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
	gopath, err := app.BuildGOPATH()
	if err != nil {
		return err
	}
	return os.Setenv("GOPATH", gopath)
}

// Disable is `gopaths disable`.
// for disable gopaths we set and recover original gopath.
func (app *App) Disable(config *AppDisableConfig) error {
	return os.Setenv("GOPATH", app.Info.GOPATH)
}

// Add is `gopaths add`.
// for adding path to gopaths manage.
func (app *App) Add(config *AppAddConfig) error {
	fileGOPATHS := app.ReposPath + GOPATHS_GOPATHS_FILE
	rinfo, err := getInfo(fileGOPATHS)
	if err != nil {
		return err
	}
	var repos []string = app.Info.Repos
	for _, path := range config.Paths {
		pathAbs, err := abs(path)
		if err != nil {
			return err
		}
		if !isElem(pathAbs, repos) {
			repos = append(repos, pathAbs)
		}
	}
	rinfo.Repos = repos
	return setInfo(fileGOPATHS, rinfo)
}

func isElem(dist string, paths []string) bool {
	for _, p := range paths {
		if strings.Compare(dist, p) == 0 {
			return true
		}
	}
	return false
}

// Remove is `gopaths remove`
// for removing path to gopaths manage.
func (app *App) Remove(config *AppRemoveConfig) error {
	var repos []string
	for _, r := range app.Info.Repos {
		if !isElem(r, config.Paths) {
			repos = append(repos, r)
		}
	}
	app.Info.Repos = repos
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
