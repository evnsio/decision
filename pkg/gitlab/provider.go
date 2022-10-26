package gitlab

import "github.com/xanzy/go-gitlab"

type Provider struct {
	client *gitlab.Client
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
