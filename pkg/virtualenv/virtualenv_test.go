package virtualenv_test

import (
	"fmt"
	"github.com/oslokommune/okctl/pkg/storage"
	"github.com/oslokommune/okctl/pkg/virtualenv"
	"github.com/oslokommune/okctl/pkg/virtualenv/commandlineprompter"
	"github.com/stretchr/testify/assert"
	"path"
	"testing"
)

// TODO: Test:
// 0 Tabletest, with OKCTL_NO_PS1=true to make test easier.
// [ ] venv.Environ() should contain os.Environ()
// [ ] merged path

// 1
// Tabletest, shellCommand er som forventet (fra /etc/passwd)
// [ ] /etc/passwd -> whateverShell

// 2
// Tabletest, PS1 er satt som forventet
// [ ] login shell = zsh
// [ ] login shell = bash
// [ ] login shell = sh

// Etter 1+2, få QA.

// [ ] OKCTL_SHELL ->
//     shellCommand er som forventet.
// 	   ingen OKCTL_PS1 er satt
// 4 [x] OKCTL_NO_PS1
// 5 [ ] OKCTL_PS1
//         [ ] Settes for bash
//         [ ] Settes for zsh
//         [ ] Settes for unknown
// 5.5 [ ] Scann PATH etter venv_ps1

// ZSH special cases
// 3 [ ] ZDOTDIR set correctly, and that dir contains .zshrc file with the expected contents.
// 		Vurder om denne trengs hvis de to nedenfor er satt.
// 6 Tabletest
//   [ ] If ~/.zshrc exists, /tmp/x/.zshrc contains "source ~/.zshrc"
//   [ ] If not ~/.zshrc exists, /tmp/x/.zshrc doesn't contain "source ~/.zshrc"
// 7 [ ] ZDOTDIR already set -> gir warning

// TODO path 				"LS_COLORS": "rs=0:di=01;34:ln=01:*.tar=01;31",

func TestCreateVirtualEnvironment(t *testing.T) {
	testHelper := NewTestHelper(t)

	testCases := []struct {
		name              string
		osEnvVars         map[string]string
		loginShellCmd     string
		expectedOsEnvVars map[string]string
		customAssert      func(commandlineprompter.CommandLinePromptOpts, *virtualenv.VirtualEnvironment)
	}{
		{ // TODO consider replacing this with "Scann PATH etter venv_ps1"
			name: "Should return merged PATH starting with path of a new PS1 executable",
			osEnvVars: map[string]string{
				"PATH":        "/somepath:/somepath2",
				"OKCTL_PS1":   "(some ps1 so it's simple to do the test assertion)",
				"OKCTL_SHELL": "/bin/fish",
			},
			loginShellCmd: "", // Not relevant since OKCTL_SHELL is primary source of shell command
			expectedOsEnvVars: map[string]string{
				"PATH":        fmt.Sprintf("%s:%s:%s", testHelper.ps1Dir, "/somepath", "/somepath2"),
				"OKCTL_PS1":   "(some ps1 so it's simple to do the test assertion)",
				"OKCTL_SHELL": "/bin/fish",
			},
			customAssert: func(opts commandlineprompter.CommandLinePromptOpts, venv *virtualenv.VirtualEnvironment) {
				// Make sure executable venv_ps1 is a file that exists on the PATH
				content, err := opts.UserDirStorage.ReadAll(path.Join(commandlineprompter.Ps1Dir, commandlineprompter.Ps1Filename))
				assert.Nil(t, err)

				// TODO jeg kom hit

			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			osEnvVars := tc.osEnvVars

			etcStorage := testHelper.CreateEtcStorage(testHelper.currentUsername, tc.loginShellCmd)
			userDirStorage := testHelper.CreateUserDirStorage(fmt.Sprintf("/home/%s/.okctl", testHelper.currentUsername))
			userHomeDirStorage := storage.NewEphemeralStorage()
			tmpStorage := storage.NewEphemeralStorage()
			environment := "myenv"

			opts := commandlineprompter.CommandLinePromptOpts{
				OsEnvVars:          osEnvVars,
				EtcStorage:         etcStorage,
				UserDirStorage:     userDirStorage.storage,
				UserHomeDirStorage: userHomeDirStorage,
				TmpStorage:         tmpStorage,
				Environment:        environment,
				CurrentUsername:    testHelper.currentUsername,
			}
			venv, err := virtualenv.CreateVirtualEnvironment(opts)
			assert.Nil(t, err)

			expectedVenv := testHelper.toSlice(tc.expectedOsEnvVars)
			assert.Equal(t, expectedVenv, venv.Environ())
		})
	}
	//
	//
	//t.Run("should return merged path", func(t *testing.T) {
	//	testHelper := NewTestHelper(t)
	//
	//	osEnvVars := make(map[string]string)
	//	osEnvVars["PATH"] = "/usr/local/go/bin:/usr/local/sbin:/home/johndoe/go/bin"
	//	osEnvVars["LS_COLORS"] = "rs=0:di=01;34:ln=01:*.tar=01;31"
	//	osEnvVars["OKCTL_PS1"] = "(some ps1 so it's simple to do the test assertion)"
	//	osEnvVars["OKCTL_SHELL"] = "/bin/fish"
	//
	//	currentUsername := "mickeymouse"
	//	etcStorage := testHelper.CreateEtcStorage(currentUsername, "/bin/bash")
	//	userDirStorage := testHelper.CreateUserDirStorage(currentUsername)
	//	userHomeDirStorage := storage.NewEphemeralStorage()
	//	tmpStorage := storage.NewEphemeralStorage()
	//	environment := "myenv"
	//
	//	opts := commandlineprompter.CommandLinePromptOpts{
	//		OsEnvVars:          osEnvVars,
	//		EtcStorage:         etcStorage,
	//		UserDirStorage:     userDirStorage.storage,
	//		UserHomeDirStorage: userHomeDirStorage,
	//		TmpStorage:         tmpStorage,
	//		Environment:        environment,
	//		CurrentUsername:    currentUsername,
	//	}
	//	venv, err := virtualenv.CreateVirtualEnvironment(opts)
	//	assert.Nil(t, err)
	//
	//	expectedPath := fmt.Sprintf("%s:%s", userDirStorage.basePath + "/venv", osEnvVars["PATH"])
	//	expectedVenv := []string{
	//		"LS_COLORS=rs=0:di=01;34:ln=01:*.tar=01;31",
	//		"OKCTL_PS1=(my_special_ps1)",
	//		"OKCTL_SHELL=/bin/fish",
	//		"PATH=" + expectedPath, // TODO: Må få trigga merged PATH, evt ny test
	//	}
	//
	//	assert.Equal(t, expectedVenv, venv.Environ())
	//})

	t.Run("should not change environment variables when OKCTL_NO_PS1 is set", func(t *testing.T) {
		testHelper := NewTestHelper(t)

		osEnvVars := make(map[string]string)
		osEnvVars["PATH"] = "/usr/local/go/bin:/usr/local/sbin:/home/johndoe/go/bin"
		osEnvVars["LS_COLORS"] = "rs=0:di=01;34:ln=01:*.tar=01;31"
		osEnvVars["OKCTL_NO_PS1"] = "true"

		currentUsername := "mickeymouse"
		etcStorage := testHelper.CreateEtcStorage(currentUsername, "/bin/bash")
		userDirStorage := storage.NewEphemeralStorage()
		userHomeDirStorage := storage.NewEphemeralStorage()
		tmpStorage := storage.NewEphemeralStorage()
		environment := "myenv"

		opts := commandlineprompter.CommandLinePromptOpts{
			OsEnvVars:          osEnvVars,
			EtcStorage:         etcStorage,
			UserDirStorage:     userDirStorage,
			UserHomeDirStorage: userHomeDirStorage,
			TmpStorage:         tmpStorage,
			Environment:        environment,
			CurrentUsername:    currentUsername,
		}
		venv, err := virtualenv.CreateVirtualEnvironment(opts)
		assert.Nil(t, err)

		expectedVenv := []string{
			"LS_COLORS=" + osEnvVars["LS_COLORS"],
			"OKCTL_NO_PS1=" + osEnvVars["OKCTL_NO_PS1"],
			"PATH=" + osEnvVars["PATH"],
		}

		assert.Equal(t, expectedVenv, venv.Environ())
	})

}

/*
func TestGetVirtualEnvironment(t *testing.T) {
	t.Run("should return expected environment variables", func(t *testing.T) {
		opts := virtualenv.VirtualEnvironmentOpts{
			UserDataDir:            mock.DefaultUserDataDir,
			Debug:                  false,
			Region:                 mock.DefaultRegion,
			AWSAccountID:           mock.DefaultAWSAccountID,
			Environment:            mock.DefaultEnv,
			Repository:             mock.DefaultRepositoryName,
			ClusterName:            mock.DefaultClusterName,
			KubectlBinaryDir:       mock.DefaultKubectlBinaryDir,
			AwsIamAuthenticatorDir: mock.DefaultAwsIamAuthenticatorDir,
		}

		userPath := "/usr/local/go/bin:/usr/local/sbin:/home/johndoe/go/bin"
		osEnvVars := []string{
			"SOME_VAR=A",
			"PATH=" + userPath,
			"LS_COLORS=rs=0:di=01;34:ln=01:*.tar=01;31",
		}

		venv, err := virtualenv.GetVirtualEnvironment(&opts, osEnvVars)
		assert.Nil(t, err)

		expectedVenv := []string{
			"AWS_CONFIG_FILE=" + path.Join(opts.UserDataDir, config.DefaultCredentialsDirName, opts.ClusterName, config.DefaultClusterAwsConfig),
			"AWS_PROFILE=default",
			"AWS_SHARED_CREDENTIALS_FILE=" + path.Join(opts.UserDataDir, config.DefaultCredentialsDirName, opts.ClusterName, config.DefaultClusterAwsCredentials),
			"HELM_CACHE_HOME=" + path.Join(opts.UserDataDir, config.DefaultHelmBaseDir),
			"HELM_CONFIG_DATA_HOME=" + path.Join(opts.UserDataDir, config.DefaultHelmBaseDir),
			"HELM_CONFIG_HOME=" + path.Join(opts.UserDataDir, config.DefaultHelmBaseDir),
			"HELM_DEBUG=" + strconv.FormatBool(opts.Debug),
			"HELM_PLUGINS=" + path.Join(opts.UserDataDir, config.DefaultHelmBaseDir, config.DefaultHelmPluginsDirectory),
			"HELM_REGISTRY_CONFIG=" + path.Join(opts.UserDataDir, config.DefaultHelmBaseDir, config.DefaultHelmRegistryConfig),
			"HELM_REPOSITORY_CACHE=" + path.Join(opts.UserDataDir, config.DefaultHelmBaseDir, config.DefaultHelmRepositoryCache),
			"HELM_REPOSITORY_CONFIG=" + path.Join(opts.UserDataDir, config.DefaultHelmBaseDir, config.DefaultHelmRepositoryConfig),
			"KUBECONFIG=" + path.Join(opts.UserDataDir, config.DefaultCredentialsDirName, opts.ClusterName, config.DefaultClusterKubeConfig),
			"LS_COLORS=rs=0:di=01;34:ln=01:*.tar=01;31",
			fmt.Sprintf("PATH=%s:%s:%s",
				mock.DefaultKubectlBinaryDir,
				mock.DefaultAwsIamAuthenticatorDir,
				userPath,
			),
			"SOME_VAR=A",
		}

		assert.Equal(t, expectedVenv, venv.Environ())
	})
}

func TestCreatePs1File(t *testing.T) {
	t.Run("should return the directory of the PS1 file", func(t *testing.T) {
		store := storage.NewEphemeralStorage()
		ps1Dir, err := virtualenv.createPs1ExecutableIfNotExists(store)

		assert.Nil(t, err)

		assert.Equal(t, "venv", ps1Dir)
	})

	t.Run("should create a file 'venv_ps1' if it doesn't exist", func(t *testing.T) {
		store := storage.NewEphemeralStorage()
		ps1Dir, err := virtualenv.createPs1ExecutableIfNotExists(store)

		assert.Nil(t, err)

		ps1Path := path.Join(ps1Dir, "venv_ps1")
		exists, err := store.Exists(ps1Path)

		assert.Nil(t, err)
		assert.True(t, exists)
	})
}
*/
