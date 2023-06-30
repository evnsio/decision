package github

import (
	"context"
	"github.com/evnsio/decision/pkg/git"

	"github.com/google/go-github/github"
)

func (p *Provider) CreateCommit(commitMessage string, path string, content []byte) (string, error) {
	return p.createFileOnBranch(commitMessage, path, git.CommitHeadBranch, content)
}

// createFileOnBranch will create a new file on an **existing** branch. The function will fail if the branch does not exist.
func (p *Provider) createFileOnBranch(commitMessage, path, branch string, content []byte) (string, error) {
	commitData := &github.RepositoryContentFileOptions{
		Message:   github.String(commitMessage),
		Content:   content,
		Branch:    github.String(branch),
		Committer: &github.CommitAuthor{Name: github.String(git.AuthorName), Email: github.String(git.AuthorEmail)},
	}

	res, _, err := p.client.Repositories.CreateFile(context.Background(), git.SourceOwner, git.SourceRepo, path, commitData)

	if err != nil {
		return "", err
	}

	return *res.Content.HTMLURL, nil
}
