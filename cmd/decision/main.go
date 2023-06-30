package main

import (
	"encoding/json"
	"fmt"
	"github.com/evnsio/decision/pkg/git"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/namsral/flag"

	"github.com/evnsio/decision/pkg/decision"
	"github.com/slack-go/slack"
)

var (
	signingSecret string
)

func handleSlash(w http.ResponseWriter, r *http.Request) {

	verifier, err := slack.NewSecretsVerifier(r.Header, signingSecret)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	r.Body = io.NopCloser(io.TeeReader(r.Body, &verifier))
	s, err := slack.SlashCommandParse(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err = verifier.Ensure(); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	switch s.Command {
	case decision.SlashCommand:
		decision.OpenDecisionModal(s.TriggerID, s.ChannelID)
	default:
		fmt.Printf("%v -- %v", s.Command, s.Text)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func handleActions(w http.ResponseWriter, r *http.Request) {
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		fmt.Printf("Could not parse action response JSON: %v", err)
	}

	if payload.Type == slack.InteractionTypeViewSubmission {
		if payload.View.CallbackID == decision.LogDecisionCallbackID {
			go decision.HandleModalSubmission(&payload)
		}
	}
}

func handleOptions(w http.ResponseWriter, r *http.Request) {
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		fmt.Printf("Could not parse action response JSON: %v", err)
	}

	if payload.Type == slack.InteractionTypeBlockSuggestion {
		switch payload.ActionID {
		case decision.CategorySelectID:
			categoryOptions := decision.GetCategoryOptions(&payload.Value)

			js, err := json.Marshal(categoryOptions)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		default:
			fmt.Printf("No handler found for action_id: %v", payload.ActionID)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

}

func redact(thing string) string {
	if len(thing) < 3 {
		return thing
	} else {
		return thing[0:3] + strings.Repeat(".", len(thing)-3)
	}
}

func parseFlags() {
	var config string
	flag.StringVar(&config, "config", "", "Path to config file")
	flag.BoolVar(&decision.CommitAsPRs, "commit-as-prs", false, "Commit decisions as Pull Requests")
	flag.StringVar(&decision.Token, "slack-token", "", "Your Slack API token starting xoxb-...")
	flag.StringVar(&signingSecret, "slack-signing-secret", "", "Your Slack signing secret")
	flag.StringVar(&git.ProviderType, "provider", "github", "Your git provider, github or gitlab")
	flag.StringVar(&git.Token, "provider-token", "", "Your access token for your Git provider")
	flag.StringVar(&git.SourceOwner, "source-owner", "", "The owner / organisation of the repo where decisions will be committed")
	flag.StringVar(&git.SourceRepo, "source-repo", "", "The repo where decisions will be committed")
	flag.StringVar(&git.CommitHeadBranch, "branch", "main", "The branch where decisions will be committed")
	flag.StringVar(&git.AuthorName, "commit-author", "", "The author name to use for commits")
	flag.StringVar(&git.AuthorEmail, "commit-email", "", "The author email to use for commits")
	flag.Parse()

	if decision.Token == "" {
		fmt.Fprintln(os.Stderr, "missing required argument -slack-token")
		os.Exit(2)
	} else {
		fmt.Printf("Slack Token: %v\n", redact(decision.Token))
	}

	if !strings.HasPrefix(decision.Token, "xoxb") {
		fmt.Fprintln(os.Stderr, "-slack-token should be a bot token starting 'xoxb-'")
		os.Exit(2)
	}

	if signingSecret == "" {
		fmt.Fprintln(os.Stderr, "missing required argument -slack-signing-secret")
		os.Exit(2)
	} else {
		fmt.Printf("Slack Signing Secret: %v\n", redact(signingSecret))
	}

	if git.Token == "" {
		fmt.Fprintln(os.Stderr, "missing required argument -provider-token")
		os.Exit(2)
	} else {
		fmt.Printf("Github Token: %v\n", redact(git.Token))
	}

	if git.SourceOwner == "" {
		fmt.Fprintln(os.Stderr, "missing required argument -source-owner")
		os.Exit(2)
	} else {
		fmt.Printf("Source Owner: %v\n", redact(git.SourceOwner))
	}

	if git.SourceRepo == "" {
		fmt.Fprintln(os.Stderr, "missing required argument -source-repo")
		os.Exit(2)
	} else {
		fmt.Printf("Source Repo: %v\n", redact(git.SourceRepo))
	}

	if git.AuthorName == "" {
		fmt.Fprintln(os.Stderr, "missing required argument -commit-author")
		os.Exit(2)
	} else {
		fmt.Printf("Commit Author: %v\n", redact(git.AuthorName))
	}

	if git.AuthorEmail == "" {
		fmt.Fprintln(os.Stderr, "missing required argument -commit-email")
		os.Exit(2)
	} else {
		fmt.Printf("Commit Email: %v\n", redact(git.AuthorEmail))
	}
}

func main() {
	parseFlags()

	http.HandleFunc("/action", handleActions)
	http.HandleFunc("/slash", handleSlash)
	http.HandleFunc("/options", handleOptions)
	fmt.Println("[INFO] Server listening")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Printf("%v\n", err)
	}
}
