// Package domain validates that the provided domain is not taken
// and matches what the user expects.
package domain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/oslokommune/okctl/pkg/config/constant"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/miekg/dns"
)

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
	Status  int `json:"Status"`
	Answer  []DNSAnswerSection
	Comment string `json:"Comment"`
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
		return fmt.Errorf(constant.InvalidDomainError, domain)
	case http.StatusInternalServerError:
		return fmt.Errorf(constant.HolyCrapError)
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
		return fmt.Errorf(constant.UnhandledDNSReponseCodeError, dnsResponse.Status)
	}

	for _, a := range dnsResponse.Answer {
		if a.Type == dns.TypeNS {
			return fmt.Errorf(constant.DomainAlreadyInUseError, domain)
		}
	}

	return nil
}

// NameServers returns the name servers for the domain
// or an empty list.
// nolint: funlen
func NameServers(domain string) ([]string, error) {
	client := &http.Client{
		Timeout: 5 * time.Second, // nolint: gomnd
	}

	// Use DNS over HTTPS service provided by google:
	// - https://developers.google.com/speed/public-dns/docs/doh/json
	req, err := http.NewRequest(http.MethodGet, "https://dns.google/resolve", nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("name", domain)
	q.Add("type", fmt.Sprintf("%d", dns.TypeNS))
	q.Add("ct", "application/x-javascript")

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	switch resp.StatusCode {
	case http.StatusBadRequest:
		return nil, fmt.Errorf(constant.InvalidDomainError, domain)
	case http.StatusInternalServerError:
		return nil, fmt.Errorf(constant.HolyCrapError)
	}

	dnsResponse := &DNSResponse{}

	err = json.NewDecoder(resp.Body).Decode(dnsResponse)
	if err != nil {
		return nil, err
	}

	switch dnsResponse.Status {
	case dns.RcodeSuccess:
		break
	case dns.RcodeServerFailure:
		return nil, fmt.Errorf("%s", dnsResponse.Comment)
	case dns.RcodeNameError:
		return nil, fmt.Errorf(constant.GetNSRecordsForDomainError, domain)
	default:
		return nil, fmt.Errorf(constant.UnhandledDNSReponseCodeError, dnsResponse.Status)
	}

	var nameservers []string

	for _, a := range dnsResponse.Answer {
		if a.Type == dns.TypeNS {
			nameservers = append(nameservers, a.Data)
		}
	}

	return nameservers, nil
}

// ShouldHaveNameServers returns if there are name servers
// nolint: funlen
func ShouldHaveNameServers(domain string, expectedNameservers []string) error {
	nameservers, err := NameServers(domain)
	if err != nil {
		return err
	}

	for i, ns := range expectedNameservers {
		expectedNameservers[i] = dns.Fqdn(ns)
	}

	if len(nameservers) == 0 {
		return fmt.Errorf(constant.GetNSRecordsForDomainError, domain)
	}

	diff := compare(expectedNameservers, nameservers)

	if len(diff) >= len(expectedNameservers) {
		return fmt.Errorf(constant.NameServerNotMatchingError, expectedNameservers, nameservers)
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
