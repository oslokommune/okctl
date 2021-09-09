package core

import (
	"context"
	"fmt"
	"strings"

	merrors "github.com/mishudark/errors"

	"github.com/oslokommune/okctl/pkg/cfn"

	"github.com/oslokommune/okctl/pkg/cfn/components"

	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/jsonpatch"
	"github.com/oslokommune/okctl/pkg/kube/securitygrouppolicy/api/types/v1beta1"
	"github.com/oslokommune/okctl/pkg/scaffold"
	"github.com/oslokommune/okctl/pkg/scaffold/resources"

	"github.com/oslokommune/okctl/pkg/api"

	"github.com/oslokommune/okctl/pkg/client"
)

const (
	dnsPort      = 53
	postgresPort = 5432
)

type applicationPostgresService struct {
	manifestService  client.ApplicationManifestService
	pgService        client.ComponentService
	securityGroupAPI client.SecurityGroupAPI
	vpcService       client.VPCService
}

// AddPostgresToApplication does the required steps for allowing EKS pod to RDS traffic
func (a *applicationPostgresService) AddPostgresToApplication(ctx context.Context, opts client.AddPostgresToApplicationOpts) error {
	clusterID := clusterMetaAsID(opts.Cluster.Metadata)

	vpc, err := a.vpcService.GetVPC(ctx, clusterID)
	if err != nil {
		return fmt.Errorf("fetching VPC: %w", err)
	}

	_, err = a.pgService.GetPostgresDatabase(ctx, client.GetPostgresDatabaseOpts{
		ClusterID:    clusterID,
		DatabaseName: opts.Application.Postgres,
	})
	if err != nil {
		return fmt.Errorf("fetching database: %w", err)
	}

	securityGroup, err := a.securityGroupAPI.CreateSecurityGroup(ctx, api.CreateSecurityGroupOpts{
		ClusterID:     clusterID,
		VPCID:         vpc.VpcID,
		Name:          opts.Application.Metadata.Name,
		Description:   fmt.Sprintf("Defines network access for %s", opts.Application.Metadata.Name),
		InboundRules:  generateInboundRules(opts.Application, vpc.Cidr),
		OutboundRules: generateOutboundRules(vpc.DatabaseSubnets),
	})
	if err != nil {
		return fmt.Errorf("creating security group: %w", err)
	}

	stackName := cfn.NewStackNamer().RDSPostgres(opts.DatabaseName, opts.Cluster.Metadata.Name)

	resourceName := components.NewRDSPostgresComposer(components.RDSPostgresComposerOpts{
		ApplicationDBName: opts.DatabaseName,
		ClusterName:       opts.Cluster.Metadata.Name,
	}).CloudFormationResourceName("RDSPostgresIncoming")

	err = a.generateSecurityGroupPolicy(ctx, opts.Cluster, opts.Application, securityGroup)
	if err != nil {
		return fmt.Errorf("generating security group policy: %w", err)
	}

	_, err = a.securityGroupAPI.AddRule(ctx, api.AddRuleOpts{
		ClusterName:               opts.Cluster.Metadata.Name,
		SecurityGroupStackName:    stackName,
		SecurityGroupResourceName: resourceName,
		RuleType:                  api.RuleTypeIngress,
		Rule: api.Rule{
			Description:           fmt.Sprintf("Allow postgres traffic from %s", opts.Application.Metadata.Name),
			FromPort:              postgresPort,
			ToPort:                postgresPort,
			Protocol:              api.RuleProtocolTCP,
			SourceSecurityGroupID: securityGroup.ID,
		},
	})
	if err != nil {
		return fmt.Errorf("adding incoming rule for database security group: %w", err)
	}

	return nil
}

// RemovePostgresFromApplication cleans up required configuration for achieving communication between a EKS pod and RDS
func (a *applicationPostgresService) RemovePostgresFromApplication(ctx context.Context, opts client.RemovePostgresFromApplicationOpts) error {
	securityGroup, err := a.securityGroupAPI.GetSecurityGroup(ctx, api.GetSecurityGroupOpts{
		ClusterName: opts.Cluster.Metadata.Name,
		Name:        opts.Application.Metadata.Name,
	})
	if err != nil {
		return fmt.Errorf("getting security group: %w", err)
	}

	err = a.securityGroupAPI.DeleteSecurityGroup(ctx, api.DeleteSecurityGroupOpts{
		ClusterName: opts.Cluster.Metadata.Name,
		Name:        opts.Application.Metadata.Name,
	})
	if err != nil {
		return fmt.Errorf("deleting security group: %w", err)
	}

	stackName := cfn.NewStackNamer().RDSPostgres(opts.DatabaseName, opts.Cluster.Metadata.Name)

	resourceName := components.NewRDSPostgresComposer(components.RDSPostgresComposerOpts{
		ApplicationDBName: opts.Application.Postgres,
		ClusterName:       opts.Cluster.Metadata.Name,
	}).CloudFormationResourceName("RDSPostgresIncoming")

	err = a.removeSecurityGroupPolicy(ctx, opts.Cluster, opts.Application, securityGroup)
	if err != nil {
		return fmt.Errorf("removing security group policy: %w", err)
	}

	err = a.securityGroupAPI.RemoveRule(ctx, api.RemoveRuleOpts{
		ClusterName:               opts.Cluster.Metadata.Name,
		SecurityGroupStackName:    stackName,
		SecurityGroupResourceName: resourceName,
		RuleType:                  api.RuleTypeIngress,
		Rule: api.Rule{
			FromPort:              postgresPort,
			ToPort:                postgresPort,
			Protocol:              api.RuleProtocolTCP,
			SourceSecurityGroupID: securityGroup.ID,
		},
	})
	if err != nil {
		return fmt.Errorf("removing rule from inbound database security group: %w", err)
	}

	return nil
}

// HasPostgresIntegration knows if an application has an existing integration with a database
func (a *applicationPostgresService) HasPostgresIntegration(ctx context.Context, opts client.HasPostgresIntegrationOpts) (bool, error) {
	sg, err := a.securityGroupAPI.GetSecurityGroup(ctx, api.GetSecurityGroupOpts{
		Name:        opts.Application.Metadata.Name,
		ClusterName: opts.Cluster.Metadata.Name,
	})
	if err != nil {
		if merrors.IsKind(err, merrors.NotExist) {
			return false, nil
		}

		return false, fmt.Errorf("acquiring security group: %w", err)
	}

	patch, err := a.manifestService.GetPatch(ctx, client.GetPatchOpts{
		ApplicationName: opts.Application.Metadata.Name,
		ClusterName:     opts.Cluster.Metadata.Name,
		Kind:            v1beta1.SecurityGroupPolicyKind,
	})
	if err != nil {
		if merrors.IsKind(err, merrors.NotExist) {
			return false, nil
		}

		return false, fmt.Errorf("acquiring patch: %w", err)
	}

	return patch.HasOperation(jsonpatch.Operation{
		Type:  jsonpatch.OperationTypeAdd,
		Path:  "/spec/securityGroups/groupIds/0",
		Value: sg.ID,
	}), nil
}

func (a *applicationPostgresService) generateSecurityGroupPolicy(ctx context.Context, cluster v1alpha1.Cluster, app v1alpha1.Application, sg api.SecurityGroup) error {
	baseManifest, err := scaffold.ResourceAsBytes(resources.CreateSecurityGroupPolicy(app))
	if err != nil {
		return fmt.Errorf("converting security group policy to bytes: %w", err)
	}

	err = a.manifestService.SaveManifest(ctx, client.SaveManifestOpts{
		ApplicationName: app.Metadata.Name,
		Filename:        "security-group-policy.yaml",
		Content:         baseManifest,
	})
	if err != nil {
		return fmt.Errorf("saving security group base manifest: %w", err)
	}

	err = a.manifestService.SavePatch(ctx, client.SavePatchOpts{
		ApplicationName: app.Metadata.Name,
		ClusterName:     cluster.Metadata.Name,
		Kind:            v1beta1.SecurityGroupPolicyKind,
		Patch: jsonpatch.Patch{
			Operations: []jsonpatch.Operation{
				{
					Type:  jsonpatch.OperationTypeAdd,
					Path:  "/spec/securityGroups/groupIds/0",
					Value: sg.ID,
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("saving security group patch: %w", err)
	}

	return nil
}

func (a *applicationPostgresService) removeSecurityGroupPolicy(
	ctx context.Context,
	cluster v1alpha1.Cluster,
	application v1alpha1.Application,
	securityGroup api.SecurityGroup,
) error {
	patch, err := a.manifestService.GetPatch(ctx, client.GetPatchOpts{
		ApplicationName: application.Metadata.Name,
		ClusterName:     cluster.Metadata.Name,
		Kind:            v1beta1.SecurityGroupPolicyKind,
	})
	if err != nil {
		return fmt.Errorf("could not find existing security group patch: %w", err)
	}

	operationIndex := getOperationIndexForSecurityGroupID(patch.Operations, securityGroup.ID)
	if operationIndex == -1 {
		return nil
	}

	patch.Operations[operationIndex] = patch.Operations[len(patch.Operations)-1]
	patch.Operations = patch.Operations[:len(patch.Operations)-1]

	err = a.manifestService.SavePatch(ctx, client.SavePatchOpts{
		ApplicationName: application.Metadata.Name,
		ClusterName:     cluster.Metadata.Name,
		Kind:            v1beta1.SecurityGroupPolicyKind,
		Patch:           patch,
	})
	if err != nil {
		return fmt.Errorf("saving updated security group policy patch: %w", err)
	}

	return nil
}

func generateInboundRules(app v1alpha1.Application, vpcCidr string) []api.Rule {
	rules := []api.Rule{
		{
			Description: "Required DNS/tcp entrypoint for control plane",
			FromPort:    dnsPort,
			ToPort:      dnsPort,
			CidrIP:      vpcCidr,
			Protocol:    api.RuleProtocolTCP,
		},
		{
			Description: "Required DNS/udp entrypoint for control plane",
			FromPort:    dnsPort,
			ToPort:      dnsPort,
			CidrIP:      vpcCidr,
			Protocol:    api.RuleProtocolUDP,
		},
	}

	if app.HasIngress() {
		rules = append(rules, api.Rule{
			Description: "Inbound traffic to the application",
			FromPort:    int(app.Port),
			ToPort:      int(app.Port),
			CidrIP:      vpcCidr,
			Protocol:    api.RuleProtocolTCP,
		})
	}

	return rules
}

func generateOutboundRules(dbSubnets []client.VpcSubnet) []api.Rule {
	rules := make([]api.Rule, len(dbSubnets))

	for index, subnet := range dbSubnets {
		rules[index] = api.Rule{
			Description: "Allow Postgres traffic to database subnet",
			FromPort:    postgresPort,
			ToPort:      postgresPort,
			CidrIP:      subnet.Cidr,
			Protocol:    api.RuleProtocolTCP,
		}
	}

	rules = append(rules,
		api.Rule{
			Description: "Allow all outbound IPv4 traffic.",
			CidrIP:      "0.0.0.0/0",
			Protocol:    api.RuleProtocolAll,
		},
	)

	return rules
}

func clusterMetaAsID(meta v1alpha1.ClusterMeta) api.ID {
	return api.ID{
		Region:       meta.Region,
		AWSAccountID: meta.AccountID,
		ClusterName:  meta.Name,
	}
}

func getOperationIndexForSecurityGroupID(operations []jsonpatch.Operation, securityGroupID string) int {
	for index, operation := range operations {
		value, ok := operation.Value.(string)

		if !ok {
			continue
		}

		if strings.HasPrefix(value, securityGroupID) {
			return index
		}
	}

	return -1
}

// NewApplicationPostgresService initializes a ApplicationPostgresService
func NewApplicationPostgresService(
	manifestService client.ApplicationManifestService,
	pgService client.ComponentService,
	securityGroupAPI client.SecurityGroupAPI,
	vpcService client.VPCService,
) client.ApplicationPostgresService {
	return &applicationPostgresService{
		manifestService:  manifestService,
		pgService:        pgService,
		securityGroupAPI: securityGroupAPI,
		vpcService:       vpcService,
	}
}
