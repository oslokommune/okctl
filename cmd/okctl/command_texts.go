package main

//nolint
const (
	ApplyCommandsShortDescription    = "Apply commands"
	DeleteCommandsShortDescription   = "Delete commands"
	ScaffoldCommandsShortDescription = "scaffold templates for different resources"
	ShowCommandsShortDescription     = "Show commands"

	ApplyClusterShortDescription = "apply a cluster definition to the world"
	ApplyClusterLongDescription  = "ensures your cluster reflects the declaration of it"

	ApplyApplicationShortDescription = "Applies an application.yaml to the IAC repo"

	DeleteClusterShortDescription = "Delete a cluster"
	DeleteClusterLongDescription  = `Delete all resources related to an EKS cluster,
including VPC, this is a highly destructive operation.`

	deleteApplicationShortDescription = "Delete an application"
	deleteApplicationLongDescription  = "Delete all resources related to an application"

	ForwardCommandsShortDescription = "Forward commands"
	ForwardCommandsLongDescription  = `Forward traffic to different services running in AWS

Sometimes it can be useful to forward traffic to services running
in AWS, such as a Postgres database, to inspect the schema 
or perform other simple operations.
`
	ForwardPostgresShortDescription = "Forward to the given postgres database"
	ForwardPostgresLongDescription  = `This sets up PgBouncer in the Kubernetes cluster, which
makes it possible to operate on a Postgres database running 
on AWS via a local port.

Be very mindful if using this functionality with a production
cluster. If you issue the wrong commands you can potentially
delete important data.
`
	OkctlShortDescription = "Opinionated and effortless infrastructure and application management"
	OkctlLongDescription  = `A highly opinionated CLI for creating a Kubernetes cluster in AWS with
a set of applications that ensure tighter integration between AWS and
Kubernetes, e.g., aws-alb-ingress-controller, external-secrets, etc.

Also comes pre-configured with ArgoCD for managing deployments, etc.
We also use the prometheus-operator for ensuring metrics and logs are
being captured. Together with slack and slick.`

	ScaffoldShortDescription = "Scaffold an application template"
	ScaffoldLongDescription  = "Scaffolds an application.yaml template which can be used to produce necessary Kubernetes and ArgoCD resources"

	ScaffoldClusterShortDescription = "Scaffold cluster resource template"
	ScaffoldClusterLongDescription  = "Scaffolds a cluster resource which can be used to control cluster resources"

	ShowShortDescription = "Show the location of the credentials"
	ShowLongDescription  = `This makes it possible to source the output from this command to run with kubectl`

	VenvShortDescription = "Runs a virtual environment"
	VenvLongDescription  = `Runs a sub shell with all needed environmental variables set.

The variables are the same as shown in "okctl show credentials". The shell command to run is retrieved from the first
environment variable that is set of the following: $OKCTL_SHELL, $SHELL. If none is set, "/bin/sh" is used.

So to override, you can run for instance:

export OKCTL_SHELL=/bin/bash
okctl venv -c my-cluster.yaml
`
	VersionShortDescription = "Prints version info"
	VersionLongDescription  = `Prints version number and git-sha of okctl-build.`

	CompletionShortDescription = "Generate completion script"
	CompletionLongDescription  = `To load completions:

Bash:

	$ source <(okctl completion bash)

	# To load completions for each session, execute once:
	# Linux:
	$ sudo sh -c 'okctl completion bash > /etc/bash_completion.d/okctl'
	# macOS:
	$ okctl completion bash > /usr/local/etc/bash_completion.d/okctl

Zsh:

	# If shell completion is not already enabled in your environment,
	# you will need to enable it.  You can execute the following once:

	$ echo "autoload -U compinit; compinit" >> ~/.zshrc

	# To load completions for each session, execute once:
	$ okctl completion zsh > "${fpath[1]}/_okctl"

	# You will need to start a new shell for this setup to take effect.

fish:

	$ okctl completion fish | source

	# To load completions for each session, execute once:
	$ okctl completion fish > ~/.config/fish/completions/okctl.fish

PowerShell:

	PS> okctl completion powershell | Out-String | Invoke-Expression

	# To load completions for every new session, run:
	PS> okctl completion powershell > okctl.ps1
	# and source this file from your PowerShell profile.`
)
