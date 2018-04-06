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

	PathInfo struct {
	}
)

func getInfo(fpath string) (*PathInfo, error) {
	var pinfo PathInfo
	if _, err := toml.DecodeFile(fpath, &pinfo); err != nil {
		return nil, errors.Wrap(err, "toml")
	}
	return &pinfo, nil
}
