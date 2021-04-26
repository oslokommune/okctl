package git_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/go-git/go-git/v5/storage/memory"

	"github.com/go-git/go-billy/v5/memfs"

	gitpkg "github.com/go-git/go-git/v5"
	"github.com/oslokommune/okctl/pkg/git"

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
	got := git.CreateRecord("test.oslo.systems.", []string{
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
	tm, err := time.Parse(layout, "2014-Feb-04")
	assert.NoError(t, err)
	rand.Seed(tm.UnixNano())

	fs := memfs.New()
	mem := memory.NewStorage()

	testCases := []struct {
		name        string
		delegator   *git.NameserverDelegator
		fqdn        string
		nameservers []string
		expect      interface{}
		expectErr   bool
	}{
		{
			name: "Should work",
			delegator: git.NewNameserverDelegator(
				false,
				"",
				git.RepositoryStagerInit("initial commit", "README.md", "# initial commit", mem),
				fs,
			),
			fqdn: "test.oslo.systems.",
			nameservers: []string{
				"ns1.something.com.",
				"ns2.something.com.",
				"ns3.something.com.",
				"ns4.something.com.",
			},
			expect:    "✅ Add nameserver record for test.oslo.systems.",
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.delegator.CreateDelegation(tc.fqdn, tc.nameservers, true)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)

				r, err := gitpkg.Open(mem, fs)
				assert.NoError(t, err)

				last, err := r.Head()
				assert.NoError(t, err)

				commit, err := r.CommitObject(last.Hash())
				assert.NoError(t, err)

				assert.Equal(t, tc.expect, commit.Message)
				assert.Equal(t, &git.Result{
					Branch:             "test-oslo-systems-SpqY",
					ModifiedRepository: true,
				}, got)
			}

			fs = memfs.New()
			mem = memory.NewStorage()
		})
	}
}

// nolint: funlen
func TestNameserverDelegatorRevoke(t *testing.T) {
	tm, err := time.Parse(layout, "2014-Feb-04")
	assert.NoError(t, err)
	rand.Seed(tm.UnixNano())

	fs := memfs.New()
	mem := memory.NewStorage()

	testCases := []struct {
		name        string
		delegator   *git.NameserverDelegator
		fqdn        string
		nameservers []string
		expect      interface{}
		expectErr   bool
	}{
		{
			name: "Should work",
			delegator: git.NewNameserverDelegator(
				false,
				"",
				git.RepositoryStagerInit("something", "test-oslo-systems.tf", "something", mem),
				fs,
			),
			fqdn:      "test.oslo.systems.",
			expect:    "❌ Remove nameserver record for test.oslo.systems.",
			expectErr: false,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			got, err := tc.delegator.RevokeDelegation(tc.fqdn, true)

			if tc.expectErr {
				assert.Error(t, err)
				assert.Equal(t, tc.expect, err.Error())
			} else {
				assert.NoError(t, err)

				r, err := gitpkg.Open(mem, fs)
				assert.NoError(t, err)

				last, err := r.Head()
				assert.NoError(t, err)

				commit, err := r.CommitObject(last.Hash())
				assert.NoError(t, err)

				assert.Equal(t, tc.expect, commit.Message)
				assert.Equal(t, &git.Result{
					Branch:             "test-oslo-systems-SpqY",
					ModifiedRepository: true,
				}, got)
			}
		})

		fs = memfs.New()
		mem = memory.NewStorage()
	}
}
