package dbsubnetgroup

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/rds"
	"github.com/oslokommune/okctl/pkg/cfn"
)

type dbSubnetGroup struct {
	name    string
	subnets []cfn.Referencer
}

func (g *dbSubnetGroup) Resource() cloudformation.Resource {
	subnets := make([]string, len(g.subnets))

	for i, s := range g.subnets {
		subnets[i] = s.Ref()
	}

	return &rds.DBSubnetGroup{
		DBSubnetGroupDescription: g.name,
		SubnetIds:                subnets,
	}
}

func (g *dbSubnetGroup) Name() string {
	return g.name
}

func (g *dbSubnetGroup) Ref() string {
	return cloudformation.Ref(g.name)
}

func New(subnets []cfn.Referencer) *dbSubnetGroup {
	return &dbSubnetGroup{
		name:    "DatabaseSubnetGroup",
		subnets: subnets,
	}
}
