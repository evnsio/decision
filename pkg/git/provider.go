package git

type Provider interface {
	// RaisePullRequest will automatically create a commit, create a branch, and open a pull request, and return the URL to the PR.
	RaisePullRequest(branch string, commitMessage string, path string, content []byte) (string, error)

	//CreateCommit creates a commit with the given content
	CreateCommit(commitMessage string, path string, content []byte) (string, error)
}
