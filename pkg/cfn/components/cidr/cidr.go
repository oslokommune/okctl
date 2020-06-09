package cidr

import (
	"fmt"
	"math"
	"net"
	"strings"

	cidrpkg "github.com/apparentlymart/go-cidr/cidr"
	"github.com/oslokommune/okctl/pkg/cfn/components/subnet"
)

// RequiredHosts calculates the required number of hosts
// and rounds up to the nearest power of two
func RequiredHosts(subnets, prefixLen int) uint64 {
	const ipv4bits = 32
	size := float64(subnets * (1 << (uint64(ipv4bits) - uint64(prefixLen))))

	// nolint
	return uint64(math.Pow(2, math.Ceil(math.Log(size)/math.Log(2))))
}

func PrivateCidrRanges() []string {
	return []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}
}

type Cidr struct {
	Block *net.IPNet
}

// NewDefault parses the provided Cidr given the default restrictions
func NewDefault(block string) (*Cidr, error) {
	return New(block, RequiredHosts(subnet.DefaultSubnets, subnet.DefaultPrefixLen), PrivateCidrRanges())
}

// New parses a provided cidr and ensures that it
// is valid
func New(from string, requiredHosts uint64, validRanges []string) (*Cidr, error) {
	// Ensure that it is a valid cidr
	addr, network, err := net.ParseCIDR(from)
	if err != nil {
		return nil, err
	}

	// Ensure that it is a ipv4 CIDR
	v4addr := addr.To4()
	if v4addr == nil {
		return nil, fmt.Errorf("cidr (%s) is not of type IPv4", network.String())
	}

	// Ensure that the address space is large enough
	availableHosts := cidrpkg.AddressCount(network)
	if availableHosts < requiredHosts {
		return nil, fmt.Errorf("address space of cidr (%s) is less than required: %d < %d", network.String(), availableHosts, requiredHosts)
	}

	// Ensure that the provided cidr falls within a valid range
	var isValid bool

	for _, v := range validRanges {
		_, outer, err := net.ParseCIDR(v)
		if err != nil {
			return nil, err
		}

		outerOnes, _ := outer.Mask.Size()
		innerOnes, _ := network.Mask.Size()

		first, last := cidrpkg.AddressRange(network)

		if outerOnes <= innerOnes && outer.Contains(first) && outer.Contains(last) {
			isValid = true
			break
		}
	}

	if !isValid {
		return nil, fmt.Errorf("provided cidr (%s) is not in the legal ranges: %s", network.String(), strings.Join(PrivateCidrRanges(), ", "))
	}

	return &Cidr{
		Block: network,
	}, nil
}
