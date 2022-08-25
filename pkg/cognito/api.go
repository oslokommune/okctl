package cognito

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/aws/aws-sdk-go/service/ssm/ssmiface"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

const maximumMaximumUserPoolResults = 60

// RegisterMFADeviceOpts are the options for registering an MFA device
type RegisterMFADeviceOpts struct {
	Ctx                    context.Context
	CognitoProvider        cognitoidentityprovideriface.CognitoIdentityProviderAPI
	ParameterStoreProvider ssmiface.SSMAPI
	Cluster                v1alpha1.Cluster
	UserEmail              string
}

// RegisterMFADevice knows how to register an MFA device with a user
func RegisterMFADevice(opts RegisterMFADeviceOpts) error {
	cognitoUserPoolclient, err := getCognitoClientForCluster(opts.Ctx, opts.CognitoProvider, opts.Cluster)
	if err != nil {
		return fmt.Errorf("acquiring Cognito client ID: %w", err)
	}

	clientSecret, err := getCognitoClientSecretForClient(opts.Ctx, opts.ParameterStoreProvider, cognitoUserPoolclient.Name)
	if err != nil {
		return fmt.Errorf("acquiring Cognito client secret: %w", err)
	}

	userPassword, err := prompt(fmt.Sprintf("Enter password for user %s", opts.UserEmail), true)
	if err != nil {
		return fmt.Errorf("prompting for password: %w", err)
	}

	secretHash, err := computeSecretHash(cognitoUserPoolclient.ID, clientSecret, opts.UserEmail)
	if err != nil {
		return fmt.Errorf("computing hash: %w", err)
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
		return fmt.Errorf("initiating auth: %w", err)
	}

	if *initiateAuthResult.ChallengeName != cognitoidentityprovider.ChallengeNameTypeMfaSetup {
		return fmt.Errorf("MFA already configured for this user. Use --force to setup a new device")
	}

	associateSoftwareTokenResult, err := opts.CognitoProvider.AssociateSoftwareTokenWithContext(opts.Ctx, &cognitoidentityprovider.AssociateSoftwareTokenInput{
		Session: initiateAuthResult.Session,
	})
	if err != nil {
		return fmt.Errorf("associating: %w", err)
	}

	printDeviceSecret(os.Stdout, *associateSoftwareTokenResult.SecretCode)

	otpCode, err := prompt("Configure your MFA client with the information above and enter the one-time-password", false)
	if err != nil {
		return fmt.Errorf("prompting for OTP: %w", err)
	}

	verifySoftwareTokenResult, err := opts.CognitoProvider.VerifySoftwareTokenWithContext(opts.Ctx, &cognitoidentityprovider.VerifySoftwareTokenInput{
		FriendlyDeviceName: aws.String("code generator"),
		Session:            associateSoftwareTokenResult.Session,
		UserCode:           aws.String(otpCode),
	})
	if err != nil {
		return fmt.Errorf("verifying software token: %w", err)
	}

	if *verifySoftwareTokenResult.Status != cognitoidentityprovider.VerifySoftwareTokenResponseTypeSuccess {
		return fmt.Errorf("verifying OTP: %w", err)
	}

	fmt.Printf("Software token setup complete\n")

	return nil
}
