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

)
