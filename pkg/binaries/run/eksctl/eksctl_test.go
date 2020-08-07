package eksctl_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/oslokommune/okctl/pkg/api"
	"github.com/oslokommune/okctl/pkg/binaries/run"
	"github.com/oslokommune/okctl/pkg/binaries/run/eksctl"
	"github.com/oslokommune/okctl/pkg/credentials/aws"
	"github.com/oslokommune/okctl/pkg/mock"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestEksctlDeleteCluster(t *testing.T) {
	testCases := []struct {
		name        string
		eksctl      *eksctl.Eksctl
		cfg         *api.ClusterConfig
		expect      interface{}
		expectError bool
	}{
		{
			name: "Should work",
			cfg:  &api.ClusterConfig{},
			eksctl: eksctl.New(
				storage.NewEphemeralStorage(),
				ioutil.Discard,
				"eksctl",
				aws.New(aws.NewAuthStatic(mock.DefaultValidCredentials())),
				fakeExecCommandSuccess(),
			),
			// nolint: lll
			expect: "wd=/, path=eksctl, env=AWS_ACCESS_KEY_ID=ASIAV3ZUEFP6EXAMPLE,AWS_SECRET_ACCESS_KEY=XXXXXXX,AWS_SESSION_TOKEN=XXXXXXX,AWS_DEFAULT_REGION=eu-west-1, args=delete,cluster,--config-file,/cluster-config.yml",
		},
		{
			name: "Should fail",
			cfg:  &api.ClusterConfig{},
			eksctl: eksctl.New(
				storage.NewEphemeralStorage(),
				ioutil.Discard,
				"eksctl",
				aws.New(aws.NewAuthStatic(mock.DefaultValidCredentials())),
				fakeExecCommandFailure(),
			),
			// nolint: lll
			expect:      "failed to delete: wd=/, path=eksctl, env=AWS_ACCESS_KEY_ID=ASIAV3ZUEFP6EXAMPLE,AWS_SECRET_ACCESS_KEY=XXXXXXX,AWS_SESSION_TOKEN=XXXXXXX,AWS_DEFAULT_REGION=eu-west-1, args=delete,cluster,--config-file,/cluster-config.yml: exit status 1",
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
		kubePath    string
		cfg         *api.ClusterConfig
		expect      interface{}
		expectError bool
	}{
		{
			name:     "Should work",
			cfg:      &api.ClusterConfig{},
			kubePath: "/some/path",
			eksctl: eksctl.New(
				storage.NewEphemeralStorage(),
				ioutil.Discard,
				"eksctl",
				aws.New(aws.NewAuthStatic(mock.DefaultValidCredentials())),
				fakeExecCommandSuccess(),
			),
			// nolint: lll
			expect: "wd=/, path=eksctl, env=AWS_ACCESS_KEY_ID=ASIAV3ZUEFP6EXAMPLE,AWS_SECRET_ACCESS_KEY=XXXXXXX,AWS_SESSION_TOKEN=XXXXXXX,AWS_DEFAULT_REGION=eu-west-1, args=create,cluster,--write-kubeconfig=true,--kubeconfig=/some/path,--config-file,/cluster-config.yml",
		},
		{
			name:     "Should fail",
			cfg:      &api.ClusterConfig{},
			kubePath: "/some/path",
			eksctl: eksctl.New(
				storage.NewEphemeralStorage(),
				ioutil.Discard,
				"eksctl",
				aws.New(aws.NewAuthStatic(mock.DefaultValidCredentials())),
				fakeExecCommandFailure(),
			),
			// nolint: lll
			expect:      "failed to create: wd=/, path=eksctl, env=AWS_ACCESS_KEY_ID=ASIAV3ZUEFP6EXAMPLE,AWS_SECRET_ACCESS_KEY=XXXXXXX,AWS_SESSION_TOKEN=XXXXXXX,AWS_DEFAULT_REGION=eu-west-1, args=create,cluster,--write-kubeconfig=true,--kubeconfig=/some/path,--config-file,/cluster-config.yml: exit status 1",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.eksctl.CreateCluster(tc.kubePath, tc.cfg)
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
			fmt.Sprintf("env=%s", strings.Join(run.AnonymizeEnv(env), ",")),
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
			fmt.Sprintf("env=%s", strings.Join(run.AnonymizeEnv(env), ",")),
			fmt.Sprintf("args=%s", strings.Join(args, ",")),
		}

		//nolint: gosec
		cmd := exec.Command(os.Args[0], cs...)

		cmd.Env = []string{"GO_TEST_PROCESS=1"}

		return cmd
	}
}
