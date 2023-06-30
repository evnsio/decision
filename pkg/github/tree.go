package github

import (
	"context"
	"fmt"
	"github.com/evnsio/decision/pkg/git"
)

func (p *Provider) GetFolders() ([]string, error) {
	tree, _, err := p.client.Git.GetTree(context.Background(), git.SourceOwner, git.SourceRepo, fmt.Sprintf("refs/heads/%s", git.CommitHeadBranch), false)
	if err != nil {
		fmt.Printf("Error getting tree: %s", err)
		return nil, nil
	}

	folders := make([]string, 0)
	for _, entry := range tree.Entries {
		if *entry.Type == "tree" {
			folders = append(folders, *entry.Path)
		}
	}

	return folders, nil
}
