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
	OutputFormat           string
	Force                  bool
}

// RegisterMFADevice knows how to register an MFA device with a user
func RegisterMFADevice(opts RegisterMFADeviceOpts) error {
	session, err := acquireSession(opts)
	if err != nil {
		return fmt.Errorf("acquiring session: %w", err)
	}

	associateSoftwareTokenResult, err := opts.CognitoProvider.AssociateSoftwareTokenWithContext(opts.Ctx, &cognitoidentityprovider.AssociateSoftwareTokenInput{
		Session: aws.String(session),
	})
	if err != nil {
		return fmt.Errorf("associating: %w", err)
	}

	switch {
	case opts.OutputFormat == MFAOutputFormatQRCode:
		qrCodePath, err := generateDeviceSecretQRCode(opts.Cluster, opts.UserEmail, *associateSoftwareTokenResult.SecretCode)
		if err != nil {
			return fmt.Errorf("generating QR code: %w", err)
		}

		openbrowser(qrCodePath)

		fmt.Printf("Scan the QR code in %s with your MFA device to configure your device.\n", qrCodePath)
	case opts.OutputFormat == MFAOutputFormatText:
		printDeviceSecret(os.Stdout, *associateSoftwareTokenResult.SecretCode)
	default:
		return fmt.Errorf("unknown output format: %s", opts.OutputFormat)
	}

	otpCode, err := prompt("Enter one-time-password", false)
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
