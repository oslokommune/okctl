// Package aws knows how to orchestrate a login to AWS using various methods
package aws

import (
	"bytes"
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/spf13/afero"
	"gopkg.in/ini.v1"
)

const defaultAWSServiceUserCredentialsDuration = 24 * time.Hour

// Credentials contains all data required for using AWS
type Credentials struct {
	AwsProfile      string
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
	SecurityToken   string
	PrincipalARN    string
	Expires         time.Time
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

// Persister defines the operations required for a concrete
// implementation for persisting the credentials
type Persister interface {
	Save(credentials *Credentials) error
	Get() (*Credentials, error)
}

// StsProviderFn knows how to create an STS API client
type StsProviderFn func(session *session.Session) stsiface.STSAPI

// Auth stores state for fetching credentials
type Auth struct {
	Retrievers []Retriever
	Persister  Persister
	creds      *Credentials
}

// AsEnv returns the AWS credentials as env vars
func (a *Auth) AsEnv() ([]string, error) {
	creds, err := a.Raw()
	if err != nil {
		return nil, err
	}

	return []string{
		fmt.Sprintf("AWS_PROFILE=%s", creds.AwsProfile),
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

		// Save the credentials for future use
		err = a.Persister.Save(a.creds)
		if err != nil {
			return nil, err
		}
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
	creds, err := a.Persister.Get()
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
func New(persister Persister, retriever Retriever, retrievers ...Retriever) *Auth {
	return &Auth{
		Persister:  persister,
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

// KeyGetter defines an interface for retrieving string values based on a key
type KeyGetter func(key string) (value string)

// NewAuthEnvironment creates a retriever that fetches credentials from
// environment variables
func NewAuthEnvironment(region string, getter KeyGetter) (Retriever, error) {
	awsAccessKeyID := getter("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := getter("AWS_SECRET_ACCESS_KEY")

	if awsAccessKeyID == "" || awsSecretAccessKey == "" {
		return nil, fmt.Errorf(
			"environment variables AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY not set, see %s for more information",
			constant.DefaultAwsAuthDocumentationURL,
		)
	}

	credentials := &Credentials{
		AccessKeyID:     awsAccessKeyID,
		SecretAccessKey: awsSecretAccessKey,
		Region:          region,
		Expires:         time.Now().Add(defaultAWSServiceUserCredentialsDuration),
	}

	return &AuthStatic{
		Credentials: credentials,
		IsValid:     credentials.AccessKeyID != "" && credentials.SecretAccessKey != "",
	}, nil
}

// NewAuthProfile creates a retriever that fetches credentials from AWS profile
// environment variable
func NewAuthProfile(region string, getter KeyGetter) (Retriever, error) {
	awsProfile := getter("AWS_PROFILE")

	if awsProfile == "" {
		return nil, fmt.Errorf(
			"environment variable AWS_PROFILE not set, see %s for more information",
			constant.DefaultAwsAuthDocumentationURL,
		)
	}

	credentials := &Credentials{
		AwsProfile: awsProfile,
		Region:     region,
		Expires:    time.Now().Add(defaultAWSServiceUserCredentialsDuration),
	}

	return &AuthStatic{
		Credentials: credentials,
		IsValid:     credentials.AwsProfile != "",
	}, nil
}

// IniStorer defines the operations required for writing and reading
// the serialised credentials
type IniStorer interface {
	Write(*IniStorerData) error
	Read() (*IniStorerData, error)
}

// IniStorerData contains the data to be read and written
type IniStorerData struct {
	AwsCredentials []byte
	AwsConfig      []byte
}

// FileSystemIniStorer maintains the required state for reading and writing
// the aws credentials from a file system
type FileSystemIniStorer struct {
	FileSystem             *afero.Afero
	BaseDir                string
	AwsCredentialsFileName string
	AwsConfigFileName      string
}

// NewFileSystemIniStorer returns an initialises file system ini storer
func NewFileSystemIniStorer(awsConfigFileName, awsCredentialsFileName, baseDir string, fileSystem *afero.Afero) *FileSystemIniStorer {
	return &FileSystemIniStorer{
		FileSystem:             fileSystem,
		BaseDir:                baseDir,
		AwsCredentialsFileName: awsCredentialsFileName,
		AwsConfigFileName:      awsConfigFileName,
	}
}

// Write the data to the filesystem
func (f *FileSystemIniStorer) Write(data *IniStorerData) error {
	err := f.FileSystem.MkdirAll(f.BaseDir, 0o744)
	if err != nil {
		return err
	}

	err = f.FileSystem.WriteFile(path.Join(f.BaseDir, f.AwsConfigFileName), data.AwsConfig, 0o644)
	if err != nil {
		return err
	}

	return f.FileSystem.WriteFile(path.Join(f.BaseDir, f.AwsCredentialsFileName), data.AwsCredentials, 0o644)
}

// Read the data from the filesystem
func (f *FileSystemIniStorer) Read() (*IniStorerData, error) {
	cfg, err := f.FileSystem.ReadFile(path.Join(f.BaseDir, f.AwsConfigFileName))
	if err != nil {
		return nil, err
	}

	creds, err := f.FileSystem.ReadFile(path.Join(f.BaseDir, f.AwsCredentialsFileName))
	if err != nil {
		return nil, err
	}

	return &IniStorerData{
		AwsCredentials: creds,
		AwsConfig:      cfg,
	}, nil
}

// IniPersister knows how to serialise the credentials to a format
// compatible with the aws-cli
type IniPersister struct {
	store IniStorer
}

// NewIniPersister creates a new ini storer
func NewIniPersister(store IniStorer) *IniPersister {
	return &IniPersister{
		store: store,
	}
}

// IniCredentials serialises the credentials into a ~/.aws/credentials format
type IniCredentials struct {
	AccessKeyID     string    `ini:"aws_access_key_id"`
	SecretAccessKey string    `ini:"aws_secret_access_key"`
	SessionToken    string    `ini:"aws_session_token"`
	SecurityToken   string    `ini:"aws_security_token"`
	PrincipalARN    string    `ini:"x_principal_arn"`
	Expires         time.Time `ini:"x_security_token_expires"`
}

// IniConfig serialises the credentials into a ~/.aws/config format
type IniConfig struct {
	Region string `ini:"region"`
}

// IniProfileName sets the aws profile name, we use the default, umm, default
const IniProfileName = "default"

func serialiseAsIni(v interface{}) ([]byte, error) {
	cfg := ini.Empty()

	profile, err := cfg.NewSection(IniProfileName)
	if err != nil {
		return nil, err
	}

	err = profile.ReflectFrom(v)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	_, err = cfg.WriteTo(buf)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Save serialises and stores the provided credentials
func (s *IniPersister) Save(credentials *Credentials) error {
	creds, err := serialiseAsIni(&IniCredentials{
		AccessKeyID:     credentials.AccessKeyID,
		SecretAccessKey: credentials.SecretAccessKey,
		SessionToken:    credentials.SessionToken,
		SecurityToken:   credentials.SecurityToken,
		PrincipalARN:    credentials.PrincipalARN,
		Expires:         credentials.Expires,
	})
	if err != nil {
		return err
	}

	cfg, err := serialiseAsIni(&IniConfig{
		Region: credentials.Region,
	})
	if err != nil {
		return err
	}

	return s.store.Write(&IniStorerData{
		AwsCredentials: creds,
		AwsConfig:      cfg,
	})
}

func deserialiseFromIni(to interface{}, from interface{}) error {
	cfg, err := ini.Load(from)
	if err != nil {
		return err
	}

	return cfg.Section(IniProfileName).MapTo(to)
}

// Get retrieves credentials from store and deserializes them
func (s *IniPersister) Get() (*Credentials, error) {
	data, err := s.store.Read()
	if err != nil {
		return nil, err
	}

	creds := &IniCredentials{}

	err = deserialiseFromIni(creds, data.AwsCredentials)
	if err != nil {
		return nil, err
	}

	cfg := &IniConfig{}

	err = deserialiseFromIni(cfg, data.AwsConfig)
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
		Region:          cfg.Region,
	}, nil
}

// InMemoryPersister is useful for tests and stores the
// credentials in memory
type InMemoryPersister struct {
	creds *Credentials
}

// Save the credentials in memory
func (n *InMemoryPersister) Save(credentials *Credentials) error {
	n.creds = credentials

	return nil
}

// Get the credentials from memory
func (n *InMemoryPersister) Get() (*Credentials, error) {
	if n.creds == nil {
		return nil, fmt.Errorf("no credentials available")
	}

	return n.creds, nil
}

// NewInMemoryStorage creates a new in memory persister
func NewInMemoryStorage() *InMemoryPersister {
	return &InMemoryPersister{}
}
