package scaffold_test

import (
	"testing"

	kaex "github.com/oslokommune/kaex/pkg/api"
	"github.com/oslokommune/okctl/pkg/scaffold"
	"github.com/stretchr/testify/assert"
)

func TestCertificateARNInjection(t *testing.T) {
	testCases := []struct {
		name string

		withApp    kaex.Application
		withCertFn scaffold.CertificateCreatorFn

		expectCertAnnotationValue string
	}{
		{
			name: "Should work",

			withApp: kaex.Application{
				Name:      "dummy-app",
				Namespace: "dummy-ns",
				Url:       "https://dummy.app",
				Port:      3000,
			},
			withCertFn: func(_ string) (string, error) {
				return "dummyarn", nil
			},

			expectCertAnnotationValue: "dummyarn",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			result, err := scaffold.NewApplicationDeployment(tc.withApp, tc.withCertFn, "", "")
			if err != nil {
				t.Error(err)
			}

			assert.Equal(t, tc.expectCertAnnotationValue, result.Ingress.Annotations["alb.ingress.kubernetes.io/certificate-arn"])
		})
	}
}
