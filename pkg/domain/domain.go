// Package domain validates that the provided domain is not taken
// and matches what the user expects.
package domain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/miekg/dns"
)

// Default returns the default domain name
func Default(repo, env string) string {
	return fmt.Sprintf("%s-%s.oslo.systems", repo, env)
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
	q.Add("type", fmt.Sprintf("%d", dns.TypeANY))
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

	switch dnsResponse.Status {
	case dns.RcodeSuccess:
		break
	case dns.RcodeNameError:
		return nil
	default:
		return fmt.Errorf("don't know how to handle DNS response code: %d", dnsResponse.Status)
	}

	for _, a := range dnsResponse.Answer {
		if a.Type == dns.TypeNS {
			return fmt.Errorf("domain '%s' already in use, found DNS records", domain)
		}
	}

	return nil
}

// ShouldHaveNameServers returns if there are name servers
func ShouldHaveNameServers(domain string, expectedNameservers []string) error {
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

	switch dnsResponse.Status {
	case dns.RcodeSuccess:
		break
	case dns.RcodeNameError:
		return fmt.Errorf("unable to get NS records for domain '%s', does not appear to be delegated yet", domain)
	default:
		return fmt.Errorf("don't know how to handle DNS response code: %d", dnsResponse.Status)
	}

	var nameservers []string

	for _, a := range dnsResponse.Answer {
		if a.Type == dns.TypeNS {
			nameservers = append(nameservers, a.Data)
		}
	}

	if len(nameservers) == 0 {
		return fmt.Errorf("unable to get NS records for domain '%s', does not appear to be delegated yet", domain)
	}

	diff := compare(expectedNameservers, nameservers)

	if len(diff) >= len(expectedNameservers) {
		return fmt.Errorf("nameservers do not match, expected: %s, but got: %s", expectedNameservers, nameservers)
	}

	return nil
}

// compare is copied from:
// https://gist.github.com/arxdsilva/7392013cbba7a7090cbcd120b7f5ca31
func compare(a, b []string) []string {
	for i := len(a) - 1; i >= 0; i-- {
		for _, vD := range b {
			if a[i] == vD {
				a = append(a[:i], a[i+1:]...)
				break
			}
		}
	}
	return a
}
