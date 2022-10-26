package github

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

type Provider struct {
	client *github.Client
}

func NewProvider(accessToken string) *Provider {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: accessToken})
	tc := oauth2.NewClient(context.Background(), ts)

	return &Provider{
		client: github.NewClient(tc),
	}
}
