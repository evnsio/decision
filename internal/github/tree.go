package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func GetFolders() ([]*string, error) {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: Token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	tree, _, err := client.Git.GetTree(ctx, SourceOwner, SourceRepo, "refs/heads/master", false)
	if err != nil {
		fmt.Printf("Error getting tree: %s", err)
		return nil, err
	}

	folders := make([]*string, 0)
	for _, entry := range tree.Entries {
		if *entry.Type == "tree" {
			folders = append(folders, entry.Path)
		}
	}

	return folders, nil
}
