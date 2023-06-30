package gitlab

import (
	"fmt"
	"github.com/evnsio/decision/pkg/git"
	"github.com/xanzy/go-gitlab"
)

func (p *Provider) CreateCommit(commitMessage string, path string, content []byte) (string, error) {
	return p.createCommitOnBranch(commitMessage, path, string(content), git.CommitHeadBranch)
}

func (p *Provider) createCommitOnBranch(commitMessage, path, content, branch string) (string, error) {
	createAction := gitlab.FileCreate
	commit, _, err := p.client.Commits.CreateCommit(
		repositoryId(),
		&gitlab.CreateCommitOptions{
			Branch:        &branch,
			StartBranch:   &git.CommitHeadBranch,
			CommitMessage: &commitMessage,
			Actions: []*gitlab.CommitActionOptions{
				{
					Action:   &createAction,
					FilePath: &path,
					Content:  &content,
				},
			},
			AuthorEmail: &git.AuthorEmail,
			AuthorName:  &git.AuthorName,
		},
	)

	if err != nil {
		fmt.Printf("Error creating commit: %s", err)
		return "", err
	}

	return commit.WebURL, nil
}

func repositoryId() string {
	return fmt.Sprintf("%s/%s", git.SourceOwner, git.SourceRepo)
}
