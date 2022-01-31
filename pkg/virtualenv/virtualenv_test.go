package virtualenv_test

import (
	"fmt"
	"testing"

	"github.com/oslokommune/okctl/pkg/virtualenv/shellgetter"

	"github.com/oslokommune/okctl/pkg/virtualenv"
	"github.com/oslokommune/okctl/pkg/virtualenv/commandlineprompter"
	"github.com/stretchr/testify/assert"
)

// nolint: funlen
func TestCreateVirtualEnvironment(t *testing.T) {
	testHelper := newTestHelper(t)

	testCases := []struct {
		name          string
		os            shellgetter.Os
		osEnvVars     map[string]string
		loginShellCmd string
		assertion     func(commandlineprompter.CommandLinePromptOpts, *virtualenv.VirtualEnvironment)
	}{
		{
			name: "Should get shell to execute from /etc/passwd",
			os:   shellgetter.OsLinux,
			osEnvVars: map[string]string{
				"OKCTL_NO_PS1": "true", // Setting this disables PS1 functionality, making the test assertions simpler
			},
			loginShellCmd: "/bin/supershell",
			assertion: func(opts commandlineprompter.CommandLinePromptOpts, venv *virtualenv.VirtualEnvironment) {
				expectedOsEnvVars := testHelper.toSlice(map[string]string{
					"OKCTL_NO_PS1": "true",
				})
				assert.Equal(t, expectedOsEnvVars, venv.Environ())
			},
		},
		{
			name: "When using bash and OKCTL_PS1, COMMAND_PROMPT should contain contents of OKCTL_PS1",
			os:   shellgetter.OsLinux,
			osEnvVars: map[string]string{
				"OKCTL_PS1": `Dir: \w $`,
			},
			loginShellCmd: "/bin/bash",
			assertion: func(opts commandlineprompter.CommandLinePromptOpts, venv *virtualenv.VirtualEnvironment) {
				expectedOsEnvVars := testHelper.toSlice(map[string]string{
					"OKCTL_PS1":      `Dir: \w $`,
					"PATH":           testHelper.ps1Dir,
					"PROMPT_COMMAND": `PS1="Dir: \w $"`,
				})
				assert.Equal(t, expectedOsEnvVars, venv.Environ())
				testHelper.assertGoldenVenvPs1(t, opts)
			},
		},
		{
			name: "When using bash, should set correct PATH and PROMPT_COMMAND, and create venv_ps1 executable",
			os:   shellgetter.OsLinux,
			osEnvVars: map[string]string{
				"PATH": "/somepath:/somepath2",
			},
			loginShellCmd: "/bin/bash",
			assertion: func(opts commandlineprompter.CommandLinePromptOpts, venv *virtualenv.VirtualEnvironment) {
				expectedOsEnvVars := testHelper.toSlice(map[string]string{
					"PATH":           fmt.Sprintf("%s:%s:%s", testHelper.ps1Dir, "/somepath", "/somepath2"),
					"PROMPT_COMMAND": fmt.Sprintf(`PS1="\[\e[0;31m\]\w \[\e[0;34m\](\$(venv_ps1 %s)) \[\e[0m\]\$ "`, opts.ClusterName),
				})
				assert.Equal(t, expectedOsEnvVars, venv.Environ())
				testHelper.assertGoldenVenvPs1(t, opts)
			},
		},
		{
			name: "When using zsh, should set correct PATH and PS1, and create venv_ps1 executable",
			os:   shellgetter.OsLinux,
			osEnvVars: map[string]string{
				"PATH": "/somepath:/somepath2",
			},
			loginShellCmd: "/bin/zsh",
			assertion: func(opts commandlineprompter.CommandLinePromptOpts, venv *virtualenv.VirtualEnvironment) {
				expectedOsEnvVars := testHelper.toSlice(map[string]string{
					"PATH": fmt.Sprintf("%s:%s:%s", testHelper.ps1Dir, "/somepath", "/somepath2"),
					"PS1":  "%F{red}%~ %f%F{blue}($(venv_ps1 myenv)%f) $ ",
				})
				assert.Equal(t, expectedOsEnvVars, venv.Environ())
				testHelper.assertGoldenVenvPs1(t, opts)
			},
		},
		{
			name: "When using a custom shell, should set correct PATH and create venv_ps1 executable",
			os:   shellgetter.OsLinux,
			osEnvVars: map[string]string{
				"PATH":        "/somepath:/somepath2",
				"OKCTL_SHELL": "/bin/fish",
			},
			loginShellCmd: "", // Not relevant since OKCTL_SHELL is primary source of shell command
			assertion: func(opts commandlineprompter.CommandLinePromptOpts, venv *virtualenv.VirtualEnvironment) {
				expectedOsEnvVars := testHelper.toSlice(map[string]string{
					"PATH":        fmt.Sprintf("%s:%s:%s", testHelper.ps1Dir, "/somepath", "/somepath2"),
					"OKCTL_SHELL": "/bin/fish",
				})
				assert.Equal(t, expectedOsEnvVars, venv.Environ())
				testHelper.assertGoldenVenvPs1(t, opts)
			},
		},
		{
			name: "When using zsh and OKCTL_PS1 is set, the environment should contain the custom PS1",
			os:   shellgetter.OsLinux,
			osEnvVars: map[string]string{
				"OKCTL_PS1": "Dir: %~ $",
			},
			loginShellCmd: "/bin/zsh",
			assertion: func(opts commandlineprompter.CommandLinePromptOpts, venv *virtualenv.VirtualEnvironment) {
				assert.Contains(t, venv.Environ(), `PS1=Dir: %~ $`)
			},
		},
		{
			name: "Should support inheriting env variables with '=' in them",
			os:   shellgetter.OsLinux,
			osEnvVars: map[string]string{
				"PATH":         "/somepath:/somepath2",
				"OKCTL_NO_PS1": "true",
				"LS_COLORS":    "rs=0:di=01;34:ln=01:*.tar=01;31",
			},
			loginShellCmd: "/bin/fish",
			assertion: func(opts commandlineprompter.CommandLinePromptOpts, venv *virtualenv.VirtualEnvironment) {
				expectedOsEnvVars := testHelper.toSlice(map[string]string{
					"PATH":         "/somepath:/somepath2",
					"OKCTL_NO_PS1": "true",
					"LS_COLORS":    "rs=0:di=01;34:ln=01:*.tar=01;31",
				})
				assert.Equal(t, expectedOsEnvVars, venv.Environ())
			},
		},
		{
			name:          "Should set PATH correctly when it does not exist already",
			os:            shellgetter.OsLinux,
			osEnvVars:     map[string]string{},
			loginShellCmd: "/bin/fish",
			assertion: func(opts commandlineprompter.CommandLinePromptOpts, venv *virtualenv.VirtualEnvironment) {
				expectedOsEnvVars := testHelper.toSlice(map[string]string{
					"PATH": testHelper.ps1Dir,
				})
				assert.Equal(t, expectedOsEnvVars, venv.Environ())
			},
		},
		{
			name: "When using bash and OKCTL_PS1, replace %env with okctl environment",
			osEnvVars: map[string]string{
				"OKCTL_PS1": `Dir: \w | \$(venv_ps1 %env) \$`,
			},
			os:            shellgetter.OsLinux,
			loginShellCmd: "/bin/bash",
			assertion: func(opts commandlineprompter.CommandLinePromptOpts, venv *virtualenv.VirtualEnvironment) {
				expectedOsEnvVars := testHelper.toSlice(map[string]string{
					"OKCTL_PS1":      `Dir: \w | \$(venv_ps1 %env) \$`,
					"PATH":           testHelper.ps1Dir,
					"PROMPT_COMMAND": fmt.Sprintf(`PS1="Dir: \w | \$(venv_ps1 %s) \$"`, opts.ClusterName),
				})
				assert.Equal(t, expectedOsEnvVars, venv.Environ())
				testHelper.assertGoldenVenvPs1(t, opts)
			},
		},
		{
			name: "When using zsh and OKCTL_PS1, replace %env with okctl environment",
			os:   shellgetter.OsLinux,
			osEnvVars: map[string]string{
				"OKCTL_PS1": `Dir: %~ | \$(venv_ps1 %env) \$`,
			},
			loginShellCmd: "/bin/zsh",
			assertion: func(opts commandlineprompter.CommandLinePromptOpts, venv *virtualenv.VirtualEnvironment) {
				expectedOsEnvVars := testHelper.toSlice(map[string]string{
					"OKCTL_PS1": `Dir: %~ | \$(venv_ps1 %env) \$`,
					"PATH":      testHelper.ps1Dir,
					"PS1":       fmt.Sprintf(`Dir: %%~ | \$(venv_ps1 %s) \$`, opts.ClusterName),
				})
				assert.Equal(t, expectedOsEnvVars, venv.Environ())
				testHelper.assertGoldenVenvPs1(t, opts)
			},
		},
		{
			name: "When using macOS, get login shell from running command dscl",
			os:   shellgetter.OsDarwin,
			assertion: func(opts commandlineprompter.CommandLinePromptOpts, venv *virtualenv.VirtualEnvironment) {
				assert.Equal(t, MacOsUserLoginShell, venv.ShellCommand)
			},
		},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			etcStorage, err := testHelper.createEtcStorage(testHelper.currentUsername, tc.loginShellCmd)
			if err != nil {
				assert.Fail(t, "could not create etc storage: %w", err)
			}

			userDirStorage := testHelper.createUserDirStorage(fmt.Sprintf("/home/%s/.okctl", testHelper.currentUsername))

			tmpStorage := testHelper.createTmpStorage()

			opts := commandlineprompter.CommandLinePromptOpts{
				Os:                   tc.os,
				MacOsUserShellGetter: testHelper.NewTestMacOsLoginShellGetter(),
				OsEnvVars:            tc.osEnvVars,
				EtcStorage:           etcStorage,
				UserDirStorage:       userDirStorage.storage,
				TmpStorage:           tmpStorage,
				ClusterName:          "myenv",
				CurrentUsername:      testHelper.currentUsername,
			}
			venv, err := virtualenv.CreateVirtualEnvironment(opts)
			assert.Nil(t, err)

			if tc.assertion != nil {
				tc.assertion(opts, venv)
			}
		})
	}
}
