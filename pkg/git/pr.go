package git

import "fmt"

func PullRequestBody(subject string) string {
	return fmt.Sprintf("Logging decision for \"%s\"", subject)
}
