// Package github knows how to retrieve valid Github credentials
package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/pkg/browser"

	"github.com/oslokommune/okctl/pkg/keyring"

	"github.com/AlecAivazis/survey/v2"

	"github.com/google/go-github/v32/github"

	oauth2device "github.com/oslokommune/okctl/pkg/oauth2"
	"golang.org/x/oauth2"
	githuboauth2 "golang.org/x/oauth2/github"
)

// DefaultGithubOauthClientID is the oauth application setup
// in the oslokommune org, this ID is considered a public
// identifier and is therefore safe to add verbatim
const DefaultGithubOauthClientID = "3e9b474f17b2bf31b07c"

// DefaultDeviceCodeURL is the default URL for entering the device
// code URL
const DefaultDeviceCodeURL = "https://github.com/login/device/code"

const (
	// CredentialsTypeDeviceFlow indicate that these are device flow credentials
	CredentialsTypeDeviceFlow = "device-flow"
	// CredentialsTypePersonalAccessToken indicate that these are personal access token
	CredentialsTypePersonalAccessToken = "personal-access-token"
)

// RequiredScopes returns the scopes required by okctl
// to perform its operations towards the Github API, see for all:
// - https://docs.github.com/en/developers/apps/scopes-for-oauth-apps
func RequiredScopes() []string {
	return []string{
		string(github.ScopeRepo),
		string(github.ScopeReadOrg),
	}
}

// ReviewURL returns the github review URL for the oauth
// permissions
func ReviewURL(clientID string) string {
	return fmt.Sprintf("https://github.com/settings/connections/applications/%s", clientID)
}

// Credentials contains the credentials
type Credentials struct {
	AccessToken string
	ClientID    string
	Type        string
}

// Authenticator provides the client interface
// for retrieving a set of valid Github credentials
type Authenticator interface {
	Raw() (*Credentials, error)
}

// Retriever defines the operations required
// for the auth orchestrator
type Retriever interface {
	Retrieve() (*Credentials, error)
	Invalidate()
	Valid() bool
}

// Persister defines the operations for storing
// and retrieving Github credentials
type Persister interface {
	Save(credentials *Credentials) error
	Get() (*Credentials, error)
}

// Auth orchestrates fetching and returning
// credentials to an end user
type Auth struct {
	Retrievers []Retriever
	Persister  Persister
	creds      *Credentials
	client     HTTPClient
}

// Raw returns the credentials as is
func (a *Auth) Raw() (*Credentials, error) {
	if a.creds != nil {
		err := AreValid(a.creds, a.client)
		if err != nil {
			return a.creds, nil
		}

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

// HTTPClient defines the http client interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// AreValid checks to see if the credentials are still good
func AreValid(credentials *Credentials, client HTTPClient) error {
	if credentials.Type == CredentialsTypePersonalAccessToken {
		// For now, lets just return
		return nil
	}

	apiURL := "https://api.github.com"

	r, err := http.NewRequest(http.MethodGet, apiURL, nil)
	if err != nil {
		return fmt.Errorf("failed to build token verification request: %w", err)
	}

	r.Header.Add("Accept", "application/vnd.github.v3+json")
	r.Header.Add("Authorization", fmt.Sprintf("token %s", credentials.AccessToken))

	resp, err := client.Do(r)
	if err != nil {
		return fmt.Errorf("failed to send token verification request: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP error %v (%v) when requesting token validation",
			resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	return nil
}

// Resolve the available authenticators until we succeed
func (a *Auth) Resolve() (*Credentials, error) {
	var accumulatedErrors []string

	// Lets try storage first, if there is no error and
	// they aren't expired, simply return them
	creds, err := a.Persister.Get()
	if err == nil { // We were able to get creds from storage
		if err := AreValid(creds, a.client); err == nil { // And they are still valid
			return creds, nil
		}
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

			// We just got these credentials, they shouldn't be invalid already
			// which means this retriever is static or otherwise broken
			if err := AreValid(creds, a.client); err != nil {
				retriever.Invalidate()

				accumulatedErrors = append(
					accumulatedErrors,
					fmt.Errorf("authenticator[%d]: invalid credentials, because: %w", i, err).Error(),
				)

				continue
			}

			return creds, nil
		}
	}

	return nil, fmt.Errorf("no valid credentials: %s", strings.Join(accumulatedErrors, ", "))
}

// New returns an initialised github authenticator
func New(persister Persister, client HTTPClient, retriever Retriever, retrievers ...Retriever) *Auth {
	return &Auth{
		Persister:  persister,
		Retrievers: append([]Retriever{retriever}, retrievers...),
		client:     client,
	}
}

// AuthStatic simply returns the provided
// credentials
type AuthStatic struct {
	Credentials *Credentials
	IsValid     bool
}

// Retrieve the stored credentials
func (a *AuthStatic) Retrieve() (*Credentials, error) {
	return a.Credentials, nil
}

// Invalidate the stored credentials
func (a *AuthStatic) Invalidate() {
	a.IsValid = false
}

// Valid returns true if the credentials
// are still valid
func (a *AuthStatic) Valid() bool {
	return a.IsValid
}

// NewAuthStatic returns an initialised static authenticator
func NewAuthStatic(creds *Credentials) *AuthStatic {
	return &AuthStatic{
		Credentials: creds,
		IsValid:     true,
	}
}

// AuthDeviceFlow contains the state required for performing
// a device flow authentication towards github
type AuthDeviceFlow struct {
	ClientID       string
	Credentials    *Credentials
	DeviceEndpoint oauth2device.DeviceEndpoint
	IsValid        bool
	ReviewURL      string
	Scopes         []string
}

// Retrieve the credentials from github
func (a *AuthDeviceFlow) Retrieve() (*Credentials, error) {
	cfg := &oauth2device.Config{
		Config: &oauth2.Config{
			ClientID: a.ClientID,
			Endpoint: githuboauth2.Endpoint,
			Scopes:   a.Scopes,
		},
		DeviceEndpoint: a.DeviceEndpoint,
	}

	client := &http.Client{}

	dcr, err := oauth2device.RequestDeviceCode(client, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve device code: %w", err)
	}

	err = a.Survey(dcr.VerificationURI, dcr.UserCode)
	if err != nil {
		return nil, fmt.Errorf("survey failed: %w", err)
	}

	accessToken, err := oauth2device.WaitForDeviceAuthorization(client, cfg, dcr)
	if err != nil {
		return nil, fmt.Errorf("failed getting device authorization: %w", err)
	}

	a.Credentials = &Credentials{
		ClientID:    cfg.ClientID,
		AccessToken: accessToken.AccessToken,
		Type:        CredentialsTypeDeviceFlow,
	}

	return a.Credentials, nil
}

// Survey queries the user to open the URL for entering the device code
func (a *AuthDeviceFlow) Survey(verificationURI, userCode string) error {
	ready := false

	prompt := &survey.Confirm{
		Message: "We will now start a Github device authentication flow, this requires entering a code in a browser window. Ready?",
		// nolint: lll
		Help:    "This process will create a github authentication token for your device, we use this token to prepare your github repository and fetch a list of teams from the organisation",
		Default: true,
	}

	err := survey.AskOne(prompt, &ready)
	if err != nil {
		return fmt.Errorf("user was not ready to continue: %w", err)
	}

	_ = browser.OpenURL(verificationURI)

	_, err = fmt.Fprintf(os.Stderr, "If a browser did not open, enter the following url in a new browser window: %s", verificationURI)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(os.Stderr, "Then enter the following code: %s", userCode)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintf(os.Stderr, "We are waiting for a response from github.com, which we will receive once you have entered the code above into the form.")
	if err != nil {
		return err
	}

	return nil
}

// Invalidate the authorisation flow
func (a *AuthDeviceFlow) Invalidate() {
	a.IsValid = false
}

// Valid returns true if the auth method is still valid
func (a *AuthDeviceFlow) Valid() bool {
	return a.IsValid
}

// NewAuthDeviceFlow returns an initialised authenticator that
// follows the device flow
func NewAuthDeviceFlow(clientID string, scopes []string) *AuthDeviceFlow {
	return &AuthDeviceFlow{
		Credentials: nil,
		IsValid:     true,
		DeviceEndpoint: oauth2device.DeviceEndpoint{
			CodeURL: DefaultDeviceCodeURL,
		},
		Scopes:    scopes,
		ClientID:  clientID,
		ReviewURL: ReviewURL(clientID),
	}
}

// KeyringPersister stores the access token in the user's keyring
type KeyringPersister struct {
	keyring keyring.Keyringer
}

// KeyringCredentialsState contains the state we store in the keyring
type KeyringCredentialsState struct {
	AccessToken string `json:"access_token"`
	ClientID    string `json:"client_id"`
	Type        string `json:"type"`
}

// Save the access token to the keyring
func (k *KeyringPersister) Save(credentials *Credentials) error {
	s := &KeyringCredentialsState{
		AccessToken: credentials.AccessToken,
		ClientID:    credentials.ClientID,
		Type:        credentials.Type,
	}

	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("failed to serialise credentials")
	}

	return k.keyring.Store(keyring.KeyTypeGithubToken, string(data))
}

// Get the access token from the keyring
func (k *KeyringPersister) Get() (*Credentials, error) {
	v, err := k.keyring.Fetch(keyring.KeyTypeGithubToken)
	if err != nil {
		return nil, err
	}

	var data KeyringCredentialsState

	err = json.Unmarshal([]byte(v), &data)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialise credentials")
	}

	return &Credentials{
		AccessToken: data.AccessToken,
		ClientID:    data.ClientID,
		Type:        data.Type,
	}, nil
}

// NewKeyringPersister returns an initialised keyring
func NewKeyringPersister(keyring keyring.Keyringer) *KeyringPersister {
	return &KeyringPersister{
		keyring: keyring,
	}
}

// InMemoryPersister stores the credentials in memory
type InMemoryPersister struct {
	Credentials *Credentials
}

// Save the credentials in memory
func (i *InMemoryPersister) Save(credentials *Credentials) error {
	i.Credentials = credentials
	return nil
}

// Get the in memory credentials
func (i *InMemoryPersister) Get() (*Credentials, error) {
	if i.Credentials != nil {
		return i.Credentials, nil
	}

	return nil, fmt.Errorf("no credentials exist")
}

// NewInMemoryPersister returns an initialised in memory persister
func NewInMemoryPersister() *InMemoryPersister {
	return &InMemoryPersister{}
}
