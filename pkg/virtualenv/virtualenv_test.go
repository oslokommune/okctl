package virtualenv_test

import (
	"fmt"
	"testing"

	"github.com/oslokommune/okctl/pkg/virtualenv"
	"github.com/oslokommune/okctl/pkg/virtualenv/commandlineprompter"
	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/assert"
)

// nolint: funlen
func TestCreateVirtualEnvironment(t *testing.T) {
	testHelper := newTestHelper(t)

	testCases := []struct {
		name            string
		osEnvVars       map[string]string
		loginShellCmd   string
		createZshrcFile bool
		assertion       func(commandlineprompter.CommandLinePromptOpts, *virtualenv.VirtualEnvironment)
	}{
		{
			name: "Should get shell to execute from /etc/passwd",
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
			osEnvVars: map[string]string{
				"PATH": "/somepath:/somepath2",
			},
			loginShellCmd: "/bin/bash",
			assertion: func(opts commandlineprompter.CommandLinePromptOpts, venv *virtualenv.VirtualEnvironment) {
				expectedOsEnvVars := testHelper.toSlice(map[string]string{
					"PATH":           fmt.Sprintf("%s:%s:%s", testHelper.ps1Dir, "/somepath", "/somepath2"),
					"PROMPT_COMMAND": fmt.Sprintf(`PS1="\[\e[0;31m\]\w \[\e[0;34m\](\$(venv_ps1 %s)) \[\e[0m\]\$ "`, opts.Environment),
				})
				assert.Equal(t, expectedOsEnvVars, venv.Environ())
				testHelper.assertGoldenVenvPs1(t, opts)
			},
		},
		{
			name: "When using zsh, should set correct PATH and ZDOTDIR, and create venv_ps1 executable",
			osEnvVars: map[string]string{
				"PATH": "/somepath:/somepath2",
			},
			loginShellCmd: "/bin/zsh",
			assertion: func(opts commandlineprompter.CommandLinePromptOpts, venv *virtualenv.VirtualEnvironment) {
				expectedOsEnvVars := testHelper.toSlice(map[string]string{
					"PATH":    fmt.Sprintf("%s:%s:%s", testHelper.ps1Dir, "/somepath", "/somepath2"),
					"ZDOTDIR": testHelper.tmpBasedir,
				})
				assert.Equal(t, expectedOsEnvVars, venv.Environ())
				testHelper.assertGoldenVenvPs1(t, opts)
			},
		},
		{
			name: "When using a custom shell, should set correct PATH and create venv_ps1 executable",
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
			name: "When using zsh, a temporary .zshrc file should have been generated",
			osEnvVars: map[string]string{
				"PATH": "/somepath:/somepath2",
			},
			loginShellCmd: "/bin/zsh",
			assertion: func(opts commandlineprompter.CommandLinePromptOpts, venv *virtualenv.VirtualEnvironment) {
				zshrc, err := opts.TmpStorage.ReadAll(".zshrc")
				assert.Nil(t, err)

				g := goldie.New(t)
				g.Assert(t, "zshrc_no_existing_zshrc", zshrc)
			},
		},
		{
			name: "When using zsh and .zshrc already exists, the tmp .zshrc should source original",
			osEnvVars: map[string]string{
				"PATH": "/somepath:/somepath2",
			},
			loginShellCmd:   "/bin/zsh",
			createZshrcFile: true,
			assertion: func(opts commandlineprompter.CommandLinePromptOpts, venv *virtualenv.VirtualEnvironment) {
				zshrc, err := opts.TmpStorage.ReadAll(".zshrc")
				assert.Nil(t, err)

				g := goldie.New(t)
				g.Assert(t, "zshrc_existing_zshrc", zshrc)
			},
		},
		{
			name: "When using zsh and OKCTL_PS1 is set, temp .zshrc should contain the custom PS1",
			osEnvVars: map[string]string{
				"OKCTL_PS1": "Dir: %~ $",
			},
			loginShellCmd: "/bin/zsh",
			assertion: func(opts commandlineprompter.CommandLinePromptOpts, venv *virtualenv.VirtualEnvironment) {
				zshrc, err := opts.TmpStorage.ReadAll(".zshrc")
				assert.Nil(t, err)

				g := goldie.New(t)
				g.Assert(t, "zshrc_custom_ps1", zshrc)
			},
		},
		{
			name: "When using zsh and ZDOTDIR is already set, a warning is returned",
			osEnvVars: map[string]string{
				"PATH":    "/somepath:/somepath2",
				"ZDOTDIR": "/somewhere",
			},
			loginShellCmd:   "/bin/zsh",
			createZshrcFile: true,
			assertion: func(opts commandlineprompter.CommandLinePromptOpts, venv *virtualenv.VirtualEnvironment) {
				assert.True(t, len(venv.Warning) > 0)
			},
		},
		{
			name: "Should support inheriting env variables with '=' in them",
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
			osEnvVars:     map[string]string{},
			loginShellCmd: "/bin/fish",
			assertion: func(opts commandlineprompter.CommandLinePromptOpts, venv *virtualenv.VirtualEnvironment) {
				expectedOsEnvVars := testHelper.toSlice(map[string]string{
					"PATH": testHelper.ps1Dir,
				})
				assert.Equal(t, expectedOsEnvVars, venv.Environ())
			},
		},
	}
	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			etcStorage, err := testHelper.CreateEtcStorage(testHelper.currentUsername, tc.loginShellCmd)
			if err != nil {
				assert.Fail(t, "could not create etc storage: %w", err)
			}

			userDirStorage := testHelper.createUserDirStorage(fmt.Sprintf("/home/%s/.okctl", testHelper.currentUsername))

			userHomeDirStorage, err := testHelper.createUserHomeDirStorage(tc.createZshrcFile)
			if err != nil {
				assert.Fail(t, "couldn't create user home dir storage: %w", err)
			}

			tmpStorage := testHelper.createTmpStorage()

			opts := commandlineprompter.CommandLinePromptOpts{
				OsEnvVars:          tc.osEnvVars,
				EtcStorage:         etcStorage,
				UserDirStorage:     userDirStorage.storage,
				UserHomeDirStorage: userHomeDirStorage,
				TmpStorage:         tmpStorage,
				Environment:        "myenv",
				CurrentUsername:    testHelper.currentUsername,
			}
			venv, err := virtualenv.CreateVirtualEnvironment(opts)
			assert.Nil(t, err)

			if tc.assertion != nil {
				tc.assertion(opts, venv)
			}
		})
	}
}
