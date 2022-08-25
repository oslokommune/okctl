package cognito

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"strings"
)

func acquireCognitoClientIDForCluster(ctx context.Context, provider cognitoidentityprovideriface.CognitoIdentityProviderAPI, cluster v1alpha1.Cluster) (string, error) {
	relevantUserPoolID, err := getRelevantUserPoolID(ctx, provider, cluster)
	if err != nil {
		return "", fmt.Errorf("getting relevant user pool ID: %w", err)
	}

	relevantUserPoolClientID, err := getRelevantUserPoolClient(ctx, provider, relevantUserPoolID)
	if err != nil {
		return "", fmt.Errorf("getting relevant user pool client: %w", err)
	}

	return relevantUserPoolClientID, nil
}

func getRelevantUserPoolID(ctx context.Context, provider cognitoidentityprovideriface.CognitoIdentityProviderAPI, cluster v1alpha1.Cluster) (string, error) {
	var nextToken *string = nil

	for {
		listUserPoolsResult, err := provider.ListUserPoolsWithContext(ctx, &cognitoidentityprovider.ListUserPoolsInput{
			MaxResults: aws.Int64(maximumMaximumUserPoolResults),
			NextToken:  nextToken,
		})
		if err != nil {
			return "", fmt.Errorf("listing user pools: %w", err)
		}

		for _, userPool := range listUserPoolsResult.UserPools {
			if strings.Contains(*userPool.Name, cluster.Metadata.Name) {
				return *userPool.Id, nil
			}
		}

		if listUserPoolsResult.NextToken == nil {
			break
		}

		nextToken = listUserPoolsResult.NextToken
	}

	return "", fmt.Errorf("no user pool found for cluster %s", cluster.Metadata.Name)
}

func getRelevantUserPoolClient(ctx context.Context, provider cognitoidentityprovideriface.CognitoIdentityProviderAPI, userPoolClientID string) (string, error) {
	var nextToken *string = nil

	for {
		listUserPoolsResult, err := provider.ListUserPoolClientsWithContext(ctx, &cognitoidentityprovider.ListUserPoolClientsInput{
			MaxResults: aws.Int64(1),
			NextToken:  nextToken,
			UserPoolId: aws.String(userPoolClientID),
		})
		if err != nil {
			return "", fmt.Errorf("listing user pools: %w", err)
		}

		for _, userPoolClient := range listUserPoolsResult.UserPoolClients {
			return *userPoolClient.ClientId, nil
		}

		if listUserPoolsResult.NextToken == nil {
			break
		}

		nextToken = listUserPoolsResult.NextToken
	}

	return "", fmt.Errorf("no clients found for user pool %s", userPoolClientID)
}

func computeSecretHash(clientId string, clientSecret string, username string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientId))

	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func prompt(question string, hidden bool) (string, error) {
	var result string
	var prompter survey.Prompt

	cleanQuestion := fmt.Sprintf("%s:", question)

	if hidden {
		prompter = &survey.Password{Message: cleanQuestion}
	} else {
		prompter = &survey.Input{Message: cleanQuestion}
	}

	err := survey.AskOne(prompter, &result)
	if err != nil {
		return "", fmt.Errorf("prompting: %w", err)
	}

	return result, nil
}
