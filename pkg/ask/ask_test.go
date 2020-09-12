package ask_test

import (
	"bytes"
	"testing"

	"github.com/bmizerany/assert"

	"github.com/oslokommune/okctl/pkg/github"

	"github.com/AlecAivazis/survey/v2/core"
	"github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	"github.com/oslokommune/okctl/pkg/ask"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

// This test is based on the tests from:
// - https://github.com/AlecAivazis/survey
// nolint: funlen
func TestConfirmPostingNameServers(t *testing.T) {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true

	testCases := []struct {
		name        string
		domain      string
		nameServers []string
		golden      string
		procedure   func(console *expect.Console, state *vt10x.State)
	}{
		{
			name:        "Confirmed",
			domain:      "test.oslo.systems",
			nameServers: []string{"ns1.something.com", "ns2.something.com"},
			procedure: func(c *expect.Console, _ *vt10x.State) {
				_, _ = c.ExpectString("? Have you sent us the information outlined above? [? for help] (Y/n)")
				_, _ = c.SendLine("")
				_, _ = c.ExpectEOF()
			},
			golden: "message",
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

			a := ask.New()
			a.In = c.Tty()
			a.Err = c.Tty()
			a.Out = c.Tty()

			msgBuf := new(bytes.Buffer)
			_, err = a.ConfirmPostingNameServers(msgBuf, tc.domain, tc.nameServers)
			require.Nil(t, err)

			g := goldie.New(t)
			g.Assert(t, tc.golden, msgBuf.Bytes())

			// Close the slave end of the pty, and read the remaining bytes from the master end.
			_ = c.Tty().Close()
			<-donec

			t.Logf("Raw output: %q", buf.String())

			// Dump the terminal's screen.
			t.Logf("\n%s", expect.StripTrailingEmptyLines(state.String()))
		})
	}
}

// nolint: funlen lll
func TestAskSelectInfrastructureRepository(t *testing.T) {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true

	testCases := []struct {
		name      string
		def       string
		repos     []*github.Repository
		procedure func(console *expect.Console, state *vt10x.State)
		expect    interface{}
	}{
		{
			name: "Selected",
			def:  "oslokommune/else",
			repos: []*github.Repository{
				{
					FullName: github.StringPtr("oslokommune/something"),
					Private:  github.BoolPtr(true),
				},
				{
					FullName: github.StringPtr("oslokommune/else"),
					Private:  github.BoolPtr(true),
				},
			},
			procedure: func(c *expect.Console, state *vt10x.State) {
				_, _ = c.ExpectString("? Select infrastructure as code repository:  [Use arrows to move, type to filter, ? for more help]\r\n  oslokommune/something\r\n> oslokommune/else\r\n")
				_, _ = c.SendLine("")
				_, _ = c.ExpectEOF()
			},
			expect: &github.Repository{
				FullName: github.StringPtr("oslokommune/else"),
				Private:  github.BoolPtr(true),
			},
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

			a := ask.New()
			a.In = c.Tty()
			a.Err = c.Tty()
			a.Out = c.Tty()

			got, err := a.SelectInfrastructureRepository(tc.def, tc.repos)
			require.Nil(t, err)
			assert.Equal(t, tc.expect, got)

			// Close the slave end of the pty, and read the remaining bytes from the master end.
			_ = c.Tty().Close()
			<-donec

			t.Logf("Raw output: %q", buf.String())

			// Dump the terminal's screen.
			t.Logf("Lines: \n%s", expect.StripTrailingEmptyLines(state.String()))
		})
	}
}

// nolint: funlen lll
func TestAskSelectTeam(t *testing.T) {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true

	testCases := []struct {
		name      string
		teams     []*github.Team
		procedure func(console *expect.Console, state *vt10x.State)
		expect    interface{}
	}{
		{
			name: "Selected",
			teams: []*github.Team{
				{
					Name: github.StringPtr("Super team"),
				},
				{
					Name: github.StringPtr("Awesome team"),
				},
			},
			procedure: func(c *expect.Console, state *vt10x.State) {
				_, _ = c.ExpectString("? Select team that is authorised to access the Argo CD UI:  [Use arrows to move, type to filter, ? for more help]\r\n> Super team\r\n  Awesome team\r\n")
				_, _ = c.SendLine("")
				_, _ = c.ExpectEOF()
			},
			expect: &github.Team{
				Name: github.StringPtr("Super team"),
			},
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

			a := ask.New()
			a.In = c.Tty()
			a.Err = c.Tty()
			a.Out = c.Tty()

			got, err := a.SelectTeam(tc.teams)
			require.Nil(t, err)
			assert.Equal(t, tc.expect, got)

			// Close the slave end of the pty, and read the remaining bytes from the master end.
			_ = c.Tty().Close()
			<-donec

			t.Logf("Raw output: %q", buf.String())

			// Dump the terminal's screen.
			t.Logf("Lines: \n%s", expect.StripTrailingEmptyLines(state.String()))
		})
	}
}

// nolint: funlen lll
func TestAskCreateOauthApp(t *testing.T) {
	// disable color output for all prompts to simplify testing
	core.DisableColor = true

	testCases := []struct {
		name      string
		golden    string
		opts      ask.OauthAppOpts
		procedure func(console *expect.Console, state *vt10x.State)
		expect    *ask.OauthApp
	}{
		{
			name: "Confirmed",
			opts: ask.OauthAppOpts{
				Organisation: "oslokommune",
				Name:         "okctl",
				URL:          "https://something",
				CallbackURL:  "https://something/callback",
			},
			procedure: func(c *expect.Console, _ *vt10x.State) {
				_, _ = c.ExpectString("? Attempt to open browser window to github oauth apps? [? for help] (Y/n)")
				_, _ = c.SendLine("n")
				_, _ = c.ExpectString("? Enter Client ID of the oauth app:")
				_, _ = c.SendLine("client_id")
				_, _ = c.ExpectString("? Enter Client Secret of the oauth app:")
				_, _ = c.SendLine("client_secret")
				_, _ = c.ExpectEOF()
			},
			golden: "oauth-message",
			expect: &ask.OauthApp{
				Name:         "okctl",
				Organisation: "oslokommune",
				URL:          "https://something",
				CallbackURL:  "https://something/callback",
				ClientID:     "client_id",
				ClientSecret: "client_secret",
			},
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

			a := ask.New()
			a.In = c.Tty()
			a.Err = c.Tty()
			a.Out = c.Tty()

			msgBuf := new(bytes.Buffer)
			got, err := a.CreateOauthApp(msgBuf, tc.opts)
			require.Nil(t, err)
			assert.Equal(t, tc.expect, got)

			g := goldie.New(t)
			g.Assert(t, tc.golden, msgBuf.Bytes())

			// Close the slave end of the pty, and read the remaining bytes from the master end.
			_ = c.Tty().Close()
			<-donec

			t.Logf("Raw output: %q", buf.String())

			// Dump the terminal's screen.
			t.Logf("\n%s", expect.StripTrailingEmptyLines(state.String()))
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
		domain    string
		ask       *ask.Ask
		procedure func(console *expect.Console, state *vt10x.State)
		expect    string
	}{
		{
			name:   "Validate domain",
			domain: "domain-not-in-use.oslo.systems",
			ask:    ask.New(),
			procedure: func(c *expect.Console, _ *vt10x.State) {
				_, _ = c.ExpectString("? Provide the name of the domain you want to delegate to this cluster [? for help] (domain-not-in-use.oslo.systems)")
				_, _ = c.SendLine("")
				_, _ = c.ExpectEOF()
			},
			expect: "domain-not-in-use.oslo.systems",
		},
		{
			name:   "Invalid domain",
			domain: "test.oslo.com",
			ask:    ask.New(),
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
			domain: "test.oslo.systems",
			ask:    ask.New(),
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

			tc.ask.In = c.Tty()
			tc.ask.Err = c.Tty()
			tc.ask.Out = c.Tty()

			d, err := tc.ask.Domain(tc.domain)
			require.Nil(t, err)
			assert.Equal(t, tc.expect, d.Domain)

			// Close the slave end of the pty, and read the remaining bytes from the master end.
			_ = c.Tty().Close()
			<-donec

			t.Logf("Raw output: %q", buf.String())

			// Dump the terminal's screen.
			t.Logf("\n%s", expect.StripTrailingEmptyLines(state.String()))
		})
	}
}
