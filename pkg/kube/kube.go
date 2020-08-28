// Package kube implements a kubernetes client
// Parts of this file have been stolen from:
// - https://github.com/helm/helm/blob/master/pkg/kube/wait.go
package kube

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	deploymentUtil "github.com/oslokommune/okctl/internal/third_party/k8s.io/kubernetes/deployment/util"
	"github.com/sirupsen/logrus"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/api/rbac/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Kuber provides the methods that are available
// by a concrete implementation
type Kuber interface {
	Apply(fn ApplyFn, fns ...ApplyFn)
}

// Kube contains state for communicating with
// a kubernetes cluster
type Kube struct {
	KubeConfigPath string
	ClientSet      *kubernetes.Clientset
	Ctx            context.Context
	Log            *logrus.Logger
}

// ApplyFn defines the signature of a function that applies
// some operation to the kubernetes cluster
type ApplyFn func(clientSet kubernetes.Interface) (interface{}, error)

// New returns an initialised kubernetes client
func New(kubeConfigPath string) (*Kube, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		return nil, err
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.Out = ioutil.Discard

	return &Kube{
		KubeConfigPath: kubeConfigPath,
		ClientSet:      clientSet,
		Ctx:            context.Background(),
		Log:            logger,
	}, nil
}

// WithLogger sets a logger
func (k *Kube) WithLogger(log *logrus.Logger) *Kube {
	k.Log = log

	return k
}

// Apply all the functions to the cluster
func (k *Kube) Apply(first ApplyFn, rest ...ApplyFn) ([]interface{}, error) {
	fns := append([]ApplyFn{first}, rest...)
	values := make([]interface{}, len(fns))

	for i, fn := range fns {
		v, err := fn(k.ClientSet)
		if err != nil {
			return nil, err
		}

		values[i] = v
	}

	return values, nil
}

// Watch a set of resources
func (k *Kube) Watch(resources []interface{}, timeout time.Duration) error {
	var err error

	// Move the wait to this point if we starting getting many resources
	// that we can check the health of, like helm do.
	for _, resource := range resources {
		switch r := resource.(type) {
		case *appsv1.Deployment:
			err = k.WatchDeployment(r, timeout)
			if err != nil {
				return err
			}
		case *v1beta1.ClusterRole, *v1beta1.ClusterRoleBinding:
			continue
		default:
			return fmt.Errorf("unknown resource type: %s", resource)
		}
	}

	return nil
}

// WatchDeployment and wait until it is ready or we hit the timeout
// most of this code, together with deploymentReady is taken from:
// - https://github.com/helm/helm
//
// Luckily we have control over the types of resources we want to watch, so we can
// simplify a little.
func (k *Kube) WatchDeployment(deployment *appsv1.Deployment, timeout time.Duration) error {
	return wait.Poll(2*time.Second, timeout, func() (bool, error) { // nolint: gomnd
		currentDeployment, err := k.ClientSet.AppsV1().Deployments(deployment.Namespace).Get(context.Background(), deployment.Name, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		newReplicaSet, err := deploymentUtil.GetNewReplicaSet(currentDeployment, k.ClientSet.AppsV1())
		if err != nil || newReplicaSet == nil {
			return false, err
		}

		if !k.deploymentReady(newReplicaSet, currentDeployment) {
			return false, nil
		}

		return true, nil
	})
}

func (k *Kube) deploymentReady(rs *appsv1.ReplicaSet, dep *appsv1.Deployment) bool {
	expectedReady := *dep.Spec.Replicas - deploymentUtil.MaxUnavailable(*dep)
	if !(rs.Status.ReadyReplicas >= expectedReady) {
		k.Log.Infof("Deployment is not ready: %s/%s. %d out of %d expected pods are ready", dep.Namespace, dep.Name, rs.Status.ReadyReplicas, expectedReady)
		return false
	}

	return true
}

// Debug a namespace
func (k *Kube) Debug(namespace string) (map[string][]string, error) {
	events, err := k.ClientSet.CoreV1().Events(namespace).List(k.Ctx, metav1.ListOptions{})
	if err != nil {
		return nil, nil
	}

	eventStrings := make([]string, len(events.Items))

	for i, event := range events.Items {
		j, err := json.MarshalIndent(event, "", "  ")
		if err != nil {
			return nil, err
		}

		eventStrings[i] = string(j)
	}

	pods, err := k.ClientSet.CoreV1().Pods(namespace).List(k.Ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	podLogs := make([]string, len(pods.Items))
	podSpec := make([]string, len(pods.Items))

	for i, pod := range pods.Items {
		request := k.ClientSet.CoreV1().Pods(namespace).GetLogs(pod.Name, &v1.PodLogOptions{})

		raw, err := request.DoRaw(k.Ctx)
		if err != nil {
			return nil, err
		}

		podLogs[i] = string(raw)

		j, err := json.MarshalIndent(pod, "", "  ")
		if err != nil {
			return nil, err
		}

		podSpec[i] = string(j)
	}

	return map[string][]string{
		"podLogs":  podLogs,
		"podSpecs": podSpec,
		"events":   eventStrings,
	}, nil
}
