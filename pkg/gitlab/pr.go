package gitlab

import (
	"github.com/evnsio/decision/pkg/decision"
	"github.com/evnsio/decision/pkg/git"
	"github.com/xanzy/go-gitlab"
)

func (p *Provider) RaisePullRequest(branch string, commitMessage string, path string, content []byte) (string, error) {
	_, err := p.createCommitOnBranch(commitMessage, path, string(content), branch)

	if err != nil {
		return "", err
	}

	removeBranch := true
	squash := true
	description := decision.PullRequestBody(commitMessage)
	mr, _, err := p.client.MergeRequests.CreateMergeRequest(repositoryId(), &gitlab.CreateMergeRequestOptions{
		Title:              &commitMessage,
		Description:        &description,
		SourceBranch:       &branch,
		TargetBranch:       &git.CommitHeadBranch,
		RemoveSourceBranch: &removeBranch,
		Squash:             &squash,
	})

	if err != nil {
		return "", err
	}

	return mr.WebURL, nil
}
