package eksctl_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/sebdah/goldie/v2"

	"github.com/oslokommune/okctl/pkg/api/okctl.io/v1alpha1"

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
		clusterName string
		golden      string
		expectError bool
	}{
		{
			name:        "Should work",
			clusterName: "myCluster",
			eksctl: eksctl.New(
				nil,
				storage.NewEphemeralStorage(),
				ioutil.Discard,
				"eksctl",
				aws.New(aws.NewInMemoryStorage(), aws.NewAuthStatic(mock.DefaultValidCredentials())),
				fakeExecCommandSuccess(),
			),
			golden: "delete-cluster-works",
		},
		{
			name:        "Should fail",
			clusterName: "myCluster",
			eksctl: eksctl.New(
				nil,
				storage.NewEphemeralStorage(),
				ioutil.Discard,
				"eksctl",
				aws.New(aws.NewInMemoryStorage(), aws.NewAuthStatic(mock.DefaultValidCredentials())),
				fakeExecCommandFailure(),
			),
			golden:      "delete-cluster-fails",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			g := goldie.New(t)
			got, err := tc.eksctl.DeleteCluster(tc.clusterName)

			if tc.expectError {
				assert.Error(t, err)
				g.Assert(t, tc.golden, []byte(err.Error()))
			} else {
				assert.NoError(t, err)
				g.Assert(t, tc.golden, got)
			}
		})
	}
}

func TestEksctlCreateCluster(t *testing.T) {
	testCases := []struct {
		name        string
		eksctl      *eksctl.Eksctl
		kubePath    string
		cfg         *v1alpha1.ClusterConfig
		golden      string
		expectError bool
	}{
		{
			name: "Should work",
			eksctl: eksctl.New(
				nil,
				storage.NewEphemeralStorage(),
				ioutil.Discard,
				"eksctl",
				aws.New(aws.NewInMemoryStorage(), aws.NewAuthStatic(mock.DefaultValidCredentials())),
				fakeExecCommandSuccess(),
			),
			kubePath: "/some/path",
			cfg:      &v1alpha1.ClusterConfig{},
			golden:   "create-cluster-works",
		},
		{
			name:     "Should fail",
			cfg:      &v1alpha1.ClusterConfig{},
			kubePath: "/some/path",
			eksctl: eksctl.New(
				nil,
				storage.NewEphemeralStorage(),
				ioutil.Discard,
				"eksctl",
				aws.New(aws.NewInMemoryStorage(), aws.NewAuthStatic(mock.DefaultValidCredentials())),
				fakeExecCommandFailure(),
			),
			golden:      "create-cluster-fails",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			g := goldie.New(t)
			got, err := tc.eksctl.CreateCluster(tc.cfg)

			if tc.expectError {
				assert.Error(t, err)
				g.Assert(t, tc.golden, []byte(err.Error()))
			} else {
				assert.NoError(t, err)
				g.Assert(t, tc.golden, got)
			}
		})
	}
}

func TestEksctlHasCluster(t *testing.T) {
	testCases := []struct {
		name        string
		eksctl      *eksctl.Eksctl
		kubePath    string
		cfg         *v1alpha1.ClusterConfig
		golden      string
		expect      interface{}
		expectError bool
	}{
		{
			name:     "Should work",
			cfg:      &v1alpha1.ClusterConfig{},
			kubePath: "/some/path",
			eksctl: eksctl.New(
				nil,
				storage.NewEphemeralStorage(),
				ioutil.Discard,
				"eksctl",
				aws.New(aws.NewInMemoryStorage(), aws.NewAuthStatic(mock.DefaultValidCredentials())),
				fakeExecCommandSuccess(),
			),
			expect: true,
		},
		{
			name: "Should fail",
			cfg: &v1alpha1.ClusterConfig{
				Metadata: v1alpha1.ClusterMeta{
					Name: "test",
				},
			},
			kubePath: "/some/path",
			eksctl: eksctl.New(
				nil,
				storage.NewEphemeralStorage(),
				ioutil.Discard,
				"eksctl",
				aws.New(aws.NewInMemoryStorage(), aws.NewAuthStatic(mock.DefaultValidCredentials())),
				fakeExecCommandFailure(),
			),
			golden:      "has-cluster-fails",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			g := goldie.New(t)
			got, err := tc.eksctl.HasCluster(tc.cfg)

			if tc.expectError {
				assert.Error(t, err)
				g.Assert(t, tc.golden, []byte(err.Error()))
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expect, got)
			}
		})
	}
}

func TestEksctlCreateServiceAccount(t *testing.T) {
	testCases := []struct {
		name        string
		eksctl      *eksctl.Eksctl
		kubePath    string
		cfg         *v1alpha1.ClusterConfig
		golden      string
		expectError bool
	}{
		{
			name:     "Should work",
			cfg:      &v1alpha1.ClusterConfig{},
			kubePath: "/some/path",
			eksctl: eksctl.New(
				nil,
				storage.NewEphemeralStorage(),
				ioutil.Discard,
				"eksctl",
				aws.New(aws.NewInMemoryStorage(), aws.NewAuthStatic(mock.DefaultValidCredentials())),
				fakeExecCommandSuccess(),
			),
			golden: "create-service-account-works",
		},
		{
			name:     "Should fail",
			cfg:      &v1alpha1.ClusterConfig{},
			kubePath: "/some/path",
			eksctl: eksctl.New(
				nil,
				storage.NewEphemeralStorage(),
				ioutil.Discard,
				"eksctl",
				aws.New(aws.NewInMemoryStorage(), aws.NewAuthStatic(mock.DefaultValidCredentials())),
				fakeExecCommandFailure(),
			),
			golden:      "create-service-account-fails",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			g := goldie.New(t)
			got, err := tc.eksctl.CreateServiceAccount(tc.cfg)

			if tc.expectError {
				assert.Error(t, err)
				g.Assert(t, tc.golden, []byte(err.Error()))
			} else {
				assert.NoError(t, err)
				g.Assert(t, tc.golden, got)
			}
		})
	}
}

func TestEksctlDeleteServiceAccount(t *testing.T) {
	testCases := []struct {
		name        string
		eksctl      *eksctl.Eksctl
		kubePath    string
		cfg         *v1alpha1.ClusterConfig
		golden      string
		expectError bool
	}{
		{
			name:     "Should work",
			cfg:      &v1alpha1.ClusterConfig{},
			kubePath: "/some/path",
			eksctl: eksctl.New(
				nil,
				storage.NewEphemeralStorage(),
				ioutil.Discard,
				"eksctl",
				aws.New(aws.NewInMemoryStorage(), aws.NewAuthStatic(mock.DefaultValidCredentials())),
				fakeExecCommandSuccess(),
			),
			golden: "delete-service-account-works",
		},
		{
			name:     "Should fail",
			cfg:      &v1alpha1.ClusterConfig{},
			kubePath: "/some/path",
			eksctl: eksctl.New(
				nil,
				storage.NewEphemeralStorage(),
				ioutil.Discard,
				"eksctl",
				aws.New(aws.NewInMemoryStorage(), aws.NewAuthStatic(mock.DefaultValidCredentials())),
				fakeExecCommandFailure(),
			),
			golden:      "delete-service-account-fails",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			g := goldie.New(t)
			got, err := tc.eksctl.DeleteServiceAccount(tc.cfg)

			if tc.expectError {
				assert.Error(t, err)
				g.Assert(t, tc.golden, []byte(err.Error()))
			} else {
				assert.NoError(t, err)
				g.Assert(t, tc.golden, got)
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
