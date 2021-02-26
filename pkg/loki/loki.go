// Package loki contains helper functions for managing loki
package loki

// NewDatasourceTemplate returns an initialised Loki
// datasource consumable by Grafana
func NewDatasourceTemplate() *Datasources {
	return &Datasources{
		APIVersion: 1,
		Datasources: []Datasource{
			{
				Name:      "Loki",
				Type:      "loki",
				Access:    "proxy",
				OrgID:     1,
				URL:       "http://loki:3100",
				BasicAuth: false,
				JSONData: map[string]interface{}{
					"tlsSkipVerify": true,
				},
				Version:  1,
				Editable: false,
			},
		},
	}
}

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
	Access    string                 `json:"access"`
	OrgID     int                    `json:"orgId"`
	URL       string                 `json:"url"`
	BasicAuth bool                   `json:"basicAuth"`
	JSONData  map[string]interface{} `json:"jsonData"`
	Version   int                    `json:"version"`
	Editable  bool                   `json:"editable"`
}
