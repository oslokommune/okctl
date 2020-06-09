package vpc

import (
	"github.com/oslokommune/okctl/pkg/cfn"
	cidrPkg "github.com/oslokommune/okctl/pkg/cfn/components/cidr"
	clusterPkg "github.com/oslokommune/okctl/pkg/cfn/components/cluster"
	"github.com/oslokommune/okctl/pkg/cfn/components/dbsubnetgroup"
	"github.com/oslokommune/okctl/pkg/cfn/components/eip"
	"github.com/oslokommune/okctl/pkg/cfn/components/internetgateway"
	"github.com/oslokommune/okctl/pkg/cfn/components/natgateway"
	"github.com/oslokommune/okctl/pkg/cfn/components/route"
	"github.com/oslokommune/okctl/pkg/cfn/components/routetable"
	"github.com/oslokommune/okctl/pkg/cfn/components/routetableassociation"
	"github.com/oslokommune/okctl/pkg/cfn/components/securitygroup"
	"github.com/oslokommune/okctl/pkg/cfn/components/subnet"
	vpcPkg "github.com/oslokommune/okctl/pkg/cfn/components/vpc"
	"github.com/oslokommune/okctl/pkg/cfn/components/vpcgatewayattachment"
)

type Builder struct {
	Name      string
	Env       string
	CidrBlock string
	Region    string
}

func New(name, env, cidrBlock, region string) *Builder {
	return &Builder{
		Name:      name,
		Env:       env,
		CidrBlock: cidrBlock,
		Region:    region,
	}
}

func (b *Builder) Build() ([]cfn.ResourceNamer, error) {
	resources := []cfn.ResourceNamer{}

	cluster := clusterPkg.New(b.Name, b.Env)

	cidr, err := cidrPkg.NewDefault(b.CidrBlock)
	if err != nil {
		return nil, err
	}

	vpc := vpcPkg.New(cluster, cidr.Block)
	igw := internetgateway.New()
	gwa := vpcgatewayattachment.New(vpc, igw)
	prt := routetable.NewPublic(vpc)
	pr := route.NewPublic(gwa, prt, igw)

	resources = append(resources, vpc, igw, gwa, prt, pr)

	subnets, err := subnet.NewDefault(cidr.Block, b.Region, vpc, cluster)
	if err != nil {
		return nil, err
	}

	for i, sub := range subnets.Private {
		rt := routetable.NewPrivate(i, vpc)
		e := eip.New(i, gwa)
		ngw := natgateway.New(i, gwa, e, sub)
		r := route.NewPrivate(i, gwa, rt, ngw)
		assoc := routetableassociation.NewPrivate(i, sub, rt)

		resources = append(resources, sub, rt, e, ngw, r, assoc)
	}

	for i, sub := range subnets.Public {
		assoc := routetableassociation.NewPublic(i, sub, prt)

		resources = append(resources, sub, assoc)
	}

	dbSubnets := make([]cfn.Referencer, len(subnets.Database))

	for i, sub := range subnets.Database {
		dbSubnets[i] = sub

		resources = append(resources, sub)
	}

	dsg := dbsubnetgroup.New(dbSubnets)
	cpsg := securitygroup.ControlPlane(vpc)

	resources = append(resources, dsg, cpsg)

	return resources, nil
}
