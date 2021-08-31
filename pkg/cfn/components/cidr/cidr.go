// Package cidr provides functionality for interacting with
// classless inter-domain routing
// https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing
package cidr

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"
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

	const base = 2

	size := float64(subnets * (1 << (uint64(ipv4bits) - uint64(prefixLen))))

	return uint64(math.Pow(base, math.Ceil(math.Log(size)/math.Log(base))))
}

// PrivateCidrRanges returns a set of valid private CIDRs
// https://en.wikipedia.org/wiki/Private_network#Private_IPv4_addresses
func PrivateCidrRanges() []string {
	return []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}
}

// Cidr stores the parsed range
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
		return nil, fmt.Errorf(constant.CdirNotIpv4Error, network.String())
	}

	// Ensure that the address space is large enough
	availableHosts := cidrpkg.AddressCount(network)
	if availableHosts < requiredHosts {
		return nil, fmt.Errorf(constant.CdirAddressSpaceError, network.String(), availableHosts, requiredHosts)
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
		return nil, fmt.Errorf(constant.CdirNotInLegalRangeError, network.String(), strings.Join(PrivateCidrRanges(), ", "))
	}

	return &Cidr{
		Block: network,
	}, nil
}
