package defaults

import (
	"fmt"
	"net"

	"github.com/oslokommune/okctl/pkg/cfn"
	cidrPkg "github.com/oslokommune/okctl/pkg/cfn/cidr"
	clusterPkg "github.com/oslokommune/okctl/pkg/cfn/cluster"
	"github.com/oslokommune/okctl/pkg/cfn/dbsubnetgroup"
	"github.com/oslokommune/okctl/pkg/cfn/eip"
	"github.com/oslokommune/okctl/pkg/cfn/internetgateway"
	"github.com/oslokommune/okctl/pkg/cfn/natgateway"
	"github.com/oslokommune/okctl/pkg/cfn/route"
	"github.com/oslokommune/okctl/pkg/cfn/routetable"
	"github.com/oslokommune/okctl/pkg/cfn/routetableassociation"
	"github.com/oslokommune/okctl/pkg/cfn/securitygroup"
	"github.com/oslokommune/okctl/pkg/cfn/subnet"
	vpcPkg "github.com/oslokommune/okctl/pkg/cfn/vpc"
	"github.com/oslokommune/okctl/pkg/cfn/vpcgatewayattachment"
)

// nolint
func subnets(region string, vpc cfn.Referencer, cluster cfn.Namer, network *net.IPNet) ([]cfn.ResourceNameReferencer, []cfn.ResourceNameReferencer, []cfn.ResourceNameReferencer, error) {
	azs, err := subnet.AvailabilityZonesForRegion(region)
	if err != nil {
		return nil, nil, nil, err
	}

	dist, err := subnet.NewDistributor(subnet.Types(), azs)
	if err != nil {
		return nil, nil, nil, err
	}

	subnets, err := subnet.NewSubnets(
		subnet.DefaultSubnets,
		subnet.DefaultPrefixLen,
		network,
		subnet.DefaultCreator(vpc, cluster, dist),
	)
	if err != nil {
		return nil, nil, nil, err
	}

	privateSubnets, ok := subnets[subnet.TypePrivate]
	if !ok {
		return nil, nil, nil, fmt.Errorf("failed to retrieve the generated private subnets")
	}

	publicSubnets, ok := subnets[subnet.TypePublic]
	if !ok {
		return nil, nil, nil, fmt.Errorf("failed to retrieve the generatetd public subnets")
	}

	databaseSubnets, ok := subnets[subnet.TypeDatabase]
	if !ok {
		return nil, nil, nil, fmt.Errorf("failed to retrieve the generated database subnets")
	}

	return publicSubnets, privateSubnets, databaseSubnets, nil
}

// nolint
func VPC(name, env, cidrBlock, region string) ([]byte, error) {
	m := cfn.New()

	cluster := clusterPkg.New(name, env)

	cidr, err := cidrPkg.New(
		cidrBlock,
		cidrPkg.RequiredHosts(subnet.DefaultSubnets, subnet.DefaultPrefixLen),
		cidrPkg.PrivateCidrRanges(),
	)
	if err != nil {
		return nil, err
	}

	vpc, err := vpcPkg.NewVPC(cluster, cidr.Block)
	if err != nil {
		return nil, err
	}

	err = m.Add(vpc)
	if err != nil {
		return nil, err
	}

	igw := internetgateway.New()

	err = m.Add(igw)
	if err != nil {
		return nil, err
	}

	gwa := vpcgatewayattachment.New(vpc, igw)

	err = m.Add(gwa)
	if err != nil {
		return nil, err
	}

	publicRouteTable := routetable.NewPublic(vpc)
	publicRoute := route.NewPublic(gwa, publicRouteTable, igw)

	err = m.Add(publicRouteTable, publicRoute)
	if err != nil {
		return nil, err
	}

	publicSubnets, privateSubnets, databaseSubnets, err := subnets(region, vpc, cluster, cidr.Block)
	if err != nil {
		return nil, err
	}

	for i, sub := range privateSubnets {
		rt := routetable.NewPrivate(i, vpc)
		e := eip.New(i, gwa)
		ngw := natgateway.New(i, gwa, e, sub)
		r := route.NewPrivate(i, gwa, rt, ngw)
		assoc := routetableassociation.NewPrivate(i, sub, rt)

		err = m.Add(sub, rt, e, ngw, r, assoc)
		if err != nil {
			return nil, err
		}
	}

	for i, sub := range publicSubnets {
		assoc := routetableassociation.NewPublic(i, sub, publicRouteTable)

		err = m.Add(sub, assoc)
		if err != nil {
			return nil, err
		}
	}

	dbSubnets := make([]cfn.Referencer, len(databaseSubnets))

	for i, sub := range databaseSubnets {
		dbSubnets[i] = sub

		err = m.Add(sub)
		if err != nil {
			return nil, err
		}
	}

	dbSubnetGroup := dbsubnetgroup.New(dbSubnets)

	err = m.Add(dbSubnetGroup)
	if err != nil {
		return nil, err
	}

	controlPlane := securitygroup.ControlPlane(vpc)

	err = m.Add(controlPlane)
	if err != nil {
		return nil, err
	}

	return m.YAML()
}
