package main

const (
	expectedString        = "Forwarding from"
	expectedExitCode      = 1
	defaultTimeoutSeconds = 60
)

type applicationOpts struct {
	OkctlBinaryPath     string
	ClusterManifestPath string
	DatabaseName        string
}

type clusterManifestDatabasePostgres struct {
	Name string `json:"name"`
}

type clusterManifestDatabases struct {
	Postgres []clusterManifestDatabasePostgres `json:"postgres"`
}

type clusterManifest struct {
	Databases clusterManifestDatabases `json:"databases"`
}
