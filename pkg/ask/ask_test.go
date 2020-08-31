package ask_test

import (
	"bytes"
	"testing"

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
			err = a.ConfirmPostingNameServers(msgBuf, tc.domain, tc.nameServers)
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
