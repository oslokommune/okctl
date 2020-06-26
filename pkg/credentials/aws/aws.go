// Package aws knows how to orchestrate a login to AWS using various methods
package aws

import (
	"regexp"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mishudark/errors"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/credentials/aws/scrape"
)

// AuthType defines a type for enumerating authentication types
type AuthType string

const (
	// AuthTypeSAML defines SAML as an authentication method towards AWS
	AuthTypeSAML AuthType = "SAML"
)

// PopulateFn is invoked when a login is required due
// to missing or expired credentials
type PopulateFn func(l *Auth) error

// Auth stores state for fetching credentials
type Auth struct {
	AWSAccountID string
	Username     string
	Password     string
	MFAToken     string
	Region       string

	AuthType   AuthType
	PopulateFn PopulateFn
}

// Validate the authentication data
func (a *Auth) Validate() error {
	switch a.AuthType {
	case AuthTypeSAML:
		return validation.ValidateStruct(a,
			validation.Field(&a.Username,
				validation.Match(regexp.MustCompile("^byr[0-9]{6}$")).
					Error("username must match: byrXXXXXX, replacing each X with a digit"),
			),
			validation.Field(&a.Password,
				validation.Required,
			),
			validation.Field(&a.MFAToken,
				validation.Match(regexp.MustCompile("^[0-9]{6}$")).
					Error("token must consist of 6 digits"),
			),
		)
	default:
		return errors.E(errors.Errorf("unknown aws authentication method: %s", a.AuthType), errors.Invalid)
	}
}

// Get starts a process for retrieving AWS credentials
func (a *Auth) Get() (*sts.Credentials, error) {
	switch a.AuthType {
	case AuthTypeSAML:
		return a.SAML()
	default:
		return nil, errors.E(errors.Errorf("unknown aws authentication method: %s", a.AuthType), errors.Invalid)
	}
}

// SAML initiates a saml based sts authentication
func (a *Auth) SAML() (*sts.Credentials, error) {
	err := a.PopulateFn(a)
	if err != nil {
		return nil, errors.E(errors.Errorf("failed to populate required fields"))
	}

	samlAssertion, err := scrape.New().Scrape(a.Username, a.Password, a.MFAToken)
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

	svc := sts.New(sess)

	resp, err := svc.AssumeRoleWithSAML(&sts.AssumeRoleWithSAMLInput{
		PrincipalArn:  aws.String(v1alpha1.PrincipalARN(a.AWSAccountID)),
		RoleArn:       aws.String(v1alpha1.RoleARN(a.AWSAccountID)),
		SAMLAssertion: aws.String(samlAssertion),
	})
	if err != nil {
		return nil, errors.E(err, "error retrieving STS credentials using SAML", errors.Unknown)
	}

	return resp.Credentials, nil
}

// New returns an AWS credentials provider
func New(awsAccountID, region string, fn PopulateFn, t AuthType) *Auth {
	return &Auth{
		AWSAccountID: awsAccountID,
		Region:       region,
		AuthType:     t,
		PopulateFn:   fn,
	}
}

// Interactive returns a populate method that queries the user interactively
func Interactive(userName string) PopulateFn {
	return func(a *Auth) error {
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
