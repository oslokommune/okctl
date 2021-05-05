// Package nameserver knows how to create git commits for zone delegation and revocation
package nameserver

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/go-git/go-billy/v5/memfs"

	"github.com/oslokommune/okctl/pkg/git"
	"github.com/oslokommune/okctl/pkg/github"

	"github.com/gosimple/slug"
)

const (
	// DefaultWorkingDir represents the default working directory
	DefaultWorkingDir = "infrastructure/production/dns/oslo.systems/auto_applied_subdomains"
	// DefaultAutoWorkingDir represents the default auto working directory
	DefaultAutoWorkingDir = "infrastructure/production/dns/oslo.systems/auto_applied_subdomains/auto"
)

const (
	randomChars = 4
)

// DefaultRepositoryURL represents the default repository URL
func DefaultRepositoryURL() string {
	return git.RepositoryURL(github.DefaultOrg, github.DefaultAWSInfrastructureRepository)
}

// BranchNameFromFQDN returns a random branch name based on the provided
// fqdn
func BranchNameFromFQDN(fqdn string) string {
	return fmt.Sprintf("%s-%s", slug.Make(fqdn), randSeq(randomChars))
}

// FileNameFromFQDN returns a filename based on the provided fqdn
func FileNameFromFQDN(fqdn string) string {
	return fmt.Sprintf("%s.tf", slug.Make(fqdn))
}

// Delegator contains required state
type Delegator struct{}

// Assign a delegation record
func (d *Delegator) Assign(fqdn string, nameservers []string, workingDir string) *git.ChangeSet {
	return &git.ChangeSet{
		Stager:        git.RepositoryStagerClone(DefaultRepositoryURL()),
		Branch:        BranchNameFromFQDN(fqdn),
		PushToRemote:  true,
		FileSystem:    memfs.New(),
		CommitMessage: fmt.Sprintf("✅ Add nameserver record for %s", fqdn),
		Actions: []git.ActionFn{
			git.AddFile(workingDir, FileNameFromFQDN(fqdn), CreateRecord(fqdn, nameservers)),
		},
	}
}

// Revoke a delegation record
func (d *Delegator) Revoke(fqdn string, workingDir string) *git.ChangeSet {
	return &git.ChangeSet{
		Stager:        git.RepositoryStagerClone(DefaultRepositoryURL()),
		Branch:        BranchNameFromFQDN(fqdn),
		PushToRemote:  true,
		FileSystem:    memfs.New(),
		CommitMessage: fmt.Sprintf("❌ Remove nameserver record for %s", fqdn),
		Actions: []git.ActionFn{
			git.RemoveFile(workingDir, FileNameFromFQDN(fqdn)),
		},
	}
}

// NewDelegator returns an initialised delegator
func NewDelegator(seedRandom bool) *Delegator {
	if seedRandom {
		rand.Seed(time.Now().UnixNano())
	}

	return &Delegator{}
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
