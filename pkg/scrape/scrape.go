package scrape

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/foolin/pagser"
)

const (
	DefaultURL = "https://login.oslo.kommune.no/auth/realms/AD/protocol/saml/clients/amazon-aws"
)

type FormAction struct {
	URL string `pagser:"form->attr(action)"`
}

type FormSAML struct {
	Response string `pagser:"input[name='SAMLResponse']->attr(value)"`
}

func New() *scrape {
	jar, _ := cookiejar.New(nil)

	return &scrape{
		p: pagser.New(),
		c: &http.Client{
			Jar: jar,
		},
	}
}

type scrape struct {
	p *pagser.Pagser
	c *http.Client
}

func (s *scrape) doLogin(loginURL, username, password string) (*http.Response, error) {
	var formAction FormAction

	resp, err := s.c.Get(loginURL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed, got response code: %d", resp.StatusCode)
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
		return nil, fmt.Errorf("login failed, got response code: %d", resp.StatusCode)
	}

	return resp, nil
}

func (s *scrape) doTotp(resp *http.Response, mfatoken string) (*http.Response, error) {
	var formAction FormAction

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("login failed, got response code: %d", resp.StatusCode)
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
		return nil, fmt.Errorf("login failed, got response code: %d", resp.StatusCode)
	}

	return resp, nil
}

func (s *scrape) Scrape(username, password, mfaToken string) (string, error) {
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
