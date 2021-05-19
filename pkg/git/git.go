// Package git knows how to do git operations
package git

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// RepositoryStagerFn defines the required function for
// staging a git repository
type RepositoryStagerFn func(fs billy.Filesystem) (*git.Repository, error)

// Git contains the state required for
// working with a git repository
type Git struct{}

// Result contains data about the created content
type Result struct {
	Branch             string
	ModifiedRepository bool
}

// ActionFn represents a function that makes changes
type ActionFn func(worktree *git.Worktree) error

// ChangeSet contains the required inputs for updating
// a git repository
type ChangeSet struct {
	Stager        RepositoryStagerFn
	Branch        string
	PushToRemote  bool
	FileSystem    billy.Filesystem
	CommitMessage string
	Actions       []ActionFn
}

// UpdateRepository applies the change set
//nolint:funlen,gocyclo
func (n *Git) UpdateRepository(c *ChangeSet) (*Result, error) {
	r, err := c.Stager(c.FileSystem)
	if err != nil {
		return nil, fmt.Errorf("staging repository: %w", err)
	}

	workTree, err := r.Worktree()
	if err != nil {
		return nil, fmt.Errorf("getting work tree: %w", err)
	}

	err = workTree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(c.Branch),
		Create: false,
	})
	if err != nil {
		err = workTree.Checkout(&git.CheckoutOptions{
			Branch: plumbing.NewBranchReferenceName(c.Branch),
			Create: true,
		})
		if err != nil {
			return nil, fmt.Errorf("checking out branch: %w", err)
		}
	}

	err = workTree.Pull(&git.PullOptions{
		RemoteName:    "origin",
		ReferenceName: plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", c.Branch)),
		SingleBranch:  true,
		Depth:         1,
	})
	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) || errors.Is(err, plumbing.ErrReferenceNotFound) || errors.Is(err, plumbing.ErrObjectNotFound) {
		} else {
			return nil, fmt.Errorf("pulling branch: %w", err)
		}
	}

	for _, action := range c.Actions {
		err = action(workTree)
		if err != nil {
			return nil, fmt.Errorf("running action: %w", err)
		}
	}

	status, err := workTree.Status()
	if err != nil {
		return nil, fmt.Errorf("checking status: %w", err)
	}

	if status.IsClean() {
		return &Result{
			Branch:             c.Branch,
			ModifiedRepository: false,
		}, nil
	}

	_, err = workTree.Commit(
		c.CommitMessage,
		&git.CommitOptions{
			Author: &object.Signature{
				Name:  "okctl automation",
				Email: "okctl@oslo.kommune.no",
				When:  time.Now(),
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("committing nameserver record: %w", err)
	}

	if c.PushToRemote {
		b := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", c.Branch))

		err := r.Push(&git.PushOptions{
			RemoteName: "origin",
			RefSpecs: []config.RefSpec{
				config.RefSpec(b + ":" + b),
			},
		})
		if err != nil {
			return nil, fmt.Errorf("pushing to remote: %w", err)
		}
	}

	return &Result{
		Branch:             c.Branch,
		ModifiedRepository: true,
	}, nil
}

// New initializes a new Git
func New() *Git {
	return &Git{}
}

// RepositoryURL builds url based on org and repo
func RepositoryURL(org, repo string) string {
	return fmt.Sprintf(
		"%s:%s/%s.git",
		"git@github.com",
		org,
		repo,
	)
}

// RemoveFile is a helper action for removing a file from a repository
func RemoveFile(workingDir, fileName string) ActionFn {
	return func(worktree *git.Worktree) error {
		filePath := worktree.Filesystem.Join(workingDir, fileName)

		err := worktree.Filesystem.Remove(filePath)
		if err != nil {
			if errors.Is(err, os.ErrNotExist) {
				return nil
			}

			return fmt.Errorf("removing file: %w", err)
		}

		_, err = worktree.Add(filePath)
		if err != nil {
			return fmt.Errorf("adding file: %w", err)
		}

		status, err := worktree.Status()
		if err != nil {
			return fmt.Errorf("getting status: %w", err)
		}

		if status.File(filePath).Staging != git.Deleted {
			return fmt.Errorf("file: %s, not staged as deleted", filePath)
		}

		return nil
	}
}

// AddFile is a helper action for adding a file to a repository
func AddFile(workingDir, fileName string, content []byte) ActionFn {
	return func(worktree *git.Worktree) error {
		filePath := worktree.Filesystem.Join(workingDir, fileName)

		file, err := worktree.Filesystem.Create(filePath)
		if err != nil {
			return fmt.Errorf("creating file: %w", err)
		}

		defer func(file billy.File) {
			err = file.Close()
		}(file)

		_, err = file.Write(content)
		if err != nil {
			return fmt.Errorf("writing to file: %w", err)
		}

		_, err = worktree.Add(filePath)
		if err != nil {
			return fmt.Errorf("adding file: %w", err)
		}

		return nil
	}
}

// RepositoryStagerClone knows how to clone a git repository from
// a given URL
func RepositoryStagerClone(repoURL string) RepositoryStagerFn {
	return func(fs billy.Filesystem) (*git.Repository, error) {
		repository, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
			URL:   repoURL,
			Depth: 1,
		})
		if err != nil {
			return nil, fmt.Errorf("cloning repository: %w", err)
		}

		return repository, nil
	}
}

// RepositoryStagerInit initialises an empty repository
func RepositoryStagerInit(msg, fileName, content string, mem *memory.Storage) RepositoryStagerFn {
	return func(fs billy.Filesystem) (*git.Repository, error) {
		repository, err := git.Init(mem, fs)
		if err != nil {
			return nil, fmt.Errorf("initialising repository: %w", err)
		}

		f, err := fs.Create(fileName)
		if err != nil {
			return nil, err
		}

		_, err = f.Write([]byte(content))
		if err != nil {
			return nil, err
		}

		err = f.Close()
		if err != nil {
			return nil, err
		}

		w, err := repository.Worktree()
		if err != nil {
			return nil, err
		}

		_, err = w.Commit(msg, &git.CommitOptions{
			All: true,
			Author: &object.Signature{
				Name:  "okctl",
				Email: "okctl@oslo.kommune.no",
				When:  time.Now(),
			},
		})
		if err != nil {
			return nil, err
		}

		return repository, nil
	}
}
