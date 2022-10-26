package gitlab

import "github.com/xanzy/go-gitlab"

type Provider struct {
	client *gitlab.Client
}

func (p *Provider) RaisePullRequest(branch string, commitMessage string, path string, content []byte) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (p *Provider) CreateCommit(commitMessage string, path string, content []byte) (string, error) {
	//TODO implement me
	panic("implement me")
}

func NewProvider(accessToken string) *Provider {
	c, err := gitlab.NewClient(accessToken)
	if err != nil {
		panic(err)
	}

	return &Provider{
		client: c,
	}
}
