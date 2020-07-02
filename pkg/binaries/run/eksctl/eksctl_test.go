package eksctl_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"
	"github.com/oslokommune/okctl/pkg/binaries/run"
	"github.com/oslokommune/okctl/pkg/binaries/run/eksctl"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestEksctlDeleteCluster(t *testing.T) {
	testCases := []struct {
		name        string
		eksctl      *eksctl.Eksctl
		cfg         *v1alpha1.ClusterConfig
		expect      interface{}
		expectError bool
	}{
		{
			name: "Should work",
			cfg:  v1alpha1.NewClusterConfig(),
			eksctl: func() *eksctl.Eksctl {
				e := eksctl.New(storage.NewEphemeralStorage(), ioutil.Discard, "eksctl", []string{"AWS_SOMETHING"})

				r, ok := e.Runner.(*run.Run)
				assert.True(t, ok)

				r.CmdFn = fakeExecCommandSuccess()

				return e
			}(),
			expect: "wd=/, path=eksctl, env=AWS_SOMETHING, args=delete,cluster,--config-file,/cluster-config.yml",
		},
		{
			name: "Should fail",
			cfg:  v1alpha1.NewClusterConfig(),
			eksctl: func() *eksctl.Eksctl {
				e := eksctl.New(storage.NewEphemeralStorage(), ioutil.Discard, "eksctl", []string{"AWS_SOMETHING"})

				r, ok := e.Runner.(*run.Run)
				assert.True(t, ok)

				r.CmdFn = fakeExecCommandFailure()

				return e
			}(),
			expect:      "failed to delete: wd=/, path=eksctl, env=AWS_SOMETHING, args=delete,cluster,--config-file,/cluster-config.yml: exit status 1",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.eksctl.DeleteCluster(tc.cfg)
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

func TestEksctlCreateCluster(t *testing.T) {
	testCases := []struct {
		name        string
		eksctl      *eksctl.Eksctl
		cfg         *v1alpha1.ClusterConfig
		expect      interface{}
		expectError bool
	}{
		{
			name: "Should work",
			cfg:  v1alpha1.NewClusterConfig(),
			eksctl: func() *eksctl.Eksctl {
				e := eksctl.New(storage.NewEphemeralStorage(), ioutil.Discard, "eksctl", []string{"AWS_SOMETHING"})

				r, ok := e.Runner.(*run.Run)
				assert.True(t, ok)

				r.CmdFn = fakeExecCommandSuccess()

				return e
			}(),
			expect: "wd=/, path=eksctl, env=AWS_SOMETHING, args=create,cluster,--write-kubeconfig=false,--config-file,/cluster-config.yml",
		},
		{
			name: "Should fail",
			cfg:  v1alpha1.NewClusterConfig(),
			eksctl: func() *eksctl.Eksctl {
				e := eksctl.New(storage.NewEphemeralStorage(), ioutil.Discard, "eksctl", []string{"AWS_SOMETHING"})

				r, ok := e.Runner.(*run.Run)
				assert.True(t, ok)

				r.CmdFn = fakeExecCommandFailure()

				return e
			}(),
			expect:      "failed to create: wd=/, path=eksctl, env=AWS_SOMETHING, args=create,cluster,--write-kubeconfig=false,--config-file,/cluster-config.yml: exit status 1",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.eksctl.CreateCluster(tc.cfg)
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

func TestRunProcessFailure(t *testing.T) {
	if os.Getenv("GO_TEST_PROCESS") != "1" {
		return
	}

	_, _ = fmt.Fprint(os.Stderr, strings.Join(os.Args[3:], ", "))

	os.Exit(1)
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

func fakeExecCommandFailure() run.CmdFn {
	return func(workingDir, path string, env, args []string) *exec.Cmd {
		cs := []string{
			"-test.run=TestRunProcessFailure",
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
