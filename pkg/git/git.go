package git

import (
	"fmt"
	"regexp"

	"github.com/go-git/go-git/v5"
)

// GithubRepoFullName attempts to extract the github repo name, returns
// and empty string if it doesn't succeed
func GithubRepoFullName(org, gitPath string) (string, error) {
	g, err := git.PlainOpen(gitPath)
	if err != nil {
		return "", fmt.Errorf("failed to open git repository: %w", err)
	}

	remotes, err := g.Remotes()
	if err != nil {
		return "", fmt.Errorf("failed to get remotes: %w", err)
	}

	re := regexp.MustCompile(fmt.Sprintf(`(?U)(?:.*github\.com:)(?P<repo>%s/.*)(?:\.git)`, org))

	for _, remote := range remotes {
		matches := re.FindStringSubmatch(remote.String())
		if len(matches) == 0 {
			continue
		}

		names := re.SubexpNames()
		for i, name := range names {
			if name == "repo" {
				return matches[i], nil
			}
		}
	}

	return "", nil
}
