package dryrun

import (
	"context"
	"fmt"
	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/apis/eksctl.io/v1alpha5"
	"github.com/oslokommune/okctl/pkg/client"
	"io"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type clusterService struct {
	out io.Writer
}

func (c clusterService) CreateCluster(_ context.Context, opts client.ClusterCreateOpts) (*client.Cluster, error) {
	fmt.Fprintf(c.out, formatCreate("Kubernetes cluster"))

	return &client.Cluster{
		ID:   opts.ID,
		Name: "",
		Config: &v1alpha5.ClusterConfig{
			TypeMeta:        metav1.TypeMeta{},
			Metadata:        v1alpha5.ClusterMeta{},
			IAM:             v1alpha5.ClusterIAM{},
			VPC:             nil,
			FargateProfiles: nil,
			NodeGroups:      nil,
			Status:          nil,
			CloudWatch:      nil,
			Addons:          nil,
		},
	}, nil
}

func (c clusterService) DeleteCluster(_ context.Context, _ client.ClusterDeleteOpts) error {
	fmt.Fprintf(c.out, formatDelete("Kubernetes cluster"))

	return nil
}

func (c clusterService) GetClusterSecurityGroupID(_ context.Context, _ client.GetClusterSecurityGroupIDOpts) (*api.ClusterSecurityGroupID, error) {
	panic("implement me")
}
