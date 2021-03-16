package psqlclient

import (
	"context"
	"fmt"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/oslokommune/okctl/pkg/kube/attach"

	restclient "k8s.io/client-go/rest"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// PSQLClient contains the state required for
// creating a psql client pod
type PSQLClient struct {
	Name      string
	Namespace string
	Manifest  *v1.Pod

	Client kubernetes.Interface
	Config *restclient.Config
	Ctx    context.Context
}

// New returns an initialised client for interacting with pods
func New(name, namespace string, pod *v1.Pod, clientSet kubernetes.Interface, config *restclient.Config) *PSQLClient {
	return &PSQLClient{
		Name:      name,
		Namespace: namespace,
		Manifest:  pod,
		Client:    clientSet,
		Config:    config,
		Ctx:       context.Background(),
	}
}

// Create creates a k8s pod with psql client
func (c *PSQLClient) Create() (*v1.Pod, error) {
	p, err := c.Client.CoreV1().Pods(c.Namespace).Create(c.Ctx, c.Manifest, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}

	return p, nil
}

// Watch waits until the pod is running
func (c *PSQLClient) Watch(resp *v1.Pod) error {
	status := resp.Status

	w, err := c.Client.CoreV1().Pods(c.Namespace).Watch(c.Ctx, metav1.ListOptions{
		Watch:           true,
		ResourceVersion: resp.ResourceVersion,
		FieldSelector:   fields.OneTermEqualSelector("metadata.name", c.Name).String(),
		LabelSelector:   labels.Everything().String(),
	})
	if err != nil {
		return err
	}

	func() {
		for {
			select {
			case events, ok := <-w.ResultChan():
				if !ok {
					return
				}
				resp = events.Object.(*v1.Pod)
				status = resp.Status
				if resp.Status.Phase != v1.PodPending {
					w.Stop()
				}
			case <-time.After(30 * time.Second):
				w.Stop()
			}
		}
	}()

	if status.Phase != v1.PodRunning {
		return fmt.Errorf("waiting for pod: %v", status.Phase)
	}

	return nil
}

// Delete the psql client pod
func (c *PSQLClient) Delete() error {
	policy := metav1.DeletePropagationForeground

	return c.Client.CoreV1().Pods(c.Namespace).Delete(c.Ctx, c.Name, metav1.DeleteOptions{
		PropagationPolicy: &policy,
	})
}

// Attach to the psql pod and hook up all the stds (pun intended)
func (c *PSQLClient) Attach() error {
	return attach.New(c.Client, c.Config).Run(
		c.Name,
		c.Namespace,
		"psql",
		os.Stdin,
		os.Stdout,
		os.Stderr,
	)
}

// Manifest returns the manifest
func Manifest(name, namespace, configMapName, secretName, securityGroup string) *v1.Pod {
	optional := false

	return &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:    "psqlclient",
					Image:   "jbergknoff/postgresql-client",
					Command: []string{"sh"},
					EnvFrom: []v1.EnvFromSource{
						{
							ConfigMapRef: &v1.ConfigMapEnvSource{
								LocalObjectReference: v1.LocalObjectReference{
									Name: configMapName,
								},
								Optional: &optional,
							},
							SecretRef: &v1.SecretEnvSource{
								LocalObjectReference: v1.LocalObjectReference{
									Name: secretName,
								},
								Optional: &optional,
							},
						},
					},
					Stdin: true,
					TTY:   true,
				},
			},
			ImagePullSecrets: nil,
		},
	}
}
