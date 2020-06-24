// Package login knows how to orchestrate a login to AWS via KeyCloak
package login

import (
	"fmt"
	"regexp"

	"github.com/AlecAivazis/survey/v2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/credentials/scrape"
	"github.com/pkg/errors"
)

// Loginer knows how to retrieve AWS credentials
type Loginer interface {
	Login() (*sts.Credentials, error)
}

// Login stores state for fetching credentials
type Login struct {
	AWSAccountID string
	Username     string
	Password     string
	MFAToken     string
	Region       string
}

// Validate the login data
func (l *Login) Validate() error {
	return validation.ValidateStruct(l,
		validation.Field(&l.Username,
			validation.Match(regexp.MustCompile("^byr[0-9]{6}$")).
				Error("username must match: byrXXXXXX, replacing each X with a digit"),
		),
		validation.Field(&l.Password,
			validation.Required,
		),
		validation.Field(&l.MFAToken,
			validation.Match(regexp.MustCompile("^[0-9]{6}$")).
				Error("token must consist of 6 digits"),
		),
	)
}

// Login starts a process for retrieving AWS credentials
func (l *Login) Login() (*sts.Credentials, error) {
	samlAssertion, err := scrape.New().Scrape(l.Username, l.Password, l.MFAToken)
	if err != nil {
		return nil, err
	}

	if len(samlAssertion) == 0 {
		return nil, fmt.Errorf("got invalid SAML assertion")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: &l.Region,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create session")
	}

	svc := sts.New(sess)

	resp, err := svc.AssumeRoleWithSAML(&sts.AssumeRoleWithSAMLInput{
		PrincipalArn:  aws.String(v1alpha1.PrincipalARN(l.AWSAccountID)),
		RoleArn:       aws.String(v1alpha1.RoleARN(l.AWSAccountID)),
		SAMLAssertion: aws.String(samlAssertion),
	})
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving STS credentials using SAML")
	}

	return resp.Credentials, nil
}

// Interactive returns an interactive login
func Interactive(awsAccountID, region, userName string) (*Login, error) {
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
		return nil, err
	}

	l := &Login{
		AWSAccountID: awsAccountID,
		Region:       region,
		Username:     answers.Username,
		Password:     answers.Password,
		MFAToken:     answers.MFAToken,
	}

	return l, l.Validate()
}
