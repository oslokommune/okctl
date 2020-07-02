// Package aws knows how to orchestrate a login to AWS using various methods
package aws

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/credentials/aws/scrape"
)

const (
	awsAccountIDLength = 12
)

// Authenticator knows how to orchestrate getting credentials
type Authenticator interface {
	Raw() (*sts.Credentials, error)
	AsEnv() ([]string, error)
}

// Retriever knows how to retrieve credentials
type Retriever interface {
	Retrieve() (*sts.Credentials, error)
	Invalidate()
	Valid() bool
}

// StsProviderFn knows how to create an STS API client
type StsProviderFn func(session *session.Session) stsiface.STSAPI

// Auth stores state for fetching credentials
type Auth struct {
	Retrievers []Retriever
	creds      *sts.Credentials
}

// AsEnv returns the AWS credentials as env vars
func (a *Auth) AsEnv() ([]string, error) {
	creds, err := a.Raw()
	if err != nil {
		return nil, err
	}

	return []string{
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", *creds.AccessKeyId),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", *creds.SecretAccessKey),
		fmt.Sprintf("AWS_SESSION_TOKEN=%s", *creds.SessionToken),
	}, nil
}

// Raw returns the raw credentials
func (a *Auth) Raw() (*sts.Credentials, error) {
	// Credentials have expired
	if a.creds != nil && AreExpired(a.creds) {
		a.creds = nil
	}

	// No credentials available
	if a.creds == nil {
		creds, err := a.Resolve()
		if err != nil {
			return nil, err
		}

		a.creds = creds
	}

	return a.creds, nil
}

// AreExpired returns true if the credentials have expired
func AreExpired(creds *sts.Credentials) bool {
	return time.Since(*creds.Expiration) >= 0
}

// Resolve the available authenticators until we succeed
func (a *Auth) Resolve() (*sts.Credentials, error) {
	var accumulatedErrors []string

	for i, retriever := range a.Retrievers {
		if retriever.Valid() {
			creds, err := retriever.Retrieve()

			// We got an error, but lets just accumulate it and try the
			// next authenticator
			if err != nil {
				accumulatedErrors = append(
					accumulatedErrors,
					fmt.Sprintf("authenticator[%d]: %s", i, err.Error()),
				)

				continue
			}

			// We just got these credentials, they shouldn't have expired already
			// which means they are static or expired from an AWS credentials profile
			if AreExpired(creds) {
				retriever.Invalidate()

				accumulatedErrors = append(
					accumulatedErrors,
					fmt.Sprintf("authenticator[%d]: expired credentials", i),
				)

				continue
			}

			return creds, nil
		}
	}

	return nil, fmt.Errorf("no valid credentials: %s", strings.Join(accumulatedErrors, ", "))
}

// New returns an AWS credentials provider, it will attempt to retrieve valid credentials
// by following the retrievers in the order they are provided
func New(retriever Retriever, retrievers ...Retriever) *Auth {
	return &Auth{
		Retrievers: append([]Retriever{retriever}, retrievers...),
	}
}

// AuthStatic simply returns the provided credentials
type AuthStatic struct {
	Credentials *sts.Credentials
	IsValid     bool
}

// Invalidate the authenticator
func (a *AuthStatic) Invalidate() {
	a.IsValid = false
}

// Valid returns true if the authenticator is valid
func (a *AuthStatic) Valid() bool {
	return a.IsValid
}

// Retrieve returns the stored credentials
func (a *AuthStatic) Retrieve() (*sts.Credentials, error) {
	return a.Credentials, nil
}

// NewAuthStatic returns an initialised static authenticator
func NewAuthStatic(creds *sts.Credentials) *AuthStatic {
	return &AuthStatic{
		Credentials: creds,
		IsValid:     true,
	}
}

// DefaultStsProvider returns a standard aws sts client
func DefaultStsProvider(sess *session.Session) stsiface.STSAPI {
	return sts.New(sess)
}

// PopulateFn is invoked when a login is required due
// to missing or expired credentials
type PopulateFn func(*AuthSAML) error

// AuthSAML contains the state for performing a SAML authentication with AWS
type AuthSAML struct {
	Username     string
	Password     string
	MFAToken     string
	Region       string
	AwsAccountID string

	IsValid bool

	Scraper    scrape.Scraper
	ProviderFn StsProviderFn
	PopulateFn PopulateFn
}

// NewAuthSAML returns an instantiated authenticator towards aws with saml
func NewAuthSAML(awsAccountID, region string, scraper scrape.Scraper, providerFn StsProviderFn, fn PopulateFn) *AuthSAML {
	return &AuthSAML{
		AwsAccountID: awsAccountID,
		Region:       region,
		Scraper:      scraper,
		ProviderFn:   providerFn,
		PopulateFn:   fn,
		IsValid:      false,
	}
}

// Invalidate sets the authentication method as invalid
func (a *AuthSAML) Invalidate() {
	a.IsValid = false
}

// Valid returns the status of the authentication method
func (a *AuthSAML) Valid() bool {
	return a.IsValid
}

// Validate the SAML authentication fields
func (a *AuthSAML) Validate() error {
	return validation.ValidateStruct(a,
		validation.Field(&a.Username,
			validation.Match(regexp.MustCompile("^byr[0-9]{6}$")).
				Error("username must match: byrXXXXXX, replacing each X with a digit"),
		),
		validation.Field(&a.Password,
			validation.Required,
		),
		validation.Field(&a.Region,
			validation.Required,
		),
		validation.Field(&a.AwsAccountID,
			validation.Required,
			validation.Length(awsAccountIDLength, awsAccountIDLength),
		),
		validation.Field(&a.MFAToken,
			validation.Match(regexp.MustCompile("^[0-9]{6}$")).
				Error("token must consist of 6 digits"),
		),
	)
}

// Retrieve initiates a saml based sts authentication
func (a *AuthSAML) Retrieve() (*sts.Credentials, error) {
	err := a.PopulateFn(a)
	if err != nil {
		return nil, errors.E(err, "failed to populate required fields")
	}

	samlAssertion, err := a.Scraper.Scrape(a.Username, a.Password, a.MFAToken)
	if err != nil {
		return nil, err
	}

	if len(samlAssertion) == 0 {
		return nil, errors.E(errors.Errorf("got empty SAML assertion"), errors.Unknown)
	}

	sess, err := session.NewSession(&aws.Config{
		Region: &a.Region,
	})
	if err != nil {
		return nil, errors.E(err, "failed to create aws sts session", errors.Unknown)
	}

	svc := a.ProviderFn(sess)

	resp, err := svc.AssumeRoleWithSAML(&sts.AssumeRoleWithSAMLInput{
		PrincipalArn:  aws.String(v1alpha1.PrincipalARN(a.AwsAccountID)),
		RoleArn:       aws.String(v1alpha1.RoleARN(a.AwsAccountID)),
		SAMLAssertion: aws.String(samlAssertion),
	})
	if err != nil {
		return nil, errors.E(err, "error retrieving STS credentials using SAML", errors.Unknown)
	}

	return resp.Credentials, nil
}

// Static returns a populate method that returns the statically declared credentials
// good for testing
func Static(userName, password, mfatoken string) PopulateFn {
	return func(a *AuthSAML) error {
		a.Username = userName
		a.Password = password
		a.MFAToken = mfatoken

		return a.Validate()
	}
}

// Interactive returns a populate method that queries the user interactively
func Interactive(userName string) PopulateFn {
	return func(a *AuthSAML) error {
		qs := []*survey.Question{
			{
				Name: "username",
				Prompt: &survey.Input{
					Message: "Username:",
					Default: userName,
					Help:    "Oslo kommune username (byrXXXXXX), for authentication towards Keycloak and AWS",
				},
			},
			{
				Name: "password",
				Prompt: &survey.Password{
					Message: "Password:",
					Help:    "Oslo kommune password, for authentication towards KeyCloak and AWS",
				},
			},
			{
				Name: "mfatoken",
				Prompt: &survey.Password{
					Message: "Multi-factor authentication token:",
					Help:    "Oslo kommune multi-factor token, for authentication towards KeyCloak and AWS",
				},
			},
		}

		answers := struct {
			Username string
			Password string
			MFAToken string
		}{}

		err := survey.Ask(qs, &answers)
		if err != nil {
			return err
		}

		a.Username = answers.Username
		a.Password = answers.Password
		a.MFAToken = answers.MFAToken

		return a.Validate()
	}
}
