package scrape

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httputil"
	"net/url"

	"github.com/foolin/pagser"
	"github.com/oslokommune/okctl/pkg/apis/okctl.io/v1alpha1"
)

const (
	DefaultURL = v1alpha1.OkSamlURL
)

type FormAction struct {
	URL string `pagser:"form->attr(action)"`
}

type FormSAML struct {
	Response string `pagser:"input[name='SAMLResponse']->attr(value)"`
}

func ErrorFromResponse(r *http.Response) error {
	pretty, _ := httputil.DumpResponse(r, true)
	return fmt.Errorf("http request failed, because: \n%s", pretty)
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

	err = s.p.ParseReader(&formAction, resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
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

	err := s.p.ParseReader(&formAction, resp.Body)
	if err != nil {
		return nil, err
	}

	err = resp.Body.Close()
	if err != nil {
		return nil, err
	}

	resp, err = s.c.PostForm(formAction.URL, url.Values{
		"totp": {mfatoken},
	})
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, ErrorFromResponse(resp)
	}

	return resp, nil
}

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

	err = s.p.ParseReader(&formSAML, resp.Body)
	if err != nil {
		return "", err
	}

	err = resp.Body.Close()

	return formSAML.Response, err
}
