package main

import (
	"fmt"
	"os"

	"github.com/google/uuid"

	"sigs.k8s.io/yaml"
)

func (receiver applicationOpts) Validate() error {
	_, err := os.Stat(receiver.OkctlBinaryPath)
	if err != nil {
		return fmt.Errorf("stating okctl binary path: %w", err)
	}

	_, err = os.Stat(receiver.ClusterManifestPath)
	if err != nil {
		return fmt.Errorf("stating cluster manifest path: %w", err)
	}

	rawManifest, err := os.ReadFile(receiver.ClusterManifestPath)
	if err != nil {
		return fmt.Errorf("reading cluster manifest: %w", err)
	}

	manifest := clusterManifest{}

	err = yaml.Unmarshal(rawManifest, &manifest)
	if err != nil {
		return fmt.Errorf("parsing cluster manifest: %w", err)
	}

	if !containsDatabase(manifest, receiver.DatabaseName) {
		return fmt.Errorf("%s does not exist in the cluster manifest: %w", err)
	}

	return nil
}

func containsDatabase(manifest clusterManifest, databaseName string) bool {
	for _, db := range manifest.Databases.Postgres {
		if db.Name == databaseName {
			return true
		}
	}

	return false
}

func generatePassfile() (string, error) {
	rawPassword := []byte(uuid.New().String())
	passPath := "pass.txt"

	err := os.WriteFile(passPath, rawPassword, 0x600)
	if err != nil {
		return "", fmt.Errorf("writing password to file: %w", err)
	}

	return passPath, nil
}

func crashPrint(err error) {
	fmt.Printf("Error: %s\n", err.Error())

	os.Exit(1)
}
