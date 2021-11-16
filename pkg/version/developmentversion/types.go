package developmentversion

import (
	"context"
	"github.com/google/go-github/v32/github"
)

// RepositoryRelease shadows github.RepositoryRelease
type RepositoryRelease = github.RepositoryRelease

// ListReleasesFn returns GitHub releases
type ListReleasesFn func(ctx context.Context, owner string, repo string) ([]*RepositoryRelease, error)
