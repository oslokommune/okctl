package cognito

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"github.com/stretchr/testify/assert"
)

func TestGetRelevantUserPoolClient(t *testing.T) {
	testCases := []struct {
		name         string
		withClients  []string
		expectClient string
	}{
		{
			name: "Should return the argocd client in an expected okctl setup",
			withClients: []string{
				"okctl-mock-cluster-argocd",
				"okctl-mock-cluster-grafana",
			},
			expectClient: "okctl-mock-cluster-argocd",
		},
		{
			name: "Should return the argocd client in a strange setup with a gazzilion clients",
			withClients: []string{
				"okctl-mock-cluster-client1",
				"okctl-mock-cluster-grafana",
				"okctl-mock-cluster-client2",
				"okctl-mock-cluster-client3",
				"okctl-mock-cluster-argocd",
				"okctl-mock-cluster-client4",
				"okctl-mock-cluster-client5",
			},
			expectClient: "okctl-mock-cluster-argocd",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			result, err := getRelevantUserPoolClient(
				context.Background(),
				&mockCognitoIdentityProviderAPI{clients: tc.withClients},
				"mock-userpool-id",
			)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectClient, *result.ClientName)
		})
	}
}

func TestGetCognitoClientForCluster(t *testing.T) {
	testCases := []struct {
		name             string
		withClusterName  string
		expectClientName string
	}{
		{
			name:             "Should work",
			withClusterName:  "mock-prod",
			expectClientName: "argocd",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			cluster := v1alpha1.NewCluster()
			cluster.Metadata.Name = tc.withClusterName

			client, err := getCognitoClientForCluster(
				context.Background(),
				&mockCognitoIdentityProviderAPI{
					clusterName: cluster.Metadata.Name,
					clients:     []string{fmt.Sprintf("okctl-%s-argocd", tc.withClusterName)},
				},
				cluster,
			)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectClientName, client.Name)
		})
	}
}

type mockCognitoIdentityProviderAPI struct {
	clusterName string
	clients     []string
}

func (m mockCognitoIdentityProviderAPI) ListUserPoolClientsWithContext(
	_ aws.Context,
	_ *cognitoidentityprovider.ListUserPoolClientsInput,
	_ ...request.Option,
) (*cognitoidentityprovider.ListUserPoolClientsOutput, error) {
	clients := make([]*cognitoidentityprovider.UserPoolClientDescription, len(m.clients))

	for index, name := range m.clients {
		clients[index] = &cognitoidentityprovider.UserPoolClientDescription{
			ClientId:   aws.String(name),
			ClientName: aws.String(name),
			UserPoolId: aws.String("mock-user-pool-id"),
		}
	}

	return &cognitoidentityprovider.ListUserPoolClientsOutput{
		UserPoolClients: clients,
	}, nil
}

func (m mockCognitoIdentityProviderAPI) ListUserPoolsWithContext(
	_ context.Context,
	_ *cognitoidentityprovider.ListUserPoolsInput,
	_ ...request.Option,
) (*cognitoidentityprovider.ListUserPoolsOutput, error) {
	return &cognitoidentityprovider.ListUserPoolsOutput{
		UserPools: []*cognitoidentityprovider.UserPoolDescriptionType{
			{Id: aws.String("mock-id"), Name: aws.String(m.clusterName)},
		},
	}, nil
}
