package cognito

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/aws/aws-sdk-go/aws/request"

	"github.com/google/uuid"
	qrcode "github.com/skip2/go-qrcode"

	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/logrusorgru/aurora"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

const (
	// MFAOutputFormatQRCode represents showing the MFA details as a QR code
	MFAOutputFormatQRCode = "qrcode"
	// MFAOutputFormatText represents showing the MFA details as text
	MFAOutputFormatText                  = "text"
	defaultOneTimePasswordType           = "TOTP"
	defaultOneTimePasswordDigits         = 6
	defaultOneTimePasswordAlgorithm      = "SHA1"
	defaultOneTimePasswordInterval       = 30
	defaultRelevantUserPoolClientKeyword = "argocd"
	defaultQRCodePixelSize               = 256
)

type userPoolClient struct {
	Name string
	ID   string
}

func acquireSession(opts RegisterMFADeviceOpts) (string, error) {
	cognitoUserPoolclient, err := getCognitoClientForCluster(opts.Ctx, opts.CognitoProvider, opts.Cluster)
	if err != nil {
		return "", fmt.Errorf("acquiring Cognito client ID: %w", err)
	}

	clientSecret, err := getCognitoClientSecretForClient(opts.Ctx, opts.ParameterStoreProvider, cognitoUserPoolclient.Name)
	if err != nil {
		return "", fmt.Errorf("acquiring Cognito client secret: %w", err)
	}

	userPassword, err := prompt(fmt.Sprintf("Enter password for Cognito user %s", opts.UserEmail), true)
	if err != nil {
		return "", fmt.Errorf("prompting for password: %w", err)
	}

	secretHash, err := computeSecretHash(cognitoUserPoolclient.ID, clientSecret, opts.UserEmail)
	if err != nil {
		return "", fmt.Errorf("computing hash: %w", err)
	}

	initiateAuthResult, err := opts.CognitoProvider.InitiateAuthWithContext(opts.Ctx, &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String(cognitoidentityprovider.AuthFlowTypeUserPasswordAuth),
		AuthParameters: map[string]*string{
			"USERNAME":    aws.String(opts.UserEmail),
			"PASSWORD":    aws.String(userPassword),
			"SECRET_HASH": aws.String(secretHash),
		},
		ClientId: aws.String(cognitoUserPoolclient.ID),
	})
	if err != nil {
		return "", fmt.Errorf("initiating auth: %w", err)
	}

	if *initiateAuthResult.ChallengeName != cognitoidentityprovider.ChallengeNameTypeMfaSetup {
		return "", fmt.Errorf("MFA already configured for this user")
	}

	return *initiateAuthResult.Session, nil
}

func getCognitoClientForCluster(ctx context.Context, provider cognitoidentityprovideriface.CognitoIdentityProviderAPI, cluster v1alpha1.Cluster) (userPoolClient, error) {
	relevantUserPoolID, err := getRelevantUserPoolID(ctx, provider, cluster)
	if err != nil {
		return userPoolClient{}, fmt.Errorf("getting relevant user pool ID: %w", err)
	}

	relevantUserPoolClient, err := getRelevantUserPoolClient(ctx, provider, relevantUserPoolID)
	if err != nil {
		return userPoolClient{}, fmt.Errorf("getting relevant user pool client: %w", err)
	}

	return userPoolClient{Name: *relevantUserPoolClient.ClientName, ID: *relevantUserPoolClient.ClientId}, nil
}

func getRelevantUserPoolID(ctx context.Context, provider cognitoidentityprovideriface.CognitoIdentityProviderAPI, cluster v1alpha1.Cluster) (string, error) {
	var nextToken *string

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

func getRelevantUserPoolClient(ctx context.Context, provider userpoolClientsLister, userPoolID string) (
	cognitoidentityprovider.UserPoolClientDescription,
	error,
) {
	var nextToken *string

	for {
		listUserPoolsResult, err := provider.ListUserPoolClientsWithContext(ctx, &cognitoidentityprovider.ListUserPoolClientsInput{
			MaxResults: aws.Int64(maximumMaximumUserPoolResults),
			NextToken:  nextToken,
			UserPoolId: aws.String(userPoolID),
		})
		if err != nil {
			return cognitoidentityprovider.UserPoolClientDescription{}, fmt.Errorf("listing user pools: %w", err)
		}

		for _, client := range listUserPoolsResult.UserPoolClients {
			// This should be any of the clients provisioned by okctl, but due to inconsistent naming of the Grafana client
			// secret SSM parameter and the situation regarding golden path we'll settle on picking the ArgoCD client for
			// now. This breaks MFA for environments without ArgoCD.
			if strings.Contains(*client.ClientName, defaultRelevantUserPoolClientKeyword) {
				return *client, nil
			}
		}

		if listUserPoolsResult.NextToken == nil {
			break
		}

		nextToken = listUserPoolsResult.NextToken
	}

	return cognitoidentityprovider.UserPoolClientDescription{}, fmt.Errorf("no clients found for user pool %s", userPoolID)
}

func getCognitoClientSecretForClient(ctx context.Context, provider ssmiface.SSMAPI, clientName string) (string, error) {
	parameterPath := fmt.Sprintf("/%s/client_secret", strings.ReplaceAll(clientName, "-", "/"))

	getParameterResult, err := provider.GetParameterWithContext(ctx, &ssm.GetParameterInput{
		Name:           aws.String(parameterPath),
		WithDecryption: aws.Bool(true),
	})
	if err != nil {
		return "", fmt.Errorf("retrieving parameter: %w", err)
	}

	return *getParameterResult.Parameter.Value, nil
}

func computeSecretHash(clientID string, clientSecret string, username string) (string, error) {
	mac := hmac.New(sha256.New, []byte(clientSecret))

	_, err := mac.Write([]byte(username + clientID))
	if err != nil {
		return "", fmt.Errorf("writing payload: %w", err)
	}

	return base64.StdEncoding.EncodeToString(mac.Sum(nil)), nil
}

func prompt(question string, hidden bool) (string, error) {
	var (
		result   string
		prompter survey.Prompt
	)

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

func printDeviceSecret(out io.Writer, secret string) {
	fmt.Fprintf(out, "Enter the following information in your MFA client:\n")
	fmt.Fprintf(out, "Code\t\t: %s\n", aurora.Green(secret))
	fmt.Fprintf(out, "Type\t\t: %s\n", aurora.Green(defaultOneTimePasswordType))
	fmt.Fprintf(out, "Digits\t\t: %d\n", aurora.Green(defaultOneTimePasswordDigits))
	fmt.Fprintf(out, "Algorithm\t: %s\n", aurora.Green(defaultOneTimePasswordAlgorithm))
	fmt.Fprintf(out, "Interval\t: %d\n", aurora.Green(defaultOneTimePasswordInterval))
}

func generateDeviceSecretQRCode(cluster v1alpha1.Cluster, userEmail string, secret string) (string, error) {
	qrCodePath := path.Join(os.TempDir(), fmt.Sprintf("%s.png", uuid.Must(uuid.NewUUID()).String()))

	qrCodeURI := fmt.Sprintf("otpauth://totp/%s:%s?secret=%s&issuer=%s&algorithm=%s&digits=%d&period=%d",
		cluster.Metadata.Name,
		userEmail,
		secret,
		cluster.Metadata.Name,
		defaultOneTimePasswordAlgorithm,
		defaultOneTimePasswordDigits,
		defaultOneTimePasswordInterval,
	)

	err := qrcode.WriteFile(qrCodeURI, qrcode.Medium, defaultQRCodePixelSize, qrCodePath)
	if err != nil {
		return "", fmt.Errorf("writing QR code: %w", err)
	}

	return qrCodePath, nil
}

//nolint:gosec
func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		log.Fatal(err)
	}
}

type userpoolClientsLister interface {
	ListUserPoolClientsWithContext(
		context.Context,
		*cognitoidentityprovider.ListUserPoolClientsInput,
		...request.Option,
	) (*cognitoidentityprovider.ListUserPoolClientsOutput, error)
}
