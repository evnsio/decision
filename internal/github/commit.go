package github

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func CreateCommit(commitMessage string, path string, content []byte) (fileURL *string) {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: Token})
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	opts := &github.RepositoryContentFileOptions{
		Message:   github.String(commitMessage),
		Content:   content,
		Branch:    github.String(CommitBranch),
		Committer: &github.CommitAuthor{Name: github.String(AuthorName), Email: github.String(AuthorEmail)},
	}
	res, _, err := client.Repositories.CreateFile(ctx, SourceOwner, SourceRepo, path, opts)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return res.Content.HTMLURL
}
