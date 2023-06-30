package gitlab

import (
	"fmt"
	"github.com/evnsio/decision/pkg/git"
	"github.com/xanzy/go-gitlab"
)

func (p *Provider) GetFolders() ([]string, error) {
	recurse := true
	nodes, _, err := p.client.Repositories.ListTree(repositoryId(), &gitlab.ListTreeOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
		},
		Ref:       &git.CommitHeadBranch,
		Recursive: &recurse,
	})

	if err != nil {
		fmt.Printf("Error getting tree: %s", err)
		return nil, err
	}

	folders := make([]string, 0)
	for _, node := range nodes {
		if node.Type == "tree" {
			folders = append(folders, node.Path)
		}
	}

	return folders, nil
}
