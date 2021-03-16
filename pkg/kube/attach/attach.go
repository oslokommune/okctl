// Package attach knows how to interact with a running pod
package attach

import (
	"io"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

// Attach contains the required state
// for connecting to a pod
type Attach struct {
	Client kubernetes.Interface
	Config *restclient.Config
}

// New returns an initialised attacher
func New(client kubernetes.Interface, config *restclient.Config) *Attach {
	return &Attach{
		Client: client,
		Config: config,
	}
}

// Run a given command
func (a *Attach) Run(podName string, namespace, command string, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	cmd := []string{
		"sh",
		"-c",
		command,
	}

	req := a.Client.CoreV1().RESTClient().
		Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec")

	option := &v1.PodExecOptions{
		Command: cmd,
		Stdin:   true,
		Stdout:  true,
		Stderr:  true,
		TTY:     true,
	}

	if stdin == nil {
		option.Stdin = false
	}

	req.VersionedParams(
		option,
		scheme.ParameterCodec,
	)

	exec, err := remotecommand.NewSPDYExecutor(a.Config, "POST", req.URL())
	if err != nil {
		return err
	}

	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	})
	if err != nil {
		return err
	}

	return nil
}
