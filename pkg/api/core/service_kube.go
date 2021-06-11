package core

import (
	"context"
	stderrors "errors"
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/wait"

	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api"
)

type kubeService struct {
	run api.KubeRun
}

func (k *kubeService) CreateNamespace(_ context.Context, opts api.CreateNamespaceOpts) (*api.Namespace, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating inputs", errors.Invalid)
	}

	ns, err := k.run.CreateNamespace(opts)
	if err != nil {
		return nil, errors.E(err, "creating namespace", errors.Internal)
	}

	return ns, nil
}

func (k *kubeService) ScaleDeployment(_ context.Context, opts api.ScaleDeploymentOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, "validating inputs", errors.Invalid)
	}

	err = k.run.ScaleDeployment(opts)
	if err != nil {
		kind := errors.Internal

		if strings.Contains(err.Error(), fmt.Sprintf("\"%s\" not found", opts.Name)) {
			kind = errors.NotExist
		}

		return errors.E(err, "scaling deployment", kind)
	}

	return nil
}

func (k *kubeService) CreateConfigMap(_ context.Context, opts api.CreateConfigMapOpts) (*api.ConfigMap, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating inputs", errors.Invalid)
	}

	s, err := k.run.CreateConfigMap(opts)
	if err != nil {
		return nil, errors.E(err, "creating configmap", errors.Internal)
	}

	return s, nil
}

func (k *kubeService) DeleteConfigMap(_ context.Context, opts api.DeleteConfigMapOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, "validating inputs", errors.Invalid)
	}

	err = k.run.DeleteConfigMap(opts)
	if err != nil {
		return errors.E(err, "removing configmap", errors.Internal)
	}

	return nil
}

func (k *kubeService) DeleteExternalSecrets(_ context.Context, opts api.DeleteExternalSecretsOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, "validating inputs", errors.Invalid)
	}

	err = k.run.DeleteExternalSecrets(opts)
	if err != nil {
		return errors.E(err, "removing external secrets", errors.Internal)
	}

	return nil
}

func (k *kubeService) CreateStorageClass(_ context.Context, opts api.CreateStorageClassOpts) (*api.StorageClassKube, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating inputs", errors.Invalid)
	}

	sc, err := k.run.CreateStorageClass(opts)
	if err != nil {
		return nil, errors.E(err, "creating storage class", errors.Internal)
	}

	return sc, nil
}

func (k *kubeService) DeleteNamespace(_ context.Context, opts api.DeleteNamespaceOpts) error {
	err := opts.Validate()
	if err != nil {
		return errors.E(err, "validating inputs", errors.Invalid)
	}

	err = k.run.DeleteNamespace(opts)
	if err != nil {
		return errors.E(err, "deleting namespace", errors.Internal)
	}

	return nil
}

func (k *kubeService) CreateExternalSecrets(_ context.Context, opts api.CreateExternalSecretsOpts) (*api.ExternalSecretsKube, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating inputs", errors.Invalid)
	}

	kube, err := k.run.CreateExternalSecrets(opts)
	if err != nil {
		return nil, errors.E(err, "creating external secrets", errors.Internal)
	}

	return kube, nil
}

func (k *kubeService) CreateExternalDNSKubeDeployment(_ context.Context, opts api.CreateExternalDNSKubeDeploymentOpts) (*api.ExternalDNSKube, error) {
	err := opts.Validate()
	if err != nil {
		return nil, errors.E(err, "validating inputs", errors.Invalid)
	}

	kube, err := k.run.CreateExternalDNSKubeDeployment(opts)
	if err != nil {
		var (
			wrappedErr = fmt.Errorf("creating external dns: %w", err)
			kind       errors.Kind
		)

		switch {
		case stderrors.Is(err, wait.ErrWaitTimeout):
			kind = errors.Timeout
		default:
			kind = errors.Internal
		}

		return nil, errors.E(wrappedErr, kind)
	}

	return kube, nil
}

// NewKubeService returns an initialised kube service
func NewKubeService(run api.KubeRun) api.KubeService {
	return &kubeService{
		run: run,
	}
}
