package github

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

var client *github.Client
var ctx = context.Background()

// getRef returns the commit branch reference object if it exists or creates it
// from the base branch before returning it.
func getRef(commitBranch string) (ref *github.Reference, err error) {
	ctx := context.Background()
	if ref, _, err = client.Git.GetRef(ctx, SourceOwner, SourceRepo, "refs/heads/"+commitBranch); err == nil {
		return ref, nil
	}

	var baseRef *github.Reference
	if baseRef, _, err = client.Git.GetRef(ctx, SourceOwner, SourceRepo, "refs/heads/"+CommitBranch); err != nil {
		return nil, err
	}
	newRef := &github.Reference{Ref: github.String("refs/heads/" + commitBranch), Object: &github.GitObject{SHA: baseRef.Object.SHA}}
	ref, _, err = client.Git.CreateRef(ctx, SourceOwner, SourceRepo, newRef)
	return ref, err
}

// getTree generates the tree to commit based on the given files and the commit
// of the ref you got in getRef.
func getTree(path string, content []byte, ref *github.Reference) (tree *github.Tree, err error) {
	entries := []*github.TreeEntry{}

	entries = append(entries, &github.TreeEntry{
		Path:    github.String(path),
		Type:    github.String("blob"),
		Content: github.String(string(content)),
		Mode:    github.String("100644"),
	})

	tree, _, err = client.Git.CreateTree(ctx, SourceOwner, SourceRepo, *ref.Object.SHA, entries)

	return tree, err
}

// createCommit creates the commit in the given reference using the given tree.
func createCommit(commitMessage string, ref *github.Reference, tree *github.Tree) (err error) {
	// Get the parent commit to attach the commit to.
	parent, _, err := client.Repositories.GetCommit(ctx, SourceOwner, SourceRepo, *ref.Object.SHA)
	if err != nil {
		return err
	}
	// This is not always populated, but is needed.
	parent.Commit.SHA = parent.SHA

	// Create the commit using the tree.
	date := time.Now()
	author := &github.CommitAuthor{Date: &date, Name: &AuthorName, Email: &AuthorEmail}
	commit := &github.Commit{Author: author, Message: &commitMessage, Tree: tree, Parents: []*github.Commit{parent.Commit}}
	newCommit, _, err := client.Git.CreateCommit(ctx, SourceOwner, SourceRepo, commit)
	if err != nil {
		return err
	}

	// Attach the commit to the master branch.
	ref.Object.SHA = newCommit.SHA
	_, _, err = client.Git.UpdateRef(ctx, SourceOwner, SourceRepo, ref, false)
	return err
}

func createPR(prSubject string, branch string) (url *string, err error) {
	prBody := "Logging decision for \"" + prSubject + "\""

	newPR := &github.NewPullRequest{
		Title:               &prSubject,
		Head:                &branch,
		Base:                &CommitBranch,
		Body:                &prBody,
		MaintainerCanModify: github.Bool(true),
	}

	pr, _, err := client.PullRequests.Create(ctx, SourceOwner, SourceRepo, newPR)
	if err != nil {
		return nil, err
	}

	prURL := pr.GetHTMLURL()
	fmt.Printf("PR created: %s\n", prURL)

	return &prURL, nil
}

func RaisePullRequest(branch string, commitMessage string, path string, content []byte) *string {
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: Token})
	tc := oauth2.NewClient(ctx, ts)
	client = github.NewClient(tc)

	ref, err := getRef(branch)
	if err != nil {
		log.Printf("Unable to get/create the commit reference: %s\n", err)
		return nil
	}
	if ref == nil {
		log.Printf("No error where returned but the reference is nil")
		return nil
	}

	tree, err := getTree(path, content, ref)
	if err != nil {
		log.Printf("Unable to create the tree based on the provided files: %s\n", err)
		return nil
	}

	if err := createCommit(commitMessage, ref, tree); err != nil {
		log.Printf("Unable to create the commit: %s\n", err)
		return nil
	}

	prURL, err := createPR(commitMessage, branch)
	if err != nil {
		log.Printf("Error while creating the pull request: %s", err)
		return nil
	}

	return prURL
}
