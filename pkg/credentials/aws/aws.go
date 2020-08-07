// Package aws knows how to orchestrate a login to AWS using various methods
package aws

import (
	"bytes"
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
	"github.com/oslokommune/okctl/pkg/storage/state"
	"gopkg.in/ini.v1"
)

const (
	awsAccountIDLength = 12
)

// Credentials contains all data required for using AWS
type Credentials struct {
	AccessKeyID     string    // aws_access_key_id
	SecretAccessKey string    // aws_secret_access_key
	SessionToken    string    // aws_session_token
	SecurityToken   string    // aws_security_token
	PrincipalARN    string    // x_principal_arn
	Expires         time.Time // x_security_token_expires
	Region          string
}

// Authenticator knows how to orchestrate getting credentials
type Authenticator interface {
	Raw() (*Credentials, error)
	AsEnv() ([]string, error)
}

// Retriever knows how to retrieve credentials
type Retriever interface {
	Retrieve() (*Credentials, error)
	Invalidate()
	Valid() bool
}

// Storer knows how to store and retrieve credentials
type Storer interface {
	Save(credentials *Credentials) error
	Get() (*Credentials, error)
}

// StsProviderFn knows how to create an STS API client
type StsProviderFn func(session *session.Session) stsiface.STSAPI

// Auth stores state for fetching credentials
type Auth struct {
	Retrievers []Retriever
	Store      Storer
	creds      *Credentials
}

// AsEnv returns the AWS credentials as env vars
func (a *Auth) AsEnv() ([]string, error) {
	creds, err := a.Raw()
	if err != nil {
		return nil, err
	}

	return []string{
		fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", creds.AccessKeyID),
		fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", creds.SecretAccessKey),
		fmt.Sprintf("AWS_SESSION_TOKEN=%s", creds.SessionToken),
		fmt.Sprintf("AWS_DEFAULT_REGION=%s", creds.Region),
	}, nil
}

// Raw returns the raw credentials
func (a *Auth) Raw() (*Credentials, error) {
	// Credentials have expired
	if a.creds != nil && AreExpired(a.creds.Expires) {
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

	// Save the credentials for future use
	err := a.Store.Save(a.creds)
	if err != nil {
		return nil, err
	}

	return a.creds, nil
}

// AreExpired returns true if the credentials have expired
func AreExpired(expires time.Time) bool {
	return time.Since(expires) >= 0
}

// Resolve the available authenticators until we succeed
func (a *Auth) Resolve() (*Credentials, error) {
	var accumulatedErrors []string

	// Lets try storage first, if there is no error and
	// they aren't expired, simply return them
	creds, err := a.Store.Get()
	if err == nil && !AreExpired(creds.Expires) {
		return creds, nil
	}

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

				// Invalidate the retriever
				retriever.Invalidate()

				continue
			}

			// We just got these credentials, they shouldn't have expired already
			// which means this retriever is static or otherwise broken
			if AreExpired(creds.Expires) {
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
func New(store Storer, retriever Retriever, retrievers ...Retriever) *Auth {
	return &Auth{
		Store:      store,
		Retrievers: append([]Retriever{retriever}, retrievers...),
	}
}

// AuthStatic simply returns the provided credentials
type AuthStatic struct {
	Credentials *Credentials
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
func (a *AuthStatic) Retrieve() (*Credentials, error) {
	return a.Credentials, nil
}

// NewAuthStatic returns an initialised static authenticator
func NewAuthStatic(creds *Credentials) *AuthStatic {
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
		IsValid:      true,
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
func (a *AuthSAML) Retrieve() (*Credentials, error) {
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

	return &Credentials{
		AccessKeyID:     aws.StringValue(resp.Credentials.AccessKeyId),
		SecretAccessKey: aws.StringValue(resp.Credentials.SecretAccessKey),
		SessionToken:    aws.StringValue(resp.Credentials.SessionToken),
		SecurityToken:   aws.StringValue(resp.Credentials.SessionToken),
		PrincipalARN:    aws.StringValue(resp.AssumedRoleUser.Arn),
		Expires:         resp.Credentials.Expiration.Local(),
		Region:          a.Region,
	}, nil
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

// IniStorage contains state for storage
type IniStorage struct {
	provider state.PersisterProvider
}

// IniCredentials represents aws credentials in ini format
type IniCredentials struct {
	AccessKeyID     string    `ini:"aws_access_key_id"`
	SecretAccessKey string    `ini:"aws_secret_access_key"`
	SessionToken    string    `ini:"aws_session_token"`
	SecurityToken   string    `ini:"aws_security_token"`
	PrincipalARN    string    `ini:"x_principal_arn"`
	Expires         time.Time `ini:"x_security_token_expires"`
}

// IniConfig represents the config portion of aws config
type IniConfig struct {
	Region string `ini:"region"`
}

// Save stores the credentials
func (s *IniStorage) Save(credentials *Credentials) error {
	cfg := ini.Empty()

	profile, err := cfg.NewSection("default")
	if err != nil {
		return err
	}

	iniCreds := &IniCredentials{
		AccessKeyID:     credentials.AccessKeyID,
		SecretAccessKey: credentials.SecretAccessKey,
		SessionToken:    credentials.SessionToken,
		SecurityToken:   credentials.SecurityToken,
		PrincipalARN:    credentials.PrincipalARN,
		Expires:         credentials.Expires,
	}

	err = profile.ReflectFrom(iniCreds)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)

	_, err = cfg.WriteTo(buf)
	if err != nil {
		return err
	}
	
	return s.provider.Application().WriteToDefault("aws_credentials", buf.Bytes())
}

// Get retrieves the credentials
func (s *IniStorage) Get() (*Credentials, error) {
	data, err := s.provider.Application().ReadFromDefault("aws_credentials")
	if err != nil {
		return nil, err
	}

	creds := &IniCredentials{}

	cfg, err := ini.Load(data)
	if err != nil {
		return nil, err
	}

	err = cfg.Section("default").MapTo(creds)
	if err != nil {
		return nil, err
	}

	return &Credentials{
		AccessKeyID:     creds.AccessKeyID,
		SecretAccessKey: creds.SecretAccessKey,
		SessionToken:    creds.SessionToken,
		SecurityToken:   creds.SecurityToken,
		PrincipalARN:    creds.PrincipalARN,
		Expires:         creds.Expires,
	}, nil
}

// NewIniStorage creates a new ini storer
func NewIniStorage(provider state.PersisterProvider) *IniStorage {
	return &IniStorage{
		provider: provider,
	}
}

// InMemoryStorage stores the config in memory
type InMemoryStorage struct {
	creds *Credentials
}

// Save the credentials in memory
func (n *InMemoryStorage) Save(credentials *Credentials) error {
	n.creds = credentials

	return nil
}

// Get the credentials from memory
func (n *InMemoryStorage) Get() (*Credentials, error) {
	if n.creds == nil {
		return nil, fmt.Errorf("no credentials available")
	}

	return n.creds, nil
}

// NewInMemoryStorage creates a new in memory store, useful for tests
func NewInMemoryStorage() *InMemoryStorage {
	return &InMemoryStorage{}
}
