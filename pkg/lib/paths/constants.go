package paths

const (
	// DefaultFilePermissions defines the default permission bits for okctl created files.
	// 0o600 means read and write for the file/folder owner only.
	DefaultFilePermissions = 0o600
	// DefaultDirectoryPermissions defines the default permission bits for okctl created directories.
	// 0o700 means read, write and execute for the file/folder owner only. Execute as a directory permission means
	// traverse the files and folders inside.
	DefaultDirectoryPermissions = 0o700
	// DefaultClusterArgoCDConfigDirectoryName defines the name of the directory where we configure ArgoCD for a
	// specific environment
	DefaultClusterArgoCDConfigDirectoryName = "argocd"
	// DefaultClusterArgoCDApplicationsDirectoryName defines the name of the directory where we place all ArgoCD
	// application manifests that should be automatically applied
	DefaultClusterArgoCDApplicationsDirectoryName = "applications"
	// DefaultClusterOkctlConfigurationDirectoryName defines the name of the directory where we place configuration files
	// that okctl manage
	DefaultClusterOkctlConfigurationDirectoryName = "okctl"
)
