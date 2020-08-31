// Package oauth2 implements the device flow authentication.
//
// Most of this functionality is shamelessly taken from, with some modifications:
// - https://github.com/rjw57/oauth2device
//
// Copyright (c) 2014, Rich Wareham rich.oauth2device@richwareham.com All rights reserved.
//
// redistribution and use in source and binary forms, with or without modification, are permitted
// provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this list of
//    conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions
//    and the following disclaimer in the documentation and/or other materials provided with the distribution.
//
// This software is provided by the copyright holders and contributors "as is" and any express or
// implied warranties, including, but not limited to, the implied warranties of merchantability and
// fitness for a particular purpose are disclaimed. in no event shall the copyright holder or contributors
// be liable for any direct, indirect, incidental, special, exemplary, or consequential damages (including,
// but not limited to, procurement of substitute goods or services; loss of use, data, or profits; or
// business interruption) however caused and on any theory of liability, whether in contract, strict
// liability, or tort (including negligence or otherwise) arising in any way out of the use of this software,
// even if advised of the possibility of such damage.
package oauth2

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"golang.org/x/oauth2"
)

// A DeviceCode represents the user-visible code, verification URL and
// device-visible code used to allow for user authorisation of this app. The
// app should show UserCode and VerificationURI to the user.
type DeviceCode struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationURI string `json:"verification_uri"`
	ExpiresIn       int64  `json:"expires_in"`
	Interval        int64  `json:"interval"`
}

// DeviceEndpoint contains the URLs required to initiate the OAuth2.0 flow for a
// provider's device flow.
type DeviceEndpoint struct {
	CodeURL string
}

// Config is a version of oauth2.Config augmented with device endpoints
type Config struct {
	*oauth2.Config
	DeviceEndpoint DeviceEndpoint
}

// A tokenOrError is either an OAuth2 Token response or an error indicating why
// such a response failed.
type tokenOrError struct {
	*oauth2.Token
	Error string `json:"error,omitempty"`
}

// ErrAccessDenied is an error returned when the user has denied this
// app access to their account.
var ErrAccessDenied = errors.New("access denied by user")

const (
	deviceGrantType = "urn:ietf:params:oauth:grant-type:device_code"
)

// HTTPClient defines the required http client interface
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// RequestDeviceCode will initiate the OAuth2 device authorization flow. It
// requests a device code and information on the code and URL to show to the
// user. Pass the returned DeviceCode to WaitForDeviceAuthorization.
func RequestDeviceCode(client HTTPClient, config *Config) (*DeviceCode, error) {
	form := url.Values{
		"client_id": {config.ClientID},
		"scope":     {strings.Join(config.Scopes, " ")},
	}

	r, err := http.NewRequest(http.MethodPost, config.DeviceEndpoint.CodeURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to build device code request: %w", err)
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	resp, err := client.Do(r)

	defer func() {
		err = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"request for device code authorisation returned status %v (%v)",
			resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	// Unmarshal response
	var dcr DeviceCode

	data, err := ioutil.ReadAll(resp.Body)

	err = json.NewDecoder(bytes.NewReader(data)).Decode(&dcr)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %s, because: %w", string(data), err)
	}

	return &dcr, nil
}

// WaitForDeviceAuthorization polls the token URL waiting for the user to
// authorize the app. Upon authorization, it returns the new token. If
// authorization fails then an error is returned. If that failure was due to a
// user explicitly denying access, the error is ErrAccessDenied.
//
// Modified to work with: https://docs.github.com/en/developers/apps/authorizing-oauth-apps#device-flow
func WaitForDeviceAuthorization(client HTTPClient, config *Config, code *DeviceCode) (*oauth2.Token, error) {
	form := url.Values{
		"client_id":   {config.ClientID},
		"device_code": {code.DeviceCode},
		"grant_type":  {deviceGrantType},
	}

	r, err := http.NewRequest(http.MethodPost, config.Endpoint.TokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to build access token request: %w", err)
	}

	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Content-Length", strconv.Itoa(len(form.Encode())))

	for {
		resp, err := client.Do(r)
		if err != nil {
			return nil, fmt.Errorf("failed to poll for access token: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			_ = resp.Body.Close()
			return nil, fmt.Errorf("HTTP error %v (%v) when polling for OAuth token",
				resp.StatusCode, http.StatusText(resp.StatusCode))
		}

		// Unmarshal response, checking for errors
		var token tokenOrError

		err = json.NewDecoder(resp.Body).Decode(&token)
		_ = resp.Body.Close()

		if err != nil {
			return nil, fmt.Errorf("failed to decode polling response: %w", err)
		}

		switch token.Error {
		case "":
			return token.Token, nil
		case "authorization_pending":
		case "slow_down":
			code.Interval *= 2
		case "access_denied":
			return nil, ErrAccessDenied
		default:
			return nil, fmt.Errorf("authorization failed: %v", token.Error)
		}

		time.Sleep(time.Duration(code.Interval) * time.Second)
	}
}
