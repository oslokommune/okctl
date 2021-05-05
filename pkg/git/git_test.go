package git_test

import (
	"os"
	"testing"

	"github.com/oslokommune/okctl/pkg/github"

	"github.com/go-git/go-billy/v5/memfs"

	"github.com/oslokommune/okctl/pkg/git"

	"github.com/stretchr/testify/assert"
)

func TestWithExternalRepository(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("skipping test in CI environment")
	}

	runner := git.New()

	changeSet := &git.ChangeSet{
		Stager:        git.RepositoryStagerClone(git.RepositoryURL(github.DefaultOrg, "okctl-repository-test")),
		Branch:        "my-branch",
		PushToRemote:  true,
		FileSystem:    memfs.New(),
		CommitMessage: "adding some data",
		Actions: []git.ActionFn{
			git.AddFile("", "my.file", []byte("some nice data")),
		},
	}

	_, err := runner.UpdateRepository(changeSet)
	assert.NoError(t, err)

	changeSet.CommitMessage = "removing some data"
	changeSet.Actions = []git.ActionFn{
		git.RemoveFile("", "my.file"),
	}
	changeSet.FileSystem = memfs.New()

	_, err = runner.UpdateRepository(changeSet)
	assert.NoError(t, err)
}
