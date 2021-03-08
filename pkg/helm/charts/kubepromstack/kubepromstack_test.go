package kubepromstack_test

import (
	"testing"

	"github.com/oslokommune/okctl/pkg/helm/charts/kubepromstack"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

func TestValues(t *testing.T) {
	testCases := []struct {
		name   string
		values *kubepromstack.Values
		golden string
	}{
		{
			name: "kubepromstack values are valid",
			values: &kubepromstack.Values{
				GrafanaCertificateARN:              "arn::1234567890/certificate/fake",
				GrafanaServiceAccountName:          "cloudwatch-something",
				GrafanaHostname:                    "grafana.okctl-test.oslo.systems",
				AuthHostname:                       "auth.okctl-test.oslo.system",
				ClientID:                           "12345dsfg456ty",
				SecretsConfigName:                  "grafana-secrets-cm",
				SecretsGrafanaCookieSecretKey:      "cookies-secret-key",
				SecretsGrafanaOauthClientSecretKey: "client-secret-key",
				SecretsGrafanaAdminUserKey:         "admin-user",
				SecretsGrafanaAdminPassKey:         "admin-pass",
			},
			golden: "values.yml",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			b, err := tc.values.RawYAML()
			assert.NoError(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.golden, b)
		})
	}
}
