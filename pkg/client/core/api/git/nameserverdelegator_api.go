// Package git knows how to do git operations
package git

import (
	"fmt"
	"path"
	"strings"
	"time"

	"github.com/oslokommune/okctl/pkg/config/constant"

	"github.com/go-git/go-billy/v5/memfs"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/oslokommune/okctl/pkg/client"
	"github.com/oslokommune/okctl/pkg/github"
)

type repository struct {
	Endpoint       string
	Organization   string
	RepositoryName string
}

func (r repository) String() string {
	return fmt.Sprintf("%s:%s/%s.git", r.Endpoint, r.Organization, r.RepositoryName)
}

// NameserverDelegationRequest represents a single request for a nameserver record
type NameserverDelegationRequest struct {
	githubAPI github.Githuber

	Record *client.NameserverRecord
}

// Submit creates a pull request to the IAC repo with the NS record
func (n *NameserverDelegationRequest) Submit() (err error) {
	err = n.githubAPI.CreateNSRecordPullRequest(n.Name())
	if err != nil {
		return fmt.Errorf("error creating pull request: %w", err)
	}

	return nil
}

// Name returns a value that can be used to identify a request
func (n *NameserverDelegationRequest) Name() string {
	return strings.Split(n.Record.FQDN, ".")[0]
}

// IsSubmitted checks if theres already a PR for the relevant NS record
func (n *NameserverDelegationRequest) IsSubmitted() bool {
	hasExisting, err := n.githubAPI.HasExistingNSRecordPullRequest(n.Name())
	if err != nil {
		return false
	}

	return hasExisting
}

// NewNameserverDelegationRequest instantiates a DelegationRequest struct
func NewNameserverDelegationRequest(gh github.Githuber, record *client.NameserverRecord) *NameserverDelegationRequest {
	return &NameserverDelegationRequest{
		githubAPI: gh,
		Record:    record,
	}
}

/*
NameserverDelegator handles delegation of NS records. This involves the following:
- Creating a terraform nameserver record definition for the requested FQDN in the top level domain
- Committing the definition file to a branch in the top level IAC repository
- Creating a pull request for the top level domain admins to review
*/
type NameserverDelegator struct {
	githubAPI            github.Githuber
	repositoryDetails    *repository
	currentWorkDirectory string

	repository *git.Repository
	fs         billy.Filesystem
}

// CreateNameserverDelegationRequest knows how to request a nameserver record delegation from the top level domain
func (n *NameserverDelegator) CreateNameserverDelegationRequest(fqdn string, nameservers []string) (request *NameserverDelegationRequest, err error) {
	record := &client.NameserverRecord{
		FQDN:        fqdn,
		Nameservers: nameservers,
	}

	err = record.Validate()
	if err != nil {
		return nil, fmt.Errorf("malformed record: %w", err)
	}

	request = NewNameserverDelegationRequest(n.githubAPI, record)

	err = n.setup(request.Name())
	if err != nil {
		return nil, fmt.Errorf("error preparing repo: %w", err)
	}

	if n.hasExistingRequest(request) {
		return request, nil
	}

	file, err := n.fs.Create(fmt.Sprintf("%s.tf", request.Name()))
	if err != nil {
		return nil, fmt.Errorf("error creating temporary file: %w", err)
	}

	_, err = file.Write(toHCL(record))
	if err != nil {
		return nil, fmt.Errorf("error writing to temporary file: %w", err)
	}

	err = n.commit(request.Name())
	if err != nil {
		return nil, fmt.Errorf("error committing change: %w", err)
	}

	err = n.push(request.Name())
	if err != nil {
		return nil, fmt.Errorf("error pushing change: %w", err)
	}

	return request, nil
}

func (n *NameserverDelegator) hasExistingRequest(request *NameserverDelegationRequest) bool {
	iter, err := n.repository.References()
	if err != nil {
		return false
	}

	defer iter.Close()

	desiredName := request.Name()

	for {
		current, err := iter.Next()
		if err != nil {
			break
		}

		currentRef := current.Name().String()
		desiredRef := fmt.Sprintf("refs/remotes/origin/%s", desiredName)

		if currentRef == desiredRef {
			return true
		}
	}

	return false
}

func (n *NameserverDelegator) setup(name string) (err error) {
	err = n.clone()
	if err != nil {
		return fmt.Errorf("error cloning repository: %w", err)
	}

	n.fs, err = n.fs.Chroot(n.currentWorkDirectory)
	if err != nil {
		return fmt.Errorf("error chrooting into relevant directory: %w", err)
	}

	err = n.checkout(name)
	if err != nil {
		return fmt.Errorf("error checking out new branch: %w", err)
	}

	return nil
}

func (n *NameserverDelegator) clone() (err error) {
	n.repository, err = git.Clone(memory.NewStorage(), n.fs, &git.CloneOptions{
		URL:   n.repositoryDetails.String(),
		Depth: 1,
	})
	if err != nil {
		return fmt.Errorf("error fetching nameserver repository: %w", err)
	}

	return nil
}

func (n *NameserverDelegator) checkout(branchName string) (err error) {
	workTree, err := n.repository.Worktree()
	if err != nil {
		return fmt.Errorf("error getting work tree: %w", err)
	}

	err = workTree.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(branchName),
		Create: true,
		Force:  false,
	})

	if err != nil {
		return fmt.Errorf("error creating and switching branch: %w", err)
	}

	return nil
}

func (n *NameserverDelegator) commit(name string) error {
	workTree, err := n.repository.Worktree()
	if err != nil {
		return fmt.Errorf("error getting work tree: %w", err)
	}

	filePath := path.Join(n.currentWorkDirectory, fmt.Sprintf("%s.tf", name))

	_, err = workTree.Add(filePath)
	if err != nil {
		return fmt.Errorf("error adding file: %w", err)
	}

	_, err = workTree.Commit(
		fmt.Sprintf("✅ Add nameserver record for %s", name),
		&git.CommitOptions{
			Author: &object.Signature{
				Name:  "okctl automation",
				Email: "",
				When:  time.Now(),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("error committing nameserver record: %w", err)
	}

	return nil
}

func (n *NameserverDelegator) push(branchName string) (err error) {
	err = n.repository.Push(&git.PushOptions{
		RefSpecs: []config.RefSpec{
			config.RefSpec(fmt.Sprintf("HEAD:refs/heads/%s", branchName)),
		},
	})
	if err != nil {
		return fmt.Errorf("error pushing to remote: %w", err)
	}

	return nil
}

// NewNameserverDelegator initializes a new NameserverDelegator
func NewNameserverDelegator(githubAPI github.Githuber) *NameserverDelegator {
	return &NameserverDelegator{
		currentWorkDirectory: "infrastructure/production/dns/oslo.systems/auto_applied_subdomains",
		fs:                   memfs.New(),
		githubAPI:            githubAPI,
		repositoryDetails: &repository{
			Endpoint:       "git@github.com",
			Organization:   github.DefaultOrg,
			RepositoryName: github.DefaultAWSInfrastructureRepository,
		},
	}
}

func toHCL(record *client.NameserverRecord) []byte {
	name := strings.Split(record.FQDN, ".")[0]

	result := []string{
		fmt.Sprintf("resource \"aws_route53_record\" \"%s\" {", name),
		"\tzone_id = var.top_level_id # Leave this be",
		fmt.Sprintf("\tname = \"%s\"", record.FQDN),
		"\ttype = \"NS\"",
		fmt.Sprintf("\tttl = %d", constant.DefaultNameserverRecordTTL),
		"\trecords = [",
	}

	for _, entry := range record.Nameservers {
		result = append(result, fmt.Sprintf("\t\t\"%s.\",", entry))
	}

	result = append(result, "\t]\n}\n")

	return []byte(strings.Join(result, "\n"))
}
