// Package subnet provides functionality for slicing and
// dicing subnets for use with EKS
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
)

const (
	// TypePublic represents a public subnet
	TypePublic = "public"
	// TypePrivate represents a private subnet
	TypePrivate = "private"
	// TypeDatabase represents a database subnet
	TypeDatabase = "database"
)

// Types returns all available subnet types
func Types() []string {
	return []string{
		TypePublic,
		TypePrivate,
		TypeDatabase,
	}
}

const (
	// DefaultSubnets defines how many subnets we will be needing by default in total
	DefaultSubnets = 9
	// DefaultPrefixLen defines the size of the subnets, e.g., how many IPs they contain
	DefaultPrefixLen = 24

	// DefaultPrivateSubnetsLogicalID defines the logical id for the stack outputs
	DefaultPrivateSubnetsLogicalID = "PrivateSubnetIds"
	// DefaultPublicSubnetsLogicalID defines the logical id for the stack outputs
	DefaultPublicSubnetsLogicalID = "PublicSubnetIds"
)

// Subnet stores state required for creating a
// cloud formation subnet
type Subnet struct {
	name    string
	cluster cfn.Namer
	number  int
	network *net.IPNet
	typ     string
	az      string
	vpc     cfn.Referencer
}

// Name returns the name of the resource
func (s *Subnet) Name() string {
	return s.name
}

// Ref returns a cloud formation intrinsic ref to the resource
func (s *Subnet) Ref() string {
	return cloudformation.Ref(s.Name())
}

// Resource returns a cloud formation resource for a subnet
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
		mapPublicIPonLaunch = false
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

// Distribution defines the interface required for something to be able
// to distribute subnets across various parameters
type Distribution interface {
	Next() (subnetType string, availabilityZone string, number int)
	DistinctSubnetTypes() int
	DistinctAzs() int
}

// Distributor stores state required for creating an even
// distribution of subnets
type Distributor struct {
	SubnetTypes []string
	SubnetIndex int
	Azs         []string
	AzsIndex    map[string]int
}

// DistinctSubnetTypes returns the number of distinct subnet types
func (d *Distributor) DistinctSubnetTypes() int {
	return len(d.SubnetTypes)
}

// DistinctAzs returns the number of distinct azs
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

// Next returns the next entry in the distribution
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

// CreatorFn defines a function that is invoked when
// creating a subnet
type CreatorFn func(network *net.IPNet) *Subnet

// NoopCreator simply returns the subnet with the given
// CIDR
func NoopCreator() CreatorFn {
	return func(network *net.IPNet) *Subnet {
		return &Subnet{
			network: network,
		}
	}
}

// DefaultCreator provides a simplified way of creating a new subnet
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

// Subnets stores the state of the created subnets
type Subnets struct {
	Public   []*Subnet
	Private  []*Subnet
	Database []*Subnet
}

// NamedOutputs returns the cloud formation outputs commonly
// required for the given subnets
func (s *Subnets) NamedOutputs() map[string]map[string]interface{} {
	private := cfn.NewJoined(DefaultPrivateSubnetsLogicalID)

	for _, p := range s.Private {
		private.Add(p.Ref())
	}

	public := cfn.NewJoined(DefaultPublicSubnetsLogicalID)

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
