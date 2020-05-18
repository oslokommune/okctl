package configure

import (
	"bytes"
	"io"

	"github.com/spf13/cobra"
	"github.com/versent/saml2aws/pkg/prompter"
	"gopkg.in/yaml.v2"

	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/storage"
)

func BuildConfigureCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "configure",
		Short: "Simplify the usage of okctl",
		Long: `This will help you configure okctl so you don't
need to repeat the same information over and over.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// We need to override the default PersistentPreRunE function, because
			// we want to set the configuration.
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return Configure()
		},
	}
}

func Configure() error {
	home, err := config.GetHomeDir()
	if err != nil {
		return err
	}

	store := storage.NewFileSystemStorage(home)

	writer, err := store.Create(config.DefaultAppDir, config.DefaultAppConfig)
	if err != nil {
		return err
	}

	defer func() {
		err = writer.Close()
	}()

	cfg := config.NewDefaultAppCfg()

	cfg.User.Username = prompter.StringRequired("Username (byrXXXXXX): ")

	b, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, bytes.NewReader(b))

	return err
}
