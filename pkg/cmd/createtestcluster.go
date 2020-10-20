package cmd

// CreateTestClusterEndMsg is the message shown to the user after creating a test cluster
const CreateTestClusterEndMsg = `Congratulations, your {{ .KubernetesCluster }} is now up and running.
To get started with some basic interactions, you can paste the
following exports into a terminal:

{{ .Exports }}

You can retrieve these credentials at any point by issuing the
command below, from within this repository:

$ okctl show credentials {{ .Environment }}

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
`
