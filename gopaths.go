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
