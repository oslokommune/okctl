package config

import (
	"os"
	"path"
)

func Load() (*AppConfig, *RepoConfig, error) {
	home, err := GetHomeDir()
	if err != nil {
		return nil, nil, err
	}

	appCfg, err := LoadApp(home, NewDefaultAppCfgFn(), NewDefaultNotFoundFn())
	if err != nil {
		return nil, nil, err
	}

	appCfg.BaseDir = path.Join(home, DefaultAppDir)

	err = appCfg.Host.Valid()
	if err != nil {
		return nil, nil, err
	}

	wd, err := os.Getwd()
	if err != nil {
		return nil, nil, err
	}

	repoCfg, err := LoadRepo(wd)
	if err != nil {
		return nil, nil, err
	}

	return appCfg, repoCfg, nil
}
