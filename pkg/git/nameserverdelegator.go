// Package git knows how to do git operations
package git

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/gosimple/slug"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/oslokommune/okctl/pkg/github"
)

const (
	// DefaultWorkingDir represents the default working directory
	DefaultWorkingDir = "infrastructure/production/dns/oslo.systems/auto_applied_subdomains"
	// DefaultAutoWorkingDir represents the default auto working directory
	DefaultAutoWorkingDir = "infrastructure/production/dns/oslo.systems/auto_applied_subdomains/auto"
	randomChars           = 4
)

// DefaultRepositoryURL represents the default repository URL
func DefaultRepositoryURL() string {
	return fmt.Sprintf(
		"%s:%s/%s.git",
		"git@github.com",
		github.DefaultOrg,
		github.DefaultAWSInfrastructureRepository,
	)
}

// RepositoryStagerFn defines the required function for
// staging a git repository
type RepositoryStagerFn func(fs billy.Filesystem) (*git.Repository, error)

// NameserverDelegator contains the state required for
// working with DNS zone delegations
type NameserverDelegator struct {
	repositoryStagerFn RepositoryStagerFn
	workingDir         string
	fs                 billy.Filesystem
}

// Result contains data about the created content
type Result struct {
	Branch             string
	ModifiedRepository bool
}

// RevokeDelegation creates a zone revocation branch
func (n *NameserverDelegator) RevokeDelegation(fqdn string, skipPush bool) (*Result, error) {
	repository, err := n.repositoryStagerFn(n.fs)
	if err != nil {
		return nil, fmt.Errorf("preparing repo: %w", err)
	}

	branch := fmt.Sprintf("%s-%s", slug.Make(fqdn), randSeq(randomChars))
	fileName := fmt.Sprintf("%s.tf", slug.Make(fqdn))

	fs, err := n.fs.Chroot(n.workingDir)
	if err != nil {
		return nil, fmt.Errorf("chrooting into relevant directory: %w", err)
	}

	workTree, err := n.checkout(branch, repository)
	if err != nil {
		return nil, fmt.Errorf("checking out new branch: %w", err)
	}

	err = fs.Remove(fileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Result{
				ModifiedRepository: false,
			}, nil
		}

		return nil, fmt.Errorf("removing file: %w", err)
	}

	err = n.commitRevokeRecord(fqdn, workTree)
	if err != nil {
		return nil, fmt.Errorf("committing change: %w", err)
	}

	if !skipPush {
		err = n.push(branch, repository)
		if err != nil {
			return nil, fmt.Errorf("pushing change: %w", err)
		}
	}

	return &Result{
		Branch:             branch,
		ModifiedRepository: true,
	}, nil
}

// CreateDelegation knows how to request a nameserver record delegation from the top level domain
func (n *NameserverDelegator) CreateDelegation(fqdn string, nameservers []string, skipPush bool) (*Result, error) {
	repository, err := n.repositoryStagerFn(n.fs)
	if err != nil {
		return nil, fmt.Errorf("preparing repo: %w", err)
	}

	branch := fmt.Sprintf("%s-%s", slug.Make(fqdn), randSeq(randomChars))
	fileName := fmt.Sprintf("%s.tf", slug.Make(fqdn))

	fs, err := n.fs.Chroot(n.workingDir)
	if err != nil {
		return nil, fmt.Errorf("chrooting into relevant directory: %w", err)
	}

	workTree, err := n.checkout(branch, repository)
	if err != nil {
		return nil, fmt.Errorf("checking out new branch: %w", err)
	}

	file, err := fs.Create(fileName)
	if err != nil {
		return nil, fmt.Errorf("creating temporary file: %w", err)
	}

	record := CreateRecord(fqdn, nameservers)

	_, err = file.Write(record)
	if err != nil {
		return nil, fmt.Errorf("writing to temporary file: %w", err)
	}

	filePath := path.Join(n.workingDir, fileName)

	err = n.commitCreateRecord(fqdn, filePath, workTree)
	if err != nil {
		return nil, fmt.Errorf("committing change: %w", err)
	}

	if !skipPush {
		err = n.push(branch, repository)
		if err != nil {
			return nil, fmt.Errorf("pushing change: %w", err)
		}
	}

	return &Result{
		Branch:             branch,
		ModifiedRepository: true,
	}, nil
}

func (*NameserverDelegator) checkout(branchName string, r *git.Repository) (*git.Worktree, error) {
	workTree, err := r.Worktree()
	if err != nil {
		return nil, fmt.Errorf("getting work tree: %w", err)
	}

	err = workTree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
		Create: true,
		Keep:   true,
	})
	if err != nil {
		return nil, fmt.Errorf("checking out branch: %w", err)
	}

	return workTree, nil
}

func (*NameserverDelegator) commitCreateRecord(fqdn, filePath string, workTree *git.Worktree) error {
	_, err := workTree.Add(filePath)
	if err != nil {
		return fmt.Errorf("adding file: %w", err)
	}

	_, err = workTree.Commit(
		fmt.Sprintf("✅ Add nameserver record for %s", fqdn),
		&git.CommitOptions{
			Author: &object.Signature{
				Name:  "okctl automation",
				Email: "okctl@oslo.kommune.no",
				When:  time.Now(),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("committing nameserver record: %w", err)
	}

	return nil
}

func (*NameserverDelegator) commitRevokeRecord(fqdn string, workTree *git.Worktree) error {
	_, err := workTree.Commit(
		fmt.Sprintf("❌ Remove nameserver record for %s", fqdn),
		&git.CommitOptions{
			Author: &object.Signature{
				Name:  "okctl automation",
				Email: "okctl@oslo.kommune.no",
				When:  time.Now(),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("committing nameserver record: %w", err)
	}

	return nil
}

func (*NameserverDelegator) push(branchName string, repository *git.Repository) error {
	b := plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", branchName))

	err := repository.Push(&git.PushOptions{
		RemoteName: "origin",
		RefSpecs: []config.RefSpec{
			config.RefSpec(b + ":" + b),
		},
	})
	if err != nil {
		return fmt.Errorf("pushing to remote: %w", err)
	}

	return nil
}

// NewNameserverDelegator initializes a new NameserverDelegator
func NewNameserverDelegator(seedRandom bool, workingDir string, fn RepositoryStagerFn, fs billy.Filesystem) *NameserverDelegator {
	if seedRandom {
		rand.Seed(time.Now().UnixNano())
	}

	return &NameserverDelegator{
		repositoryStagerFn: fn,
		workingDir:         workingDir,
		fs:                 fs,
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

const createRecordTemplate = `resource "aws_route53_record" "%s" {
	zone_id = var.top_level_id # Leave this be
	name = "%s"
	type = "NS"
	ttl = 300
	records = [
%s
	]
}`

// CreateRecord returns the create record template
func CreateRecord(fqdn string, nameservers []string) []byte {
	for i, ns := range nameservers {
		nameservers[i] = fmt.Sprintf("\t\t\"%s\",", ns)
	}

	return []byte(
		fmt.Sprintf(
			createRecordTemplate,
			slug.Make(fqdn),
			fqdn,
			strings.Join(nameservers, "\n"),
		),
	)
}

// from: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go/22892986#22892986
// nolint: gochecknoglobals
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)

	for i := range b {
		// nolint: gosec
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}
