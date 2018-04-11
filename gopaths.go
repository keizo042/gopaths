package gopaths

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

type RepoVersion int16

func (v RepoVersion) String() string {
	i := int16(v)
	major := (i & 0xd0) >> 6
	minor1 := (i & 0x30) >> 4
	minor2 := (i & 0x0d) >> 2
	release := i & 0x03

	return fmt.Sprintf("%s:%s:%s", major, minor1, minor2, release)
}

var (
	GOPATHS_ENV_ORIGINAL_GOPATH = "GOPATHS_GOPATH"
	GOPATHS_CONFIG_FILE         = "config.toml"
	GOPATHS_GOPATHS_FILE        = "gopaths.toml"

	REPOINFO_VERSION_NUMBER RepoVersion = 0x00000100 // 0.0.1.0

	ERR_NOTIMPL = errors.New("NotImplemented")
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
		Path    bool
	}

	AppConfigConfig struct {
		GOPATH  string
		Args    []string
		Show    bool
		Verbose bool
	}

	AppAddConfig struct {
		Paths   []string
		Verbose bool
	}

	AppRemoveConfig struct {
		Paths   []string
		Verbose bool
		All     bool
	}

	AppEnableConfig struct {
		Verbose bool
		Args    []string
	}

	AppDisableConfig struct {
		Verbose bool
		Args    []string
	}

	RepoInfo struct {
		Version      RepoVersion
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
	oldFpath := fpath + ".old"
	if err := os.Rename(fpath, oldFpath); err != nil {
		return err
	}
	f, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := toml.NewEncoder(f).Encode(rinfo); err != nil {
		return err
	}
	return os.Remove(oldFpath)
}

func NewApp(c *Config) (*App, error) {
	if err := checkGopathsConfig(c.SettingPath); err != nil {
		return nil, errors.Wrap(err, "check global config dotfile")
	}
	var gopath string = os.Getenv("GOPATH")
	if gopath == "" {
		path, err := exec.Command("go", "env", "GOPATH").Output()
		if err != nil {
			return nil, errors.Wrap(err, "fail go command")
		}
		gopath = string(path)
		if strings.HasSuffix(gopath, "\n") {
			gopath = gopath[:len(gopath)-1]
		}
	}
	info, err := getInfo(c.SettingPath + GOPATHS_GOPATHS_FILE)
	if err != nil {
		return nil, errors.Wrap(err, "repos info")
	}
	if info.GOPATH == "" {
		info.GOPATH = gopath
	}
	if info.Version == 0 {
		info.Version = REPOINFO_VERSION_NUMBER
	}
	return &App{
		GOPATH:    gopath,
		ReposPath: c.SettingPath,
		Info:      info,
	}, nil
}

func checkGopathsConfig(fpath string) error {
	if _, err := os.Stat(fpath); err != nil {
		if err := os.MkdirAll(fpath, 0777); err != nil {
			return err

		}
	}
	fileGOPATHS := fpath + GOPATHS_GOPATHS_FILE
	if _, err := os.Stat(fileGOPATHS); err != nil {
		if _, err := os.Create(fileGOPATHS); err != nil {
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
	// TODO: take care ".." relative path.
	if strings.HasPrefix(fpath, ".") {
		dir, err := os.Getwd()
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

func (app *App) SetGOPATH(gopath string) error {
	return os.Setenv("GOPATH", gopath)
}

// Init is `gopaths init`.
// for initalizing GOPATH which gopaths maintain.
func (app *App) Init(config *AppInitConfig) error {
	gopath, err := app.BuildGOPATH()
	if err != nil {
		return err
	}
	if config.Path {
		_, err := fmt.Printf("export GOPATH=%s", gopath)
		return err
	}
	return app.SetGOPATH(gopath)
}

// Config is `gopaths config`
// for manage gopaths config.
func (app *App) Config(config *AppConfigConfig) error {
	if config.Show {
		return toml.NewEncoder(os.Stdout).Encode(app.Info)
	}
	return ERR_NOTIMPL
}

// Enable is `gopaths enable`.
// for enable gopaths we set.
func (app *App) Enable() error {
	gopath, err := app.BuildGOPATH()
	if err != nil {
		return errors.Wrap(err, "building gopath")
	}
	return app.SetGOPATH(gopath)
}

// Disable is `gopaths disable`.
// for disable gopaths we set and recover original gopath.
func (app *App) Disable(config *AppDisableConfig) error {
	return app.SetGOPATH(app.Info.GOPATH)
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
	if err := setInfo(fileGOPATHS, rinfo); err != nil {
		return err
	}
	gopath, err := app.BuildGOPATH()
	if err != nil {
		return err
	}
	return app.SetGOPATH(gopath)
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
	var newGopathRepos []string
	var removedPaths []string
	fileGOPATHS := app.ReposPath + GOPATHS_GOPATHS_FILE
	if config.All {
		app.Info.Repos = []string{}
		return setInfo(fileGOPATHS, app.Info)
	}

	for _, path := range config.Paths {
		absPath, err := abs(path)
		if err != nil {
			return err
		}
		removedPaths = append(removedPaths, absPath)
	}
	for _, curRepos := range app.Info.Repos {
		if isElem(curRepos, removedPaths) {
			continue
		}
		newGopathRepos = append(newGopathRepos, curRepos)
	}
	app.Info.Repos = newGopathRepos
	if err := setInfo(fileGOPATHS, app.Info); err != nil {
		return err
	}
	// TODO: fix some configura
	gopath, err := app.BuildGOPATH()
	if err != nil {
		return err
	}
	return app.SetGOPATH(gopath)
}

func (app *App) Restore() error {
	return ERR_NOTIMPL
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
