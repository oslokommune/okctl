package route53

import (
	"strings"

	route "github.com/aws/aws-sdk-go/service/route53"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

// Route53er defines the available methods
type Route53er interface {
	HostedZones() ([]HostedZone, error)
}

// HostedZone contains information about an AWS Route53 HostedZone
type HostedZone struct {
	ID     string
	Domain string
	FQDN   string
	Public bool
}

// Route53 contains the state required for querying
// the Route53 API
type Route53 struct {
	provider v1alpha1.CloudProvider
}

// New returns an initialised route53 client
func New(provider v1alpha1.CloudProvider) *Route53 {
	return &Route53{
		provider: provider,
	}
}

// PublicHostedZones returns all the registered hosted zones
func (r *Route53) PublicHostedZones() ([]*HostedZone, error) {
	out, err := r.provider.Route53().ListHostedZones(&route.ListHostedZonesInput{})
	if err != nil {
		return nil, err
	}

	var zones []*HostedZone

	for _, hz := range out.HostedZones {
		if !*hz.Config.PrivateZone {
			if !*hz.Config.PrivateZone {
				zones = append(zones, &HostedZone{
					ID:     strings.TrimPrefix(*hz.Id, "/hostedzone/"),
					Domain: strings.TrimSuffix(*hz.Name, "."),
					FQDN:   *hz.Name,
					Public: true,
				})
			}
		}
	}

	return zones, nil
}
