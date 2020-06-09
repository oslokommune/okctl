package login

import (
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/davecgh/go-spew/spew"
	"github.com/oslokommune/okctl/pkg/scrape"
	"github.com/pkg/errors"
	"github.com/versent/saml2aws/pkg/prompter"
)

type AWSCredentials struct {
	AWSAccessKey     string    `ini:"aws_access_key_id"`
	AWSSecretKey     string    `ini:"aws_secret_access_key"`
	AWSSessionToken  string    `ini:"aws_session_token"`
	AWSSecurityToken string    `ini:"aws_security_token"`
	PrincipalARN     string    `ini:"x_principal_arn"`
	Expires          time.Time `ini:"x_security_token_expires"`
	Region           string    `ini:"region"`
}

const (
	DefaultRegion = "eu-west-1"
)

func principalARN(account string) string {
	return fmt.Sprintf("arn:aws:iam::%s:saml-provider/keycloak", account)
}

func roleARN(account string) string {
	return fmt.Sprintf("arn:aws:iam::%s:role/oslokommune/iamadmin-SAML", account)
}

type login struct {
	awsAccount string
	username   string
	password   string
	mfaToken   string
}

func (l *login) Login() (*sts.Credentials, error) {
	samlAssertion, err := scrape.New().Scrape(l.username, l.password, l.mfaToken)
	if err != nil {
		return nil, err
	}

	if len(samlAssertion) == 0 {
		return nil, fmt.Errorf("got invalid SAML assertion")
	}

	region := DefaultRegion

	sess, err := session.NewSession(&aws.Config{
		Region: &region,
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to create session")
	}

	svc := sts.New(sess)

	params := &sts.AssumeRoleWithSAMLInput{
		PrincipalArn:  aws.String(principalARN(l.awsAccount)),
		RoleArn:       aws.String(roleARN(l.awsAccount)),
		SAMLAssertion: aws.String(samlAssertion),
	}

	resp, err := svc.AssumeRoleWithSAML(params)
	if err != nil {
		return nil, errors.Wrap(err, "error retrieving STS credentials using SAML")
	}

	c := &AWSCredentials{
		AWSAccessKey:     aws.StringValue(resp.Credentials.AccessKeyId),
		AWSSecretKey:     aws.StringValue(resp.Credentials.SecretAccessKey),
		AWSSessionToken:  aws.StringValue(resp.Credentials.SessionToken),
		AWSSecurityToken: aws.StringValue(resp.Credentials.SessionToken),
		PrincipalARN:     aws.StringValue(resp.AssumedRoleUser.Arn),
		Expires:          resp.Credentials.Expiration.Local(),
		Region:           DefaultRegion,
	}

	log.Println(spew.Sdump(c))

	return resp.Credentials, nil
}

type Loginer interface {
	Login() (*sts.Credentials, error)
}

func New(account, username string) Loginer {
	return &login{
		awsAccount: account,
		username:   prompter.String("Username", username),
		password:   prompter.Password("Password"),
		mfaToken:   prompter.Password("MFA Token"),
	}
}
