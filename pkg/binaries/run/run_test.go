package run_test

import (
	"bytes"
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
		name                string
		run                 *run.Run
		args                []string
		expect              interface{}
		expectErrorContains string
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
		{
			name: "If command fails, it should return an error that contains the command's exit code",
			run: func() *run.Run {
				r := run.New(nil, "working_dir", "binary_path", []string{"env_var"}, run.Cmd())
				r.CmdFn = fakeExecCommandFailure()
				return r
			}(),
			args:                []string{"binary_args"},
			expectErrorContains: "got: exit status 1",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.run.Run(ioutil.Discard, tc.args)

			if len(tc.expectErrorContains) > 0 {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectErrorContains)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, string(got))
			}
		})
	}
}

func TestRunCmdOutput(t *testing.T) {
	testCases := []struct {
		name        string
		run         *run.Run
		output      string
		expectError bool
	}{
		{
			name:   "Should get output from command stdout",
			output: "Hello!",
			run: func() *run.Run {
				cmdFn := fakeExecCommandSuccess()

				return run.New(
					nil,
					"working_dir",
					"binary_path",
					[]string{},
					cmdFn)
			}(),
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			var progress bytes.Buffer

			out, err := tc.run.Run(&progress, []string{tc.output})
			assert.NoError(t, err)

			expectedOutput := "wd=working_dir, path=binary_path, env=, args=" + tc.output

			assert.Equal(t, expectedOutput, string(out))

			expectedProgess := []byte(expectedOutput)
			assert.Equal(t, expectedProgess, progress.Bytes())
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

func TestRunProcessFailure(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}

	_, _ = fmt.Fprint(os.Stdout, strings.Join(os.Args[3:], ", "))

	os.Exit(1)
}

// fakeExecCommandSuccess is a function that initialises a new exec.Cmd, one which will
// simply call TestRunProcessSuccess rather than the command it is provided. It will
// also pass through the command and its arguments as an argument to TestRunProcessSuccess
// https://jamiethompson.me/posts/Unit-Testing-Exec-Command-In-Golang/
// Test with
func fakeExecCommandSuccess() run.CmdFn {
	return fakeExecCommand("Success")
}

func fakeExecCommandFailure() run.CmdFn {
	return fakeExecCommand("Failure")
}

func fakeExecCommand(testName string) run.CmdFn {
	return func(workingDir, path string, env, args []string) *exec.Cmd {
		cs := []string{
			"-test.run=TestRunProcess" + testName,
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
