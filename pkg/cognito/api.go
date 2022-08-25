package cognito

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go/service/cognitoidentityprovider/cognitoidentityprovideriface"
	"github.com/logrusorgru/aurora"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
	"io"
	"os"
)

const maximumMaximumUserPoolResults = 60

// RegisterMFADevice knows how to register an MFA device with a user
func RegisterMFADevice(ctx context.Context, provider cognitoidentityprovideriface.CognitoIdentityProviderAPI, userEmail string, cluster v1alpha1.Cluster) error {
	clientID, err := acquireCognitoClientIDForCluster(ctx, provider, cluster)
	if err != nil {
		return fmt.Errorf("acquiring Cognito client ID: %w", err)
	}

	userPassword, err := prompt(fmt.Sprintf("Enter password for user %s", userEmail), true)
	if err != nil {
		return fmt.Errorf("prompting for password: %w", err)
	}

	initiateAuthResult, err := provider.InitiateAuthWithContext(ctx, &cognitoidentityprovider.InitiateAuthInput{
		AuthFlow: aws.String(cognitoidentityprovider.AuthFlowTypeUserPasswordAuth),
		AuthParameters: map[string]*string{
			"USERNAME":    aws.String(userEmail),
			"PASSWORD":    aws.String(userPassword),
			"SECRET_HASH": aws.String(computeSecretHash(clientID, clientSecret, userEmail)),
		},
		ClientId: aws.String(clientID),
	})
	if err != nil {
		return fmt.Errorf("initiating auth: %w", err)
	}

	if *initiateAuthResult.ChallengeName != cognitoidentityprovider.ChallengeNameTypeMfaSetup {
		return fmt.Errorf("MFA already configured for this user. Use --force to setup a new device")
	}

	associateSoftwareTokenResult, err := provider.AssociateSoftwareTokenWithContext(ctx, &cognitoidentityprovider.AssociateSoftwareTokenInput{
		Session: initiateAuthResult.Session,
	})
	if err != nil {
		return fmt.Errorf("associating: %w", err)
	}

	printDeviceSecret(os.Stdout, *associateSoftwareTokenResult.SecretCode)

	otpCode, err := prompt("Configure your MFA client with the information above and enter the OTP code", false)
	if err != nil {
		return fmt.Errorf("prompting for OTP: %w", err)
	}

	verifySoftwareTokenResult, err := provider.VerifySoftwareTokenWithContext(ctx, &cognitoidentityprovider.VerifySoftwareTokenInput{
		FriendlyDeviceName: aws.String("code generator"),
		Session:            associateSoftwareTokenResult.Session,
		UserCode:           aws.String(otpCode),
	})
	if err != nil {
		return fmt.Errorf("verifying software token: %w", err)
	}

	if *verifySoftwareTokenResult.Status != cognitoidentityprovider.VerifySoftwareTokenResponseTypeSuccess {
		return fmt.Errorf("verifying otp code: %w", err)
	}

	fmt.Printf("Software token setup complete\n")

	return nil
}

func printDeviceSecret(out io.Writer, secret string) {
	fmt.Fprintf(out, "Enter the following information in your MFA client:\n")
	fmt.Fprintf(out, "Code\t\t: %s\n", aurora.Green(secret))
	fmt.Fprintf(out, "Type\t\t: %s\n", aurora.Green("TOTP"))
	fmt.Fprintf(out, "Digits\t\t: %d\n", aurora.Green(6))
	fmt.Fprintf(out, "Algorithm\t: %s\n", aurora.Green("SHA1"))
	fmt.Fprintf(out, "Interval\t: %d\n", aurora.Green(30))
}
