package commands

// ShowMessageOpts contains the fields used by the Go template for the show credentials user message
type ShowMessageOpts struct {
	VenvCmd                 string
	KubectlCmd              string
	KubectlPath             string
	AwsIamAuthenticatorCmd  string
	AwsIamAuthenticatorPath string
	K8sClusterVersion       string
	ArgoCD                  string
	ArgoCDURL               string
}

// ShowMsg is the message shown to the user when using show credentials
const ShowMsg = `
Tip: Run {{ .VenvCmd }} to run a shell with these environment variables set. Then you
can avoid using full paths to executables and modifying your PATH.

Now you can use {{ .KubectlCmd }} to list nodes, pods, etc. Try out some commands:

$ {{ .KubectlPath }} get pods --all-namespaces
$ {{ .KubectlPath }} get nodes

This also requires {{ .AwsIamAuthenticatorCmd }}, which you can add to your PATH from here:

{{ .AwsIamAuthenticatorPath }}

Optionally, install kubectl and aws-iam-authenticator to your
system from:

- https://kubernetes.io/docs/tasks/tools/install-kubectl/
- https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html

The installed version of kubectl needs to be within 2 versions of the
kubernetes cluster version, which is: {{ .K8sClusterVersion }}.

We have also setup {{ .ArgoCD }} for continuous deployment, you can access
the UI at this URL by logging in with Github:

{{ .ArgoCDURL }}

`
