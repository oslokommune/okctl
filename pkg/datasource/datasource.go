// Package datasource contains helper functions for creating grafana datasource
// definitions
package datasource

import "github.com/aws/aws-sdk-go/aws"

// Datasources contains a set of datasources
type Datasources struct {
	APIVersion  int          `json:"apiVersion"`
	Datasources []Datasource `json:"datasources"`
}

// Datasource contains a single datasource
// nolint: maligned
type Datasource struct {
	Name      string                 `json:"name"`
	Type      string                 `json:"type"`
	Access    *string                `json:"access,omitempty"`
	OrgID     *int                   `json:"orgId,omitempty"`
	URL       *string                `json:"url,omitempty"`
	BasicAuth *bool                  `json:"basicAuth,omitempty"`
	JSONData  map[string]interface{} `json:"jsonData"`
	Version   *int                   `json:"version,omitempty"`
	Editable  *bool                  `json:"editable,omitempty"`
}

// NewCloudWatch returns an initialised CloudWatch
// datasource consumable by Grafana
func NewCloudWatch(region string) *Datasources {
	return &Datasources{
		APIVersion: 1,
		Datasources: []Datasource{
			{
				Name: "CloudWatch",
				Type: "cloudwatch",
				JSONData: map[string]interface{}{
					"authType":      "default",
					"defaultRegion": region,
				},
			},
		},
	}
}

// NewLoki returns an initialised Loki
// datasource consumable by Grafana
func NewLoki() *Datasources {
	return &Datasources{
		APIVersion: 1,
		Datasources: []Datasource{
			{
				Name:      "Loki",
				Type:      "loki",
				Access:    aws.String("proxy"),
				OrgID:     aws.Int(1),
				URL:       aws.String("http://loki:3100"),
				BasicAuth: aws.Bool(false),
				JSONData: map[string]interface{}{
					"tlsSkipVerify": true,
				},
				Version:  aws.Int(1),
				Editable: aws.Bool(false),
			},
		},
	}
}

// NewTempo returns an initialised tempo
// datasource consumable by Grafana
func NewTempo() *Datasources {
	return &Datasources{
		APIVersion: 1,
		Datasources: []Datasource{
			{
				Name:      "Tempo",
				Type:      "tempo",
				Access:    aws.String("proxy"),
				OrgID:     aws.Int(1),
				URL:       aws.String("http://tempo:16686"),
				BasicAuth: aws.Bool(false),
				JSONData: map[string]interface{}{
					"tlsSkipVerify": true,
				},
				Version:  aws.Int(1),
				Editable: aws.Bool(false),
			},
		},
	}
}
