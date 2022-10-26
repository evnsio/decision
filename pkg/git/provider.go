package git

import (
	"github.com/evnsio/decision/pkg/github"
	"github.com/evnsio/decision/pkg/gitlab"
)

var provider Provider

type Provider interface {
	// RaisePullRequest will automatically create a commit, create a branch, and open a pull request, and return the URL to the PR.
	RaisePullRequest(branch string, commitMessage string, path string, content []byte) (string, error)

	// CreateCommit creates a commit with the given content
	CreateCommit(commitMessage string, path string, content []byte) (string, error)

	// GetFolders returns all available folders in the configured repository
	GetFolders() ([]string, error)
}

func GetProvider() Provider {
	if provider != nil {
		return provider
	}

	switch ProviderType {
	case "github":
		provider = github.NewProvider(Token)
		return provider
	case "gitlab":
		provider = gitlab.NewProvider(Token)
	}

	return provider
}
