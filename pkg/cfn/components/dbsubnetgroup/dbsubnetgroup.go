package dbsubnetgroup

import (
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/rds"
	"github.com/oslokommune/okctl/pkg/cfn"
)

type DBSubnetGroup struct {
	StoredName string
	Subnets    []cfn.Referencer
}

func (g *DBSubnetGroup) Resource() cloudformation.Resource {
	subnets := make([]string, len(g.Subnets))

	for i, s := range g.Subnets {
		subnets[i] = s.Ref()
	}

	return &rds.DBSubnetGroup{
		DBSubnetGroupDescription: g.StoredName,
		SubnetIds:                subnets,
	}
}

func (g *DBSubnetGroup) Name() string {
	return g.StoredName
}

func (g *DBSubnetGroup) Ref() string {
	return cloudformation.Ref(g.StoredName)
}

// New creates a database subnet group
//
//https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-rds-dbsubnet-group.html
func New(subnets []cfn.Referencer) *DBSubnetGroup {
	return &DBSubnetGroup{
		StoredName: "DatabaseSubnetGroup",
		Subnets:    subnets,
	}
}
