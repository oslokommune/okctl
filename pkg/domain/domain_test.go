package domain_test

import (
	"bytes"
	"testing"

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/oslokommune/okctl/pkg/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotTaken(t *testing.T) {
	testCases := []struct {
		name        string
		fqdn        string
		expectError bool
		expect      interface{}
	}{
		{
			name: "Available domain",
			fqdn: "nosuchsubdomain.oslo.systems",
		},
		{
			name:        "Taken domain",
			fqdn:        "test.oslo.systems",
			expectError: true,
			expect:      "domain 'test.oslo.systems' already in use, found DNS records",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := domain.NotTaken(tc.fqdn)
			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidate(t *testing.T) {
	testCases := []struct {
		name        string
		fqdn        string
		expectError bool
		expect      interface{}
	}{
		{
			name: "Valid domain",
			fqdn: "test.oslo.systems",
		},
		{
			name:        "Invalid domain",
			fqdn:        "not a domain.oslo.systems",
			expectError: true,
			expect:      "'not a domain.oslo.systems' is not a valid domain",
		},
		{
			name:        "Valid domain, doesn't end with oslo.systems",
			fqdn:        "some.other.domain.com",
			expectError: true,
			expect:      "'some.other.domain.com' must end with .oslo.systems",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := domain.Validate(tc.fqdn)
			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// This test is based on the tests from:
// - https://github.com/AlecAivazis/survey
// nolint: funlen
func TestDomainSurvey(t *testing.T) {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true

	testCases := []struct {
		name      string
		domain    *domain.Domain
		procedure func(console *expect.Console, state *vt10x.State)
		expect    string
	}{
		{
			name:   "Valid domain",
			domain: domain.New("domain-not-in-use.oslo.systems"),
			procedure: func(c *expect.Console, _ *vt10x.State) {
				_, _ = c.ExpectString("? Provide the name of the domain you want to delegate to this cluster [? for help] (domain-not-in-use.oslo.systems)")
				_, _ = c.SendLine("")
				_, _ = c.ExpectEOF()
			},
			expect: "domain-not-in-use.oslo.systems",
		},
		{
			name:   "Invalid domain",
			domain: domain.New("test.oslo.com"),
			procedure: func(c *expect.Console, _ *vt10x.State) {
				_, _ = c.ExpectString("? Provide the name of the domain you want to delegate to this cluster [? for help] (test.oslo.com)")
				_, _ = c.SendLine("")
				_, _ = c.ExpectString("X Sorry, your reply was invalid: 'test.oslo.com' must end with .oslo.systems")
				_, _ = c.ExpectString("? Provide the name of the domain you want to delegate to this cluster [? for help] (test.oslo.com)")
				_, _ = c.SendLine("domain-not-in-use.oslo.systems")
				_, _ = c.ExpectEOF()
			},
			expect: "domain-not-in-use.oslo.systems",
		},
		{
			name:   "Domain taken",
			domain: domain.New("test.oslo.systems"),
			procedure: func(c *expect.Console, _ *vt10x.State) {
				_, _ = c.ExpectString("? Provide the name of the domain you want to delegate to this cluster [? for help] (test.oslo.systems)")
				_, _ = c.SendLine("")
				_, _ = c.ExpectString("X Sorry, your reply was invalid: domain 'test.oslo.systems' already in use, found DNS records")
				_, _ = c.ExpectString("? Provide the name of the domain you want to delegate to this cluster [? for help] (test.oslo.systems)")
				_, _ = c.SendLine("domain-not-in-use.oslo.systems")
				_, _ = c.ExpectEOF()
			},
			expect: "domain-not-in-use.oslo.systems",
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			c, state, err := vt10x.NewVT10XConsole(
				expect.WithStdout(buf),
				// Uncomment this line to get debug output:
				// expect.WithLogger(log.New(os.Stdout, "state", log.LstdFlags)),
			)
			require.Nil(t, err)
			defer func() {
				_ = c.Close()
			}()

			donec := make(chan struct{})
			go func() {
				defer close(donec)
				tc.procedure(c, state)
			}()

			tc.domain.In = c.Tty()
			tc.domain.Err = c.Tty()
			tc.domain.Out = c.Tty()

			err = tc.domain.Survey()
			require.Nil(t, err)
			assert.Equal(t, tc.expect, tc.domain.Domain)

			// Close the slave end of the pty, and read the remaining bytes from the master end.
			_ = c.Tty().Close()
			<-donec

			t.Logf("Raw output: %q", buf.String())

			// Dump the terminal's screen.
			t.Logf("\n%s", expect.StripTrailingEmptyLines(state.String()))
		})
	}
}
