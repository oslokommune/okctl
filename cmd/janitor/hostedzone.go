package main

import (
	"fmt"
	"os"

	"github.com/oslokommune/okctl/pkg/cloud"
	"github.com/oslokommune/okctl/pkg/domain"
	"github.com/oslokommune/okctl/pkg/janitor/hostedzone"
	"github.com/spf13/cobra"
)

func buildHostedZoneCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "hostedzone",
		Short: "Interacts with route53 hosted zones",
	}

	cmd.AddCommand(buildHostedZoneUndelegatedNameServerRecords())

	return cmd
}

func buildHostedZoneUndelegatedNameServerRecords() *cobra.Command {
	var hostedZoneID string

	cmd := &cobra.Command{
		Use:   "undelegated",
		Short: "Returns a list of undelegated hosted zones",
		RunE: func(cmd *cobra.Command, args []string) error {
			sess, err := cloud.NewSessionFromEnv("eu-west-1")
			if err != nil {
				return err
			}

			provider, err := cloud.NewFromSession("eu-west-1", "", sess)
			if err != nil {
				return err
			}

			hzID, err := cmd.Flags().GetString("hosted-zone-id")
			if err != nil {
				return err
			}

			undelegated, err := hostedzone.New(provider.Provider).UndelegatedZonesInHostedZones(hzID, domain.NameServers)
			if err != nil {
				return err
			}

			_, err = fmt.Fprintln(os.Stdout, "Undelegated hosted zones:")
			if err != nil {
				return err
			}

			for _, u := range undelegated {
				_, err := fmt.Fprintf(os.Stdout, "- %s\n", u.Name)
				if err != nil {
					return err
				}
			}

			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&hostedZoneID, "hosted-zone-id", "i", "",
		"The id of the hosted zone to evaluate")

	return cmd
}
