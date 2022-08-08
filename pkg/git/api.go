package git

import (
	"errors"
	"fmt"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/format/index"
	"github.com/go-git/go-git/v5/storage/memory"
)

// DeleteRemoteFile knows how to delete a file in a remote Git repository
func DeleteRemoteFile(repositoryURL string, path string, commitMessage string) error {
	repo, err := git.Clone(memory.NewStorage(), memfs.New(), &git.CloneOptions{
		URL:   repositoryURL,
		Depth: 1,
	})
	if err != nil {
		return fmt.Errorf("cloning repository: %w", err)
	}

	tree, _ := repo.Worktree()

	_, err = tree.Remove(path)
	if err != nil {
		if errors.Is(err, index.ErrEntryNotFound) {
			return nil
		}

		return fmt.Errorf("removing file: %w", err)
	}

	_, err = tree.Commit(commitMessage, &git.CommitOptions{})
	if err != nil {
		return fmt.Errorf("committing changes: %w", err)
	}

	err = repo.Push(&git.PushOptions{})
	if err != nil {
		return fmt.Errorf("pushing to repository: %w", err)
	}

	return nil
}
