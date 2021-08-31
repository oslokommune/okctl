// Package constant contains constants used throughout okctl
package constant

// nolint: golint
const (
	AddSchemeError       = "adding scheme: %w"
	BuildRestClientError = "building rest client: %w"

	NamespaceValidation               = "must consist of 3-64 characters (a-z, A-Z, -)"
	AwsAccountIDValidation            = "must consist of 12 digits"
	ClusterRootDomainValidation       = "with automatizeZoneDelegation enabled, must end with auto.oslo.systems"
	ClusterNameValidation             = "must consist of 3-64 characters (a-z, A-Z, -)"
	ClusterSupportedRegionsValidation = "for now, only \"eu-west-1\", \"eu-central-1\" and \"eu-north-1\" are supported"
	UserNameValidation                = "username must be in the form: yyyXXXXXX (y = letter, x = digit)"

	ErrorDescribingCert      = "describing certificate: %w"
	ErrorStoringAppResources = "storing application resources: %w"
	ErrorRemovingApp         = "removing application: %w"

	ReconcilerDetermineActionError = "determining course of action: %w"

	DeleteNotImplementedError = "deletion of applications is not implemented"
	ActionNotImplementedError = "action %s is not implemented"

	GetPrimaryHostedZoneError = "getting primary hosted zone: %w"
	GetGithubInfoError        = "retrieving Github information"

	GetContainerRepoError = "getting container repository: %w"

	GetDependencyStateError         = "acquiring dependency state: %w"
	DeterminePrimaryHostedZoneError = "determining existence of primary hosted zone for %s: %w"
	DetermineGithubRepoError        = "determining existence of a Github repository for %s: %w"
	DetermineExistsEcrRepoError     = "determining existence of a ECR repository: %w"
	ApplicationURLParseError        = "parsing application URL: %w"

	InferApplicationError     = "inferring application from stdin or file: %w"
	OptionValidationerror     = "failed validating options: %w"
	SpinnerCreationError      = "error creating spinner: %w"
	ReconcileApplicationError = "reconciling application: %w"

	OpeningApplicationFileError = "opening application file: %w"
	ReadingApplicationFileError = "reading application file: %w"

	ParisingApplicationYamlError = "parsing application yaml: %w"

	CreatePresistantVolumeClaimResourceError = "creating PersistentVolumeClaim resource: %w"

	CreateDeploymentPatchError = "creating deployment patch: %w"
	CreateIngressPatchError    = "creating ingress patch: %w"
	MarshalKustomizationError  = "marshalling kustomization: %w"
	MarshalIngressPatchError   = "marshalling ingress patch: %w"

	InferClusterError = "inferring cluster: %w"

	ValidateClusterDeclarationError = "validating cluster declaration: %w"
	LoadApplicationDataError        = "loading application data: %w"

	InitializeOkctlError = "initializing okctl: %w"

	GetServicesError              = "error getting services: %w"
	SyncDeclarationWithStateError = "synchronizing declaration with state: %w"

	ReadFileError        = "unable to read file: %w"
	CopyReaderDataError  = "copying reader data: %w"
	UnmarshalBufferError = "unmarshalling buffer: %w"

	DeleteArgoCdError     = "deleting argocd: %w"
	DeleteGithubRepoError = "deleting github repository: %w"

	FetchDeployKeyError = "fetching deploy key: %w"

	GetIdentityPoolError  = "getting identity pool: %w"
	GetArgoCdTimeoutError = "got ArgoCD timeout: %w"

	CreateaArgoCdError = "creating argocd: %w"

	CheckIfArgoCdExistsError = "acquiring ArgoCD existence: %w"

	CheckDependenciesError = "checking dependencies: %w"

	QueryStateError = "querying state: %w"

	InvalidArnError = "not a valid arn: %s"
	ParseArnError   = "parsing arn: %w"

	StopSpinnerError = "stopping spinner: %w"

	GetRepositoryConfigError = "getting repository config: %w"

	LocatePostgresDatabaseError = "finding postgres database: %w"

	ValidateInputsError = "validating inputs: %w"

	BuildKubeconfigError = "building kubeconfig: %w"

	EnablePodEniError = "enabling pod eni: %w"

	CreateSecurityGroupPolicyError = "creating security group policy: %w"

	CreatePodError = "creating pod: %w"
	WatchPodError  = "watching pod: %w"

	AttachPodError = "attaching to pod: %w"

	CreateAutoScalerError = "creating autoscaler: %w"
	DeleteAutoScalerError = "deleting autoscaler: %w"

	CheckIfClusterExistsError    = "acquiring cluster existence: %w"
	CheckIfAutoScalerExistsError = "acquiring autoscaler existence: %w"

	GetHelmReleaseError = "getting Helm release: %w"

	NoValidCredentialsError = "no valid credentials: %s"

	PopulateRequiredFieldsError = "populating required fields: %w"

	SamlAssertionEmptyError = "empty SAML assertion"
	VerifySamlError         = "verifying saml: %w"

	CreateAwsStsSessionError = "creating AWS STS session: %w"

	RetrieveStsCretentialsWithSamlError = "retrieving STS credentials using SAML: %w"

	NoCredentialsAvailableError = "no credentials available"

	AwsSomethingError = "something"
	AwsSomethingBadError = "something bad"

	GetHelmError = "getting helm release: %w"
	GetVpcError = "getting vpc: %w"

	CreateAWSLoadBalancerControllerError = "creating aws load balancer controller: %w"
	DeleteAWSLoadBalancerControllerError = "deleting aws load balancer controller: %w"

	CheckIfAWSLoadBalancerExistsError = "acquiring AWS Load Balancer Controller existence: %w"
	CheckDeleteDependenciesError = "checking deletion dependencies: %w"

	CreateBlockStorageError = "creating blockstorage: %w"
	DeleteBlockStorageError = "deleting blockstorage: %w"
	LoadStateDatabaseError = "loading state database: %w"

	OutputExistsError = "already have output with name: %s"
	ResourceExistsError = "already have resource with name: %s"
	RemoveCertificateUsagesError = "removing usages of certificate: %w"

	BuildCloudFormationTemplateError = "building cloudformation template: %w"
	ApplyCloudformationTemplateError = "applying cloudformation template: %w"

	ProcessOutputsError = "processing outputs: %w"

	HasCapacityError = "%s: required %d, but only have %d available"

	CdirNotIpv4Error = "cidr (%s) is not of type IPv4"
	CdirAddressSpaceError = "address space of cidr (%s) is less than required: %d < %d"
	CdirNotInLegalRangeError = "provided cidr (%s) is not in the legal ranges: %s"

	ListLoadBalancersError = "listing load balancers: %w"
	DeleteLoadBalancerError = "deleting load balancer: %w"

	ListSecurityGroupsError = "listing security groups for vpc: %w"
	DeleteSecurityGroupError = "deleting security group: %w"
	ListingTargetGroupsError = "listing target groups: %w"
	DescribeTargetGroupsError = "describing tags for target group: %w"
	RemovingTargetGrupError = "removing target group: %w"

	CheckIfVpcExistsError = "checking VPC existence: %w"
	FailedToWriteDebugOutputError = "failed to write debug output: %w"

	FailedToMarshalDataError   = "failed to marshal data for"
	FailedToCreateRequestError = "failed to create request for"
	RequestFailedForError = "request failed for"
	ResponseFailedForError = "failed to read response for"
	FailedTowriteDebugOutputError = "failed to write debug output: %w"
	FailedToParseResponseError  ="failed to parse response: %w"
	FailedToWriteProgressForError = "failed to write progress for"

	UnmarshalOnServerSideError = "unmarshalling error from server side: %w: %s"
	ValidateDeserializedError = "validating deserialized error with content %s: %w"

	CreateClusterError = "creating cluster: %w"
	DeleteClusterError = "deleting cluster: %w"

	DeleteDanglingALBError = "cleaning up dangling ALBs: %w"
	CheckDepedencyReadyError = "checking for dependency ready status: %w"
	CleanUpALBError = "cleaning up ALBs: %w"
	CleanUpTragetGroupError = "cleaning target groups: %w"

	GetVPCStateError = "acquiring VPC state: %w"

	CheckServiceQuotasError = "checking service quotas: %w"

	AssertExistenceError = "asserting existence: %w"

	GetAwsIamAuthBinaryError = "retrieving aws-iam-authenticator binary: %w"
	GetKubectlBinaryError = "retrieving kubectl binary: %w"
	GetEksctlBinaryError = "retrieving eksctl binary: %w"

	FailedToDeleteClusterError = "failed to delete cluster: %w"

	QueryForClusterDataError = "querying state for cluster data: %w"

	DescribeUerPoolDomainError = "describing user pool domain: %w"

	DescribeUserPoolClientError = "describing user pool client: %w"

	EableTotpfMfaError = "enabling totp mfa: %w"

	SetPS1Error = "could not set PS1: %w"

	CreateExecutablePS1Error = "could not create PS1 executable: %w"

	CheckIfPS1ExecutableExistsError = "could not check existence of PS1 helper executable: %w"

	UnableToCreatePS1FileError = "couldn't create PS1 file: %w"
	WriteContentToPS1FileError = "could not write contents to ps1 file: %w"
	ClosePS1FileError = "could not close ps1 file: %w"

	CreateCloudFormationStackError = "creating cloud formation stack: %w"

	CollectStackFormationOutputsError = "collecting stack outputs: %w"

	PatchClouFormationTemplateError = "patching cloud formation template: %w"

	RotateHookInitializeError = "initialising the file rotate hook: %v"
	GetRepositoryRootDirectoryError = "getting repository root directory: %w"

	CreateContainerRepositoryError = "creating container repository: %w"
	DeleteContainerRepositoryError = "deleting container repository: %w"

	GetStateError = "acquiring existence from state %w"
	WriteVolumeToBufferError = "writing volume to buffer: %w"
	WriteResourceDefinitionToBufferError = "writing resource definition to buffer: %w"

	CanNotFindInArchiveError = "couldn't find: %s, in archive"

	InitializeError = "initialising: %w"

	CheckDepedencyError = "checking dependency: %w"

	CheckIfPrimaryHostedZoneExistsError = "checking primary hosted zone existence: %w"

	ReaderNilError = "reader is nil"

	UnsupportedDigesterError = "unsupported digester: %s"

	InvalidDomainError = "invalid domain: %s"

	HolyCrapError = "holy crap"

	UnhandledDNSReponseCodeError = "don't know how to handle DNS response code: %d"
	DomainAlreadyInUseError = "domain '%s' already in use, found DNS records"
	GetNSRecordsForDomainError = "unable to get NS records for domain '%s', does not appear to be delegated yet"
	NameServerNotMatchingError = "nameservers do not match, expected: %s, but got: %s"

	FailedToSetTTLOnNSRecordError = "failed to set NS record TTL: %w"

	GetSecurityGroupForNodeError = "getting security group for node: %w"

	AuthorizeSecurityGroupIngressError = "authorizing security group ingress: %w"

	RevokeSecurityGroupIngressError = "revoking security group ingress: %w"

	GetEIPQuotasError = "getting eip quotas: %w"

	GetEIPCountError = "getting current eip count: %w"

	GetFargateProfileError = "getting fargate profile: %w"

	FailedToCreateServiceAccountError = "failed to create service account: %s, because: %w"
	FailedToDeleteServiceAccountError = "failed to delete service account: %s, because: %w"

	FailedToDeleteError = "failed to delete: %s, because: %w"

	GetClusterInfoError = "failed to get cluster information: %s: %w"

	DescribeListenersError = "describing listeners: %w"

	DeleteListenersError = "deleting listener: %w"

	CreateExternalDNSError = "creating external DNS: %w"
	DeleteExteralDNSError  = "deleting external DNS: %w"

	CheckIfDNSControllerExistsError = "acquiring DNS controller existence: %w"

	CreateExternalSecretSetError = "creating external secrets client set: %w"

	ListExternalSecretsError = "listing external secrets in %s: %w"
	GetExternalSecretsError = "getting external secret %s in %s: %w"

	CreateExternalSecretsError = "creating external secrets: %w"
	DeleteExternalSecretsError = "deleting external secrets: %w"

	CheckIfSecretsControllerExistsError = "acquiring secrets controller existence: %w"


	GetFargateOnDemanPodsQuotasError = "getting fargate on-demand pods quotas: %w"
	GetFargetOnDemandPodUtilization = "getting fargate on-demand pods utilization: %w"

	PkgURLValidationError = "a valid pkgURL must begin https://, got: %s"
	DownloadFailedError = "failed to download file at: %s"
	BadStatusDownloadError = "bad status: %s, failed download of: %s"
	EmptyDownloadError = "downloaded file was size: 0, for url: %s"

	FileSystemCastError = "could not cast implemenation to *fileSystem"
	StructNotFoundError = "failed to retrieve struct: no such name '%s'"
	ProcessStructError = "failed to process struct: %w"
	PostProcessDataError = "failed to postprocess data: %w"

	FileExistsError = "file exists: %w"
	UnknownOptionError = "unknown option: %s"
	ApplyAlterationError = "failed to apply alteration %s: %w"
	PreProcessDataError = "failed to preprocess data: %w"

	StoreBytesError = "failed to store bytes: %w"

	CheckIfFileExistsError        = "failed to determine if file exists: %w"
	CannotOverwriteFileError      = "file '%s' exists and overwrite is disabled"
	CreateDirectoriesError        = "failed to create directories: %w"
	CheckIfDirectoryExistsError   = "failed to determine if directory exists: %w"
	CannotOverwriteDirectoryError = "directory does not exist '%s' and create directories disabled"
	ProcessUnkownOpreationError   = "cannot process unknown operation option: %v"

	WriteToFileError = "failed to write file: %w"
	RemoveFileError = "failed to remove file: %w"

	RemoveDirectoryError = "failed to remove directory: %w"

	FileSystemTaskCastError = "failed to cast task to *fileSystemTask"
	ProcessTaskError = "failed to process task %s(%s): %w"

	StartPortForwardingError = "starting port forwarding: %w"

	StagingRepositoryError      = "staging repository: %w"
	GetWorkTreeError            = "getting work tree: %w"
	CheckoutBranchError         = "checking out branch: %w"
	PullBranchError             = "pulling branch: %w"
	RunActionError              = "running action: %w"
	CheckStatusError            = "checking status: %w"
	CommitNameServerRecordError = "committing nameserver record: %w"
	PushToRemoteError           = "pushing to remote: %w"
	RemoveTrackedFileError      = "removing file: %w"
	AddFileError = "adding file: %w"
	GetGitStatusError = "getting status: %w"
	FileNotStagedAsDeletedError = "file: %s, not staged as deleted"
	CreateFileError = "creating file: %w"
	CloneRepositoryError = "cloning repository: %w"

	InitializeRepositoryError = "initializing repository: %w"

	BuildTokenVerificationRequestError = "failed to build token verification request: %w"
	SendTokenVerificationRequestError = "failed to send token verification request: %w"
	ValidationHTTPError = "HTTP error %v (%v) when requesting token validation"
	InvalidAuthenticatorCredentilasError = "authenticator[%d]: invalid credentials, because: %w"

	RecieveDeviceCodeError = "failed to retrieve device code: %w"
	SurveyFailedError = "survey failed: %w"

	DeviceAuthorizationError = "failed getting device authorization: %w"

	UserNotReadyError = "user was not ready to continue: %w"

	SerializeCredentialsError = "failed to serialize credentials"
	MissingCredentialsError = "no credentials exist"

	GetGithubCredentialsError = "failed to get github credentials: %w"
	GetGithubTeamsError = "failed to retrieve teams: %w"
	GetGithubRepositoriesError = "failed to retrieve repositories: %w"
	CreateDeployKeyError = "creating deploy key: %w"
	DeleteDeployKeyError = "deleting deploy key: %w"
	CreateGithubPullRequestError = "creating github pull request: %w"
	AddLabelToPullRequestError = "adding labels to pull request: %w"

	GetAuthenticationDetailsError = "getting authentication details as environment: %w"
	EstablishEnvironmentError = "establishing environment: %w"
	InitializeActionConfirgurationError = "initializing action configuration: %w"
	FindReleaseError = "finding release: %w"
	BadReleaseStateError = "release is in state: %s, cannot continue"
	LocateChartError = "locating chart: %w"
	LoadChartError = "loading chart: %w"
	CheckIfChartIsInstallableError = "checking if chart is installable: %w"
	UpdateLocalChartsDirectoryError = "updating local charts directory: %w"
	GenerateValuesMap = "generating values map: %w"
	CreateDebuggerError = "creating debugger: %w"
	DebugNamespaceError = "debugging namespace: %w"
	RunHelmInstallCommandError = "running helm install command: %w"
	ChartNotInstallableError = "chart: %s is not installable"
	CreateLockError = "failed to create lock: %s"
	MarshallToYamlError = "marshalling values struct to yaml: %w"

)
