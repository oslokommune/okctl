package subnet

import (
	"fmt"
	"net"
	"strings"

	cidrpkg "github.com/apparentlymart/go-cidr/cidr"
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/awslabs/goformation/v4/cloudformation/tags"
	"github.com/oslokommune/okctl/pkg/cfn"
)

const (
	TypePublic   = "public"
	TypePrivate  = "private"
	TypeDatabase = "database"
)

func Types() []string {
	return []string{
		TypePublic,
		TypePrivate,
		TypeDatabase,
	}
}

const (
	DefaultSubnets   = 9
	DefaultPrefixLen = 24
)

func AvailabilityZonesForRegion(region string) ([]string, error) {
	var azs []string

	switch region {
	case "eu-west-1":
		azs = []string{
			"eu-west-1a",
			"eu-west-1b",
			"eu-west-1c",
		}
	default:
		return nil, fmt.Errorf("no availability zone data available for region: %s", region)
	}

	return azs, nil
}

type subnet struct {
	name    string
	cluster cfn.Namer
	number  int
	network *net.IPNet
	typ     string
	az      string
	vpc     cfn.Referencer
}

func (s *subnet) Name() string {
	return s.name
}

func (s *subnet) Ref() string {
	return cloudformation.Ref(s.Name())
}

func (s *subnet) Resource() cloudformation.Resource {
	mapPublicIPonLaunch := true

	t := []tags.Tag{
		{
			Key:   fmt.Sprintf("kubernetes.io/cluster/%s", s.cluster.Name()),
			Value: "shared",
		},
	}

	switch s.typ {
	case TypePublic:
		t = append(t, tags.Tag{
			Key:   "kubernetes.io/role/elb",
			Value: "1",
		})
	case TypePrivate:
		t = append(t, tags.Tag{
			Key:   "kubernetes.io/role/internal-elb",
			Value: "1",
		})
	case TypeDatabase:
		t = []tags.Tag{}
		mapPublicIPonLaunch = false
	}

	return &ec2.Subnet{
		AvailabilityZone:    s.az,
		CidrBlock:           s.network.String(),
		MapPublicIpOnLaunch: mapPublicIPonLaunch,
		Tags:                t,
		VpcId:               s.vpc.Ref(),
	}
}

type Distributor interface {
	Next() (subnetType string, availabilityZone string, number int)
	DistinctSubnetTypes() int
	DistinctAzs() int
}

type distributor struct {
	subnetTypes []string
	subnetIndex int
	azs         []string
	azsIndex    map[string]int
}

func (d *distributor) DistinctSubnetTypes() int {
	return len(d.subnetTypes)
}

func (d *distributor) DistinctAzs() int {
	return len(d.azs)
}

func NewDistributor(subnetTypes, azs []string) (*distributor, error) {
	if len(subnetTypes) == 0 {
		return nil, fmt.Errorf("must provide at least one subnet type")
	}

	uniqueTypes := uniqueStringsInSlice(subnetTypes)

	if len(azs) == 0 {
		return nil, fmt.Errorf("must provide at least one availability zone")
	}

	uniqueAzs := uniqueStringsInSlice(azs)

	azsIndex := make(map[string]int, len(uniqueTypes))
	for _, t := range uniqueTypes {
		azsIndex[t] = 0
	}

	return &distributor{
		subnetTypes: uniqueTypes,
		azs:         uniqueAzs,
		azsIndex:    azsIndex,
	}, nil
}

func (d *distributor) Next() (string, string, int) {
	t := d.nextSubnetType()
	n := d.azsIndex[t]
	a := d.nextAz(t)

	return t, a, n
}

func (d *distributor) nextSubnetType() string {
	next := d.subnetTypes[d.subnetIndex]

	d.subnetIndex++
	if d.subnetIndex == len(d.subnetTypes) {
		d.subnetIndex = 0
	}

	return next
}

// nextAz should not be able to fail due to the checks
// introduced when creating the distributor.
func (d *distributor) nextAz(subnetType string) string {
	next := d.azs[d.azsIndex[subnetType]]

	d.azsIndex[subnetType]++
	if d.azsIndex[subnetType] == len(d.azsIndex) {
		d.azsIndex[subnetType] = 0
	}

	return next
}

type CreatorFn func(network *net.IPNet) *subnet

func NoopCreator() CreatorFn {
	return func(network *net.IPNet) *subnet {
		return &subnet{
			network: network,
		}
	}
}

func DefaultCreator(vpc cfn.Referencer, cluster cfn.Namer, dist Distributor) CreatorFn {
	return func(network *net.IPNet) *subnet {
		subnetType, az, number := dist.Next()

		return &subnet{
			cluster: cluster,
			name:    fmt.Sprintf("%sSubnet%02d", strings.Title(subnetType), number),
			number:  number,
			network: network,
			typ:     subnetType,
			az:      az,
			vpc:     vpc,
		}
	}
}

// NewSubnets creates n new subnets from the provided cidr block with the given
// network prefix size and distributes them evenly across the subnet types and
// availability zones as given by the Distributor.
func NewSubnets(n int, prefixLen int, block *net.IPNet, createFn CreatorFn) (map[string][]cfn.ResourceNameReferencer, error) {
	subnets := map[string][]cfn.ResourceNameReferencer{}

	bits, _ := block.Mask.Size()

	subnet, err := cidrpkg.Subnet(block, prefixLen-bits, 0)
	if err != nil {
		return nil, err
	}

	for i := 0; i < n; i++ {
		s := createFn(subnet)
		subnets[s.typ] = append(subnets[s.typ], s)
		subnet, _ = cidrpkg.NextSubnet(subnet, prefixLen)
	}

	return subnets, nil
}

func uniqueStringsInSlice(values []string) []string {
	var result []string

	unique := map[string]struct{}{}

	for _, val := range values {
		if _, ok := unique[val]; !ok {
			unique[val] = struct{}{}

			result = append(result, val)
		}
	}

	return result
}
