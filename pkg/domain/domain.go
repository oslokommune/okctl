// Package domain validates that the provided domain is not taken
// and matches what the user expects.
package domain

import (
	"fmt"
	"io"
	"os"
	"regexp"

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

// NotTaken ensures that the host is available
// Well, not really, but close enough
func NotTaken(domain string) error {
	msg := new(dns.Msg)
	msg.Id = dns.Id()
	msg.RecursionDesired = true
	msg.Question = make([]dns.Question, 1)
	msg.Question[0] = dns.Question{
		Name:   dns.Fqdn(domain),
		Qtype:  dns.TypeNS,
		Qclass: dns.ClassINET,
	}

	in, err := dns.Exchange(msg, "8.8.8.8:53")
	if err != nil {
		return err
	}

	for _, a := range in.Answer {
		if _, ok := a.(*dns.NS); ok {
			return fmt.Errorf("domain '%s' already in use, found name servers", domain)
		}
	}

	return nil
}
