// Package domain validates that the provided domain is not taken
// and matches what the user expects.
package domain

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/miekg/dns"
)

// Domain contains the state for a domain
type Domain struct {
	Domain string
	FQDN   string

	In  terminal.FileReader
	Out terminal.FileWriter
	Err io.Writer
}

// NewDefaultWithSurvey returns a default domain by using a survey
func NewDefaultWithSurvey(repo, env string) (*Domain, error) {
	d := New(fmt.Sprintf("%s-%s.oslo.systems", repo, env))

	err := d.Survey()
	if err != nil {
		return nil, err
	}

	return d, nil
}

// New returns an initialised domain
func New(domain string) *Domain {
	return &Domain{
		Domain: domain,
		FQDN:   dns.Fqdn(domain),
		In:     os.Stdin,
		Out:    os.Stdout,
		Err:    os.Stderr,
	}
}

// Survey asks the user if they accept the given domain
func (d *Domain) Survey() error {
	for {
		domain, err := d.ask()
		if err != nil {
			if err == terminal.InterruptErr {
				return err
			}

			_, err = fmt.Fprintln(d.Err, err.Error())
			if err != nil {
				return err
			}

			continue
		}

		d.FQDN = dns.Fqdn(domain)
		d.Domain = domain

		return nil
	}
}

func (d *Domain) ask() (string, error) {
	domain := ""

	q := &survey.Input{
		Message: "Provide the name of the domain you want to delegate to this cluster",
		Default: d.Domain,
		Help:    "This is the domain name we will delegate to your AWS account and that the cluster will create hostnames from",
	}

	validatorFn := func(val interface{}) error {
		f, ok := val.(string)
		if !ok {
			return fmt.Errorf("could not convert input to a string")
		}

		err := Validate(f)
		if err != nil {
			return err
		}

		return NotTaken(f)
	}

	err := survey.AskOne(q, &domain,
		survey.WithStdio(d.In, d.Out, d.Err),
		survey.WithValidator(survey.Required),
		survey.WithValidator(validatorFn),
	)
	if err != nil {
		return "", err
	}

	return domain, nil
}

// Validate the provided domain
func Validate(fqdn string) error {
	return validation.Validate(&fqdn,
		is.Domain.Error(fmt.Sprintf("'%s' is not a valid domain", fqdn)),
		validation.Match(regexp.MustCompile(`^(.*)\.oslo\.systems$`)).Error(fmt.Sprintf("'%s' must end with .oslo.systems", fqdn)),
	)
}

// DNSResponse maps up the parts we are interested
// in from the response
type DNSResponse struct {
	Status int `json:"Status"`
	Answer []DNSAnswerSection
}

// DNSAnswerSection contains the answer part
type DNSAnswerSection struct {
	Name string `json:"name"`
	Type uint16 `json:"type"`
	TTL  int    `json:"TTL"`
	Data string `json:"data"`
}

// NotTaken ensures that the host is available
// Well, not really, but close enough
func NotTaken(domain string) error {
	client := &http.Client{
		Timeout: 5 * time.Second, // nolint: gomnd
	}

	// Use DNS over HTTPS service provided by google:
	// - https://developers.google.com/speed/public-dns/docs/doh/json
	req, err := http.NewRequest(http.MethodGet, "https://dns.google/resolve", nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("name", domain)
	q.Add("type", fmt.Sprintf("%d", dns.TypeNS))
	q.Add("ct", "application/x-javascript")

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return fmt.Errorf("invalid domain: %s", domain)
	case http.StatusInternalServerError:
		return fmt.Errorf("holy crap")
	}

	dnsResponse := &DNSResponse{}

	err = json.NewDecoder(resp.Body).Decode(dnsResponse)
	if err != nil {
		return err
	}

	if dnsResponse.Status != 0 {
		return fmt.Errorf("got status: %d", dnsResponse.Status)
	}

	for _, a := range dnsResponse.Answer {
		if a.Type == dns.TypeNS {
			return fmt.Errorf("domain '%s' already in use, found name servers", domain)
		}
	}

	return nil
}
