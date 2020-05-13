package login

import (
	"github.com/apex/log"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/oslokommune/okctl/pkg/validate"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/oslokommune/okctl/pkg/config"
	"github.com/oslokommune/okctl/pkg/login"
	"github.com/oslokommune/okctl/pkg/stager"
	"github.com/oslokommune/okctl/pkg/storage"
)

func BuildLoginCommand(_ *config.UserConfig, _ *logrus.Logger) *cobra.Command {
	return &cobra.Command{
		Use:   "login",
		Short: "Login to AWS",
		Run: func(cmd *cobra.Command, args []string) {
			Login()
		},
	}
}

func Login() {
	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		log.Fatalf("couldn't retrieve home directory")
	}

	cfg, err := config.LoadUserConfiguration(home)
	if err != nil {
		log.Fatalf("failed to load configuration file: %s", err)
	}

	err = validate.Host(cfg.Host)
	if err != nil {
		log.Fatalf("invalid host configuration: %s", err)
	}

	store := storage.NewFileSystemStorage(home)

	stagers, err := stager.FromConfig(cfg.Binaries, cfg.Host, store)
	if err != nil {
		log.Fatalf("failed to create binary stagers: %s", err)
	}

	for _, s := range stagers {
		err = s.Run()
		if err != nil {
			log.Fatalf("failed to stage binary: %s", err)
		}
	}

	err = login.New("922935510846", "byr299604").Login()
	if err != nil {
		log.Fatalf("failed to login: %s", err)
	}

	log.Info("Done")
}
