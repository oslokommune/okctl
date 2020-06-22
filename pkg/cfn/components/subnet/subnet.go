package subnet

import (
	"fmt"
	"net"
	"strings"

	cidrPkg "github.com/apparentlymart/go-cidr/cidr"
	"github.com/awslabs/goformation/v4/cloudformation"
	"github.com/awslabs/goformation/v4/cloudformation/ec2"
	"github.com/awslabs/goformation/v4/cloudformation/tags"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/cfn"
	"github.com/oslokommune/okctl/pkg/cfn/builder/output"
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

type Subnet struct {
	name    string
	cluster cfn.Namer
	number  int
	network *net.IPNet
	typ     string
	az      string
	vpc     cfn.Referencer
}

func (s *Subnet) Name() string {
	return s.name
}

func (s *Subnet) Ref() string {
	return cloudformation.Ref(s.Name())
}

func (s *Subnet) Resource() cloudformation.Resource {
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

type Distribution interface {
	Next() (subnetType string, availabilityZone string, number int)
	DistinctSubnetTypes() int
	DistinctAzs() int
}

type Distributor struct {
	SubnetTypes []string
	SubnetIndex int
	Azs         []string
	AzsIndex    map[string]int
}

func (d *Distributor) DistinctSubnetTypes() int {
	return len(d.SubnetTypes)
}

func (d *Distributor) DistinctAzs() int {
	return len(d.Azs)
}

// NewDistributor returns a Distributor that can be used
// to delegate a set of subnets evenly across the provided
// types and availability zones
func NewDistributor(subnetTypes, azs []string) (*Distributor, error) {
	if len(subnetTypes) == 0 {
		return nil, fmt.Errorf("must provide at least one Subnet type")
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

	return &Distributor{
		SubnetTypes: uniqueTypes,
		Azs:         uniqueAzs,
		AzsIndex:    azsIndex,
	}, nil
}

func (d *Distributor) Next() (string, string, int) {
	t := d.nextSubnetType()
	n := d.AzsIndex[t]
	a := d.nextAz(t)

	return t, a, n
}

func (d *Distributor) nextSubnetType() string {
	next := d.SubnetTypes[d.SubnetIndex]

	d.SubnetIndex++
	if d.SubnetIndex == len(d.SubnetTypes) {
		d.SubnetIndex = 0
	}

	return next
}

// nextAz should not be able to fail due to the checks
// introduced when creating the Distributor.
func (d *Distributor) nextAz(subnetType string) string {
	next := d.Azs[d.AzsIndex[subnetType]]

	d.AzsIndex[subnetType]++
	if d.AzsIndex[subnetType] == len(d.AzsIndex) {
		d.AzsIndex[subnetType] = 0
	}

	return next
}

type CreatorFn func(network *net.IPNet) *Subnet

func NoopCreator() CreatorFn {
	return func(network *net.IPNet) *Subnet {
		return &Subnet{
			network: network,
		}
	}
}

func DefaultCreator(vpc cfn.Referencer, cluster cfn.Namer, dist Distribution) CreatorFn {
	return func(network *net.IPNet) *Subnet {
		subnetType, az, number := dist.Next()

		return &Subnet{
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

type Subnets struct {
	Public   []*Subnet
	Private  []*Subnet
	Database []*Subnet
}

func (s *Subnets) NamedOutputs() map[string]map[string]interface{} {
	private := output.NewJoined("PrivateSubnetIds")

	for _, p := range s.Private {
		private.Add(p.Ref())
	}

	public := output.NewJoined("PublicSubnetIds")

	for _, p := range s.Public {
		public.Add(p.Ref())
	}

	return map[string]map[string]interface{}{
		private.Name(): private.Outputs(),
		public.Name():  public.Outputs(),
	}
}

// NewDefault creates a default Subnet distribution for the given network and region
func NewDefault(network *net.IPNet, region string, vpc cfn.Referencer, cluster cfn.Namer) (*Subnets, error) {
	azs, err := v1alpha1.SupportedAvailabilityZones(region)
	if err != nil {
		return nil, err
	}

	dist, err := NewDistributor(Types(), azs)
	if err != nil {
		return nil, err
	}

	return New(DefaultSubnets, DefaultPrefixLen, network, DefaultCreator(vpc, cluster, dist))
}

// New creates n new subnets from the provided cidr block with the given
// network prefix size and distributes them evenly across the Subnet types and
// availability zones as given by the Distribution.
func New(n int, prefixLen int, block *net.IPNet, createFn CreatorFn) (*Subnets, error) {
	subnets := &Subnets{}

	bits, _ := block.Mask.Size()

	subnet, err := cidrPkg.Subnet(block, prefixLen-bits, 0)
	if err != nil {
		return nil, err
	}

	for i := 0; i < n; i++ {
		s := createFn(subnet)

		switch s.typ {
		case TypePublic:
			subnets.Public = append(subnets.Public, s)
		case TypePrivate:
			subnets.Private = append(subnets.Private, s)
		case TypeDatabase:
			subnets.Database = append(subnets.Database, s)
		}

		subnet, _ = cidrPkg.NextSubnet(subnet, prefixLen)
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
