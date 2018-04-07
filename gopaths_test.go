package gopaths

import (
	"os"
	"testing"
)

func TestApp(t *testing.T) {
	testApp := &App{
		GOPATH:    "",
		Info:      &RepoInfo{},
		ReposPath: "",
	}
	if err := testApp.Init(); err != nil {
		t.Errorf("`gopaths init` : %v", err)
	}
	if err := testApp.Clean(); err != nil {
		t.Errorf("`gopaths clean`: %v", err)
	}
	if v := os.Getenv(GOPATHS_ENV_ORIGINAL_GOPATH); v != "" {
		t.Errorf("%s is not clean: %s", GOPATHS_ENV_ORIGINAL_GOPATH, v)
	}
}
