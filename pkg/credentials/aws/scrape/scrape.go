// Package scrape knows how to parse html to retrieve a SAML response from KeyCloak
package scrape

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/config/constant"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"

	"github.com/foolin/pagser"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

const (
	// DefaultURL provides a default starting location for
	// retrieving the credentials
	DefaultURL = v1alpha1.OkSamlURL
)

// Scraper defines methods for scraping SAML
type Scraper interface {
	Scrape(username, password, mfaToken string) (string, error)
}

// FormAction knows how to find the URL for
// the next stage of the parsing
type FormAction struct {
	URL string `pagser:"form->attr(action)"`
}

// FormSAML knows how to extract the SAML
// response for getting AWS credentials
type FormSAML struct {
	Response string `pagser:"input[name='SAMLResponse']->attr(value)"`
}

// FormError knows how to extract an error message
// from the HTML response
type FormError struct {
	Message string `pagser:"div[class='error_message'] span->text()"`
}

// ErrorFromResponse creates an error from an invalid http status
// code
func ErrorFromResponse(r *http.Response) error {
	pretty, _ := httputil.DumpResponse(r, true)
	return fmt.Errorf(constant.HTTPRequestError, pretty)
}

// HasError returns the error message embedded in the HTML,
// if one exists.
func HasError(p *pagser.Pagser, content string) error {
	var formError FormError

	err := p.Parse(&formError, content)
	if err != nil {
		return err
	}

	if len(formError.Message) > 0 {
		return fmt.Errorf("%s", formError.Message)
	}

	return nil
}

// New returns a scraper that knows how extract the SAML
// response for logging onto AWS using KeyCloak
func New() *Scrape {
	jar, _ := cookiejar.New(nil)

	return &Scrape{
		p: pagser.New(),
		c: &http.Client{
			Jar: jar,
		},
	}
}

// Scrape stores the state required for parsing the responses
type Scrape struct {
	p *pagser.Pagser
	c *http.Client
}

func (s *Scrape) doLogin(loginURL, username, password string) (*http.Response, error) {
	var formAction FormAction

	resp, err := s.c.Get(loginURL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrorFromResponse(resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	err = HasError(s.p, string(body))
	if err != nil {
		return nil, err
	}

	err = s.p.Parse(&formAction, string(body))
	if err != nil {
		return nil, err
	}

	resp, err = s.c.PostForm(formAction.URL, url.Values{
		"username": {username},
		"password": {password},
	})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrorFromResponse(resp)
	}

	return resp, nil
}

func (s *Scrape) doTotp(resp *http.Response, mfatoken string) (*http.Response, error) {
	var formAction FormAction

	if resp.StatusCode != http.StatusOK {
		return nil, ErrorFromResponse(resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	err = HasError(s.p, string(body))
	if err != nil {
		return nil, err
	}

	err = s.p.Parse(&formAction, string(body))
	if err != nil {
		return nil, err
	}

	resp, err = s.c.PostForm(formAction.URL, url.Values{
		"otp": {mfatoken},
	})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrorFromResponse(resp)
	}

	return resp, nil
}

// Scrape starts a process for fetching valid AWS credentials
func (s *Scrape) Scrape(username, password, mfaToken string) (string, error) {
	resp, err := s.doLogin(DefaultURL, username, password)
	if err != nil {
		return "", err
	}

	resp, err = s.doTotp(resp, mfaToken)
	if err != nil {
		return "", err
	}

	var formSAML FormSAML

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	err = resp.Body.Close()
	if err != nil {
		return "", err
	}

	err = HasError(s.p, string(body))
	if err != nil {
		return "", err
	}

	err = s.p.Parse(&formSAML, string(body))
	if err != nil {
		return "", err
	}

	return formSAML.Response, err
}
