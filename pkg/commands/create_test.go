package commands

import (
	"testing"

	"github.com/andreyvit/diff"
	"github.com/bmizerany/assert"
	"github.com/logrusorgru/aurora"
)

// The weird characters in this variable are color codes from the color library.
// If you need to find out the actual content to put into the expected value, use
// ioutil.WriteFile("/tmp/create_test.txt", []byte(createClusterExpectedMsg), 0644)
// nolint:stylecheck
const createClusterExpectedMsg = `Congratulations, your [32mkubernetes cluster[0m is now up and running.
To get started with some basic interactions, you can paste the
following exports into a terminal:

export HELM_CACHE_HOME=/home/johndoe/.okctl/helm
export KUBECONFIG=/home/johndoe/.okctl/credentials/test/kubeconfig

You can retrieve these credentials at any point by issuing the
command below, from within this repository:

$ okctl show credentials prod

Tip: Run [32mokctl venv[0m to run a shell with these environment variables set. Then you
can avoid using full paths to executables and modifying your PATH.

Now you can use [32mkubectl[0m to list nodes, pods, etc. Try out some commands:

$ /home/johndoe/.okctl/binaries/kubectl/1.16.8/linux/amd64/kubectl get pods --all-namespaces
$ /home/johndoe/.okctl/binaries/kubectl/1.16.8/linux/amd64/kubectl get nodes

This also requires [32maws-iam-authenticator[0m, which you can add to your PATH from here:

/home/johndoe/.okctl/binaries/aws-iam-authenticator/0.5.1/linux/amd64/aws-iam-authenticator

Optionally, install kubectl and aws-iam-authenticator to your
system from:

- https://kubernetes.io/docs/tasks/tools/install-kubectl/
- https://docs.aws.amazon.com/eks/latest/userguide/install-aws-iam-authenticator.html

The installed version of kubectl needs to be within 2 versions of the
kubernetes cluster version, which is: [32m1.17[0m.

We have also setup [32mArgoCD[0m for continuous deployment, you can access
the UI at this URL by logging in with Github:

http://argocd

It might take 5-10 minutes for the ArgoCD ALB to come up, and
about 15 minutes for the auth to come up.
`

func TestCreateClusterMessage(t *testing.T) {
	t.Run("Should get expected output", func(t *testing.T) {
		exports := `export HELM_CACHE_HOME=/home/johndoe/.okctl/helm
export KUBECONFIG=/home/johndoe/.okctl/credentials/test/kubeconfig`
		data := CreateClusterMsgOpts{
			KubernetesCluster:       aurora.Green("kubernetes cluster").String(),
			Exports:                 exports,
			Environment:             "prod",
			VenvCmd:                 aurora.Green("okctl venv").String(),
			KubectlCmd:              aurora.Green("kubectl").String(),
			AwsIamAuthenticatorCmd:  aurora.Green("aws-iam-authenticator").String(),
			KubectlPath:             kubectlPath,
			AwsIamAuthenticatorPath: awsIamAuthenticatorPath,
			K8sClusterVersion:       aurora.Green("1.17").String(),
			ArgoCD:                  aurora.Green("ArgoCD").String(),
			ArgoCDURL:               "http://argocd",
		}

		msg, err := GoTemplateToString(CreateClusterEndMsg, data)

		assert.Equal(t, nil, err)
		assert.Equal(t, createClusterExpectedMsg, msg, diff.LineDiff(createClusterExpectedMsg, msg))
	})
}
