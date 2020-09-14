package run_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/oslokommune/okctl/pkg/binaries/run"
	"github.com/stretchr/testify/assert"
)

func TestManual(t *testing.T) {
	testCases := []struct {
		name    string
		command run.CmdFn
	}{
		{
			name:    "Should work",
			command: run.Cmd(),
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Skip("skipping manual test, used for playing with implementation")
			r := run.New(nil, ".", "/usr/bin/env", nil, tc.command)
			_, err := r.Run(ioutil.Discard, nil)
			assert.NoError(t, err)
		})
	}
}

func TestAnonymizeEnv(t *testing.T) {
	testCases := []struct {
		name    string
		entries []string
		expect  []string
	}{
		{
			name: "Hide AWS env vars",
			entries: []string{
				"AWS_SECRET_ACCESS_KEY=something",
				"AWS_SESSION_TOKEN=something",
			},
			expect: []string{
				"AWS_SECRET_ACCESS_KEY=XXXXXXX",
				"AWS_SESSION_TOKEN=XXXXXXX",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := run.AnonymizeEnv(tc.entries)

			assert.Equal(t, tc.expect, got)
		})
	}
}

func TestRun(t *testing.T) {
	testCases := []struct {
		name        string
		run         *run.Run
		args        []string
		expect      interface{}
		expectError bool
	}{
		{
			name: "Should work",
			run: func() *run.Run {
				r := run.New(nil, "working_dir", "binary_path", []string{"env_var"}, run.Cmd())
				r.CmdFn = fakeExecCommandSuccess()
				return r
			}(),
			args:   []string{"binary_args"},
			expect: "wd=working_dir, path=binary_path, env=env_var, args=binary_args",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.run.Run(ioutil.Discard, tc.args)

			if tc.expectError {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, string(got))
			}
		})
	}
}

func TestRunProcessSuccess(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}

	_, _ = fmt.Fprint(os.Stdout, strings.Join(os.Args[3:], ", "))

	os.Exit(0)
}

// fakeExecCommandSuccess is a function that initialises a new exec.Cmd, one which will
// simply call TestRunProcessSuccess rather than the command it is provided. It will
// also pass through the command and its arguments as an argument to TestRunProcessSuccess
// https://jamiethompson.me/posts/Unit-Testing-Exec-Command-In-Golang/
func fakeExecCommandSuccess() run.CmdFn {
	return func(workingDir, path string, env, args []string) *exec.Cmd {
		cs := []string{
			"-test.run=TestRunProcessSuccess",
			"--",
			fmt.Sprintf("wd=%s", workingDir),
			fmt.Sprintf("path=%s", path),
			fmt.Sprintf("env=%s", strings.Join(env, ",")),
			fmt.Sprintf("args=%s", strings.Join(args, ",")),
		}

		//nolint: gosec
		cmd := exec.Command(os.Args[0], cs...)

		cmd.Env = []string{"GO_TEST_PROCESS=1"}

		return cmd
	}
}
