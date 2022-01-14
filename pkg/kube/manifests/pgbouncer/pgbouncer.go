// Package pgbouncer knows how to deploy PgBouncer to the
// Kubernetes cluster:
// - https://www.pgbouncer.org
package pgbouncer

import (
	"context"
	"crypto/md5" // nolint: gosec
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/watch"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"

	podpkg "github.com/oslokommune/okctl/pkg/kube/pod"

	"github.com/oslokommune/okctl/pkg/kube/forward"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

const (
	podWatchTimeoutInSeconds          = 120
	secretWatchTimeoutInSeconds       = 10
	secretWatchSleepIntervalInSeconds = 1
	maxClientConnections              = 5
)

// PgBouncer contains the state we need
// to deploy onto Kubernetes
type PgBouncer struct {
	ListenPort int32

	Pod    *v1.Pod
	Secret *v1.Secret

	In  io.Reader
	Out io.Writer
	Err io.Writer

	Client kubernetes.Interface
	Config *restclient.Config

	Ctx    context.Context
	Logger *logrus.Logger
}

// Config contains all the required inputs
type Config struct {
	Name                  string
	Database              string
	Namespace             string
	Username              string
	Password              string
	DBParamsSecretName    string
	DBParamsConfigmapName string
	Labels                map[string]string
	ListenPort            int32

	In  io.Reader
	Out io.Writer
	Err io.Writer

	ClientSet kubernetes.Interface
	Config    *restclient.Config

	Logger *logrus.Logger
}

// New returns an initialised PgBouncer client
func New(config *Config) *PgBouncer {
	secret := Secret(
		config.Name,
		config.Namespace,
		config.Username,
		config.Password,
	)

	pod := Pod(
		config.Name,
		config.Database,
		config.Namespace,
		secret.Name,
		config.DBParamsConfigmapName,
		config.DBParamsSecretName,
		config.Labels,
		config.ListenPort,
	)

	return &PgBouncer{
		ListenPort: config.ListenPort,
		Pod:        pod,
		Secret:     secret,
		In:         config.In,
		Out:        config.Out,
		Err:        config.Err,
		Client:     config.ClientSet,
		Config:     config.Config,
		Ctx:        context.Background(),
		Logger:     config.Logger,
	}
}

// DeleteSecret deletes the secret
func (p *PgBouncer) DeleteSecret(s *v1.Secret) error {
	policy := metav1.DeletePropagationForeground

	p.Logger.Info("removing pgbouncer secret")

	err := p.Client.CoreV1().Secrets(s.Namespace).Delete(p.Ctx, s.Name, metav1.DeleteOptions{
		PropagationPolicy: &policy,
	})
	if err != nil {
		return err
	}

	w, err := p.Client.CoreV1().Secrets(s.Namespace).Watch(p.Ctx, metav1.ListOptions{
		Watch:           true,
		ResourceVersion: s.ResourceVersion,
		FieldSelector:   fields.OneTermEqualSelector("metadata.name", s.Name).String(),
		LabelSelector:   labels.Everything().String(),
	})
	if err != nil {
		return err
	}

	var eventType watch.EventType

	func() {
		for {
			select {
			case event, ok := <-w.ResultChan():
				if !ok {
					return
				}

				eventType = event.Type

				if event.Type == watch.Deleted {
					w.Stop()
				}
			case <-time.After(secretWatchTimeoutInSeconds * time.Second):
				w.Stop()
			default:
				p.Logger.Info("waiting for secret to be removed")
				time.Sleep(secretWatchSleepIntervalInSeconds * time.Second)
			}
		}
	}()

	if eventType != watch.Deleted {
		return fmt.Errorf("timed out waiting for secret to be deleted")
	}

	return nil
}

// CreateSecret first removes any existing secret then creates a new one
func (p *PgBouncer) CreateSecret() error {
	secrets, err := p.Client.CoreV1().Secrets(p.Secret.Namespace).List(p.Ctx, metav1.ListOptions{})
	if err != nil {
		return fmt.Errorf("listing secrets: %w", err)
	}

	for _, s := range secrets.Items {
		if s.Name == p.Secret.Name {
			s := s

			err = p.DeleteSecret(&s)
			if err != nil {
				return fmt.Errorf("removing existing secret: %w", err)
			}
		}
	}

	p.Logger.Info("creating pgbouncer secret")

	_, err = p.Client.CoreV1().Secrets(p.Secret.Namespace).Create(p.Ctx, p.Secret, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating secret: %w", err)
	}

	return nil
}

// Create all the PgBouncer resources
func (p *PgBouncer) Create() error {
	err := p.CreateSecret()
	if err != nil {
		return err
	}

	policy := metav1.DeletePropagationForeground

	pods, err := p.Client.CoreV1().Pods(p.Pod.Namespace).List(p.Ctx, metav1.ListOptions{})
	if err != nil {
		return err
	}

	for _, po := range pods.Items {
		if po.Name == p.Pod.Name {
			po := po

			p.Logger.Info("deleting existing pgbouncer pod")

			err = p.Client.CoreV1().Pods(po.Namespace).Delete(p.Ctx, po.Name, metav1.DeleteOptions{
				PropagationPolicy: &policy,
			})
			if err != nil {
				return err
			}

			p.Logger.Info("waiting for pod to terminate")

			err = podpkg.New(p.Ctx, p.Logger, podWatchTimeoutInSeconds*time.Second, p.Client).WaitFor(podpkg.StateDeleted, &po)
			if err != nil {
				return err
			}
		}
	}

	pod, err := p.Client.CoreV1().Pods(p.Pod.Namespace).Create(p.Ctx, p.Pod, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating pod: %w", err)
	}

	err = podpkg.New(p.Ctx, p.Logger, podWatchTimeoutInSeconds*time.Second, p.Client).WaitFor(podpkg.StateRunning, pod)
	if err != nil {
		return err
	}

	return p.Forward(p.ListenPort, pod)
}

// Forward start forwarding the connection
func (p *PgBouncer) Forward(listenPort int32, pod *v1.Pod) error {
	p.Logger.Info("starting port forwarder")

	err := forward.New(os.Stdin, os.Stdout, os.Stderr, p.Config).Start(listenPort, pod)
	if err != nil {
		return err
	}

	return nil
}

// Delete removes all PgBouncer Kubernetes resources
func (p *PgBouncer) Delete() error {
	policy := metav1.DeletePropagationForeground

	p.Logger.Info("deleting pgbouncer pod")

	err := p.Client.CoreV1().Pods(p.Pod.Namespace).Delete(p.Ctx, p.Pod.Name, metav1.DeleteOptions{
		PropagationPolicy: &policy,
	})
	if err != nil {
		return fmt.Errorf("deleting pgbouncer client pod: %w", err)
	}

	err = p.DeleteSecret(p.Secret)
	if err != nil {
		return fmt.Errorf("deleting pgbouncer secret: %w", err)
	}

	return nil
}

// Secret returns the Secret for setting up the `userlist.txt`
func Secret(name, namespace, username, secret string) *v1.Secret {
	hash := md5.Sum([]byte(fmt.Sprintf("%s%s", secret, username))) // nolint: gosec

	return &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		StringData: map[string]string{
			"userlist.txt": fmt.Sprintf(
				`"%s" "md5%s"
`,
				username,
				hex.EncodeToString(hash[:]),
			),
		},
		Type: v1.SecretTypeOpaque,
	}
}

// Pod returns the Kubernetes Pod definition for PgBouncer
// - https://github.com/edoburu/docker-pgbouncer
// nolint: funlen
func Pod(
	name, database, namespace, pgBouncerSecret, dbParamsConfigMap, dbParamsSecret string,
	labels map[string]string,
	listenPort int32,
) *v1.Pod {
	// Pods using security groups must contain terminationGracePeriodSeconds in their pod spec
	// - https://docs.aws.amazon.com/eks/latest/userguide/security-groups-for-pods.html
	var terminationGracePeriodSeconds int64 = 30

	var mode int32 = 0o666

	optional := false

	db := v1.EnvVar{
		Name: "DB_NAME",
		ValueFrom: &v1.EnvVarSource{
			ConfigMapKeyRef: &v1.ConfigMapKeySelector{
				LocalObjectReference: v1.LocalObjectReference{
					Name: dbParamsConfigMap,
				},
				Key:      "PGDATABASE",
				Optional: &optional,
			},
		},
	}

	if len(database) > 0 {
		db = v1.EnvVar{
			Name:  "DB_NAME",
			Value: database,
		}
	}

	return &v1.Pod{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			Labels:    labels,
		},
		Spec: v1.PodSpec{
			TerminationGracePeriodSeconds: &terminationGracePeriodSeconds,
			Volumes: []v1.Volume{
				{
					Name: "config-volume",
					VolumeSource: v1.VolumeSource{
						Secret: &v1.SecretVolumeSource{
							SecretName:  pgBouncerSecret,
							DefaultMode: &mode,
							Optional:    &optional,
						},
					},
				},
				{
					Name: "pgbouncer",
					VolumeSource: v1.VolumeSource{
						EmptyDir: &v1.EmptyDirVolumeSource{},
					},
				},
			},
			InitContainers: []v1.Container{
				{
					Name:  "copy-ro-config",
					Image: "busybox:1.28",
					Command: []string{
						"/bin/sh", "-c", "cp /pgbouncer/userlist.txt /etc/pgbouncer/ && chown 70:70 /etc/pgbouncer/userlist.txt",
					},
					VolumeMounts: []v1.VolumeMount{
						{
							Name:      "config-volume",
							MountPath: "/pgbouncer/",
						},
						{
							Name:      "pgbouncer",
							MountPath: "/etc/pgbouncer/",
						},
					},
				},
			},
			Containers: []v1.Container{
				{
					Name:  "pgbouncer",
					Image: "edoburu/pgbouncer:1.15.0",
					VolumeMounts: []v1.VolumeMount{
						{
							Name:      "pgbouncer",
							ReadOnly:  false,
							MountPath: "/etc/pgbouncer/",
						},
					},
					Env: []v1.EnvVar{
						{
							Name:  "LISTEN_PORT",
							Value: fmt.Sprintf("%d", listenPort),
						},
						{
							// Default is 100, which is way more than needed.
							// Also, for some reason, pgbouncer has been observed to create a lot of connections,
							// which starves the actual application for connections.
							Name:  "MAX_CLIENT_CONN",
							Value: fmt.Sprintf("%d", maxClientConnections),
						},
						{
							Name: "DB_USER",
							ValueFrom: &v1.EnvVarSource{
								SecretKeyRef: &v1.SecretKeySelector{
									LocalObjectReference: v1.LocalObjectReference{
										Name: dbParamsSecret,
									},
									Key:      "PGUSER",
									Optional: &optional,
								},
							},
						},
						{
							Name: "DB_PASSWORD",
							ValueFrom: &v1.EnvVarSource{
								SecretKeyRef: &v1.SecretKeySelector{
									LocalObjectReference: v1.LocalObjectReference{
										Name: dbParamsSecret,
									},
									Key:      "PGPASSWORD",
									Optional: &optional,
								},
							},
						},
						{
							Name: "DB_HOST",
							ValueFrom: &v1.EnvVarSource{
								ConfigMapKeyRef: &v1.ConfigMapKeySelector{
									LocalObjectReference: v1.LocalObjectReference{
										Name: dbParamsConfigMap,
									},
									Key:      "PGHOST",
									Optional: &optional,
								},
							},
						},
						{
							Name: "DB_PORT",
							ValueFrom: &v1.EnvVarSource{
								ConfigMapKeyRef: &v1.ConfigMapKeySelector{
									LocalObjectReference: v1.LocalObjectReference{
										Name: dbParamsConfigMap,
									},
									Key:      "PGPORT",
									Optional: &optional,
								},
							},
						},
						db,
					},
				},
			},
		},
	}
}
