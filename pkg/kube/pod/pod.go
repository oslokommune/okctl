// Package pod provides some reusable functions
// for interacting with the kubernetes pod API
package pod

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/watch"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

// State enumerates the states we can handle
type State string

const (
	// StateDeleted means we want to wait for this state
	StateDeleted State = "deleted"
	// StateRunning means we want to wait for this state
	StateRunning State = "running"

	sleepTimeWaitingForPodInSeconds = 10
)

// Pod contains the required state for interacting with
// the pod kubernetes API
type Pod struct {
	Client  kubernetes.Interface
	Ctx     context.Context
	Timeout time.Duration
	Logger  *logrus.Logger
}

// New returns an initialised Pod client
func New(ctx context.Context, logger *logrus.Logger, timeout time.Duration, client kubernetes.Interface) *Pod {
	return &Pod{
		Client:  client,
		Ctx:     ctx,
		Timeout: timeout,
		Logger:  logger,
	}
}

// WaitFor the desired state or return an error
// nolint: gocyclo gocognit
func (p *Pod) WaitFor(state State, pod *v1.Pod) error {
	status := pod.Status

	var eventType watch.EventType

	w, err := p.Client.CoreV1().Pods(pod.Namespace).Watch(p.Ctx, metav1.ListOptions{
		Watch:           true,
		ResourceVersion: pod.ResourceVersion,
		FieldSelector:   fields.OneTermEqualSelector("metadata.name", pod.Name).String(),
		LabelSelector:   labels.Everything().String(),
	})
	if err != nil {
		return fmt.Errorf("watching pod: %w", err)
	}

	func() {
		for {
			select {
			case event, ok := <-w.ResultChan():
				if !ok {
					return
				}

				pod = event.Object.(*v1.Pod)
				status = pod.Status
				eventType = event.Type

				switch event.Type {
				case watch.Deleted:
					if state == StateDeleted {
						w.Stop()
					}
				case watch.Added, watch.Modified:
					if state == StateRunning && pod.Status.Phase != v1.PodPending {
						w.Stop()
					}
				default:
					continue
				}
			case <-time.After(p.Timeout):
				w.Stop()
			default:
				what := "terminate"
				if state == StateRunning {
					what = "start"
				}

				p.Logger.Infof("waiting for pod to %s", what)
				time.Sleep(sleepTimeWaitingForPodInSeconds * time.Second)
			}
		}
	}()

	if state == StateDeleted && eventType == watch.Deleted {
		return nil
	}

	if state == StateRunning && status.Phase == v1.PodRunning {
		return nil
	}

	return fmt.Errorf("waiting for pod: %v", status.Phase)
}
