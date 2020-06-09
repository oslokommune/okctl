package config

import (
	"bytes"
	"io"
	"os"
	"path"

	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/versent/saml2aws/pkg/prompter"
	"gopkg.in/yaml.v2"
)

func LoadAppCfg() (*AppConfig, error) {
	home, err := GetHomeDir()
	if err != nil {
		return nil, err
	}

	appCfg, err := LoadApp(home, NewDefaultAppCfgFn(), NewDefaultNotFoundFn())
	if err != nil {
		return nil, err
	}

	appCfg.BaseDir = path.Join(home, DefaultAppDir)

	err = appCfg.Host.Valid()
	if err != nil {
		return nil, err
	}

	return appCfg, nil
}

func NewAppCfg() (*AppConfig, error) {
	home, err := GetHomeDir()
	if err != nil {
		return nil, err
	}

	store := storage.NewFileSystemStorage(home)

	writer, err := store.Create(DefaultAppDir, DefaultAppConfig)
	if err != nil {
		return nil, err
	}

	defer func() {
		err = writer.Close()
	}()

	cfg := NewDefaultAppCfg()

	cfg.User.Username = prompter.StringRequired("Username (byrXXXXXX): ")

	b, err := yaml.Marshal(cfg)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(writer, bytes.NewReader(b))

	return cfg, err
}

func LoadRepoCfg() (*RepoConfig, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	repoCfg, err := LoadRepo(wd)
	if err != nil {
		return nil, err
	}

	return repoCfg, nil
}
