package nameserver_test

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/oslokommune/okctl/pkg/github"

	"github.com/oslokommune/okctl/pkg/git"
	"github.com/oslokommune/okctl/pkg/git/nameserver"
	"github.com/stretchr/testify/assert"
)

const expect = `resource "aws_route53_record" "test-oslo-systems" {
	zone_id = var.top_level_id # Leave this be
	name = "test.oslo.systems."
	type = "NS"
	ttl = 300
	records = [
		"ns1.something.com.",
		"ns2.something.com.",
		"ns3.something.com.",
		"ns4.something.com.",
	]
}`

func TestToHCL(t *testing.T) {
	got := nameserver.CreateRecord("test.oslo.systems.", []string{
		"ns1.something.com.",
		"ns2.something.com.",
		"ns3.something.com.",
		"ns4.something.com.",
	})

	assert.Equal(t, []byte(expect), got)
}

const layout = "2006-Jan-02"

// nolint: funlen
func TestNameserverDelegatorCreate(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("skipping test in CI environment")
	}

	tm, err := time.Parse(layout, "2014-Feb-04")
	assert.NoError(t, err)
	rand.Seed(tm.UnixNano())

	runner := git.New()

	assign := nameserver.NewDelegator(false).Assign("test.oslo.systems.", []string{
		"ns1.something.com.",
		"ns2.something.com.",
		"ns3.something.com.",
		"ns4.something.com.",
	}, nameserver.DefaultWorkingDir)
	assign.Stager = git.RepositoryStagerClone(git.RepositoryURL(github.DefaultOrg, "okctl-repository-test"))
	assign.Branch = "test-branch"

	result, err := runner.UpdateRepository(assign)
	assert.NoError(t, err)
	assert.Equal(t, &git.Result{
		Branch:             "test-branch",
		ModifiedRepository: true,
	}, result)

	result, err = runner.UpdateRepository(assign)
	assert.NoError(t, err)
	assert.Equal(t, &git.Result{
		Branch:             "test-branch",
		ModifiedRepository: false,
	}, result)

	revoke := nameserver.NewDelegator(false).Revoke("test.oslo.systems.", nameserver.DefaultWorkingDir)
	revoke.Stager = git.RepositoryStagerClone(git.RepositoryURL(github.DefaultOrg, "okctl-repository-test"))
	revoke.Branch = "test-branch"

	result, err = runner.UpdateRepository(revoke)
	assert.NoError(t, err)
	assert.Equal(t, &git.Result{
		Branch:             "test-branch",
		ModifiedRepository: true,
	}, result)

	result, err = runner.UpdateRepository(revoke)
	assert.NoError(t, err)
	assert.Equal(t, &git.Result{
		Branch:             "test-branch",
		ModifiedRepository: false,
	}, result)
}
