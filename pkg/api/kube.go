package api

import "context"

type Kube struct {
}

type CreateExternalDnsKubeDeploymentOpts struct {
}

type KubeService interface {
	CreateExternalDnsKubeDeployment(ctx context.Context, opts CreateExternalDnsKubeDeploymentOpts) (*Kube, error)
}

type KubeRun interface {
	CreateExternalDnsKubeDeployment(opts CreateExternalDnsKubeDeploymentOpts) (*Kube, error)
}

type KubeStore interface {
	SaveExternalDnsKubeDeployment(kube *Kube) error
	GetExternalDnsKubeDeployment() (*Kube, error)
}
