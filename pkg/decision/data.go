package decision

var (
	Token       string
	CommitAsPRs bool
)

const (
	SlashCommand = "/decision"

	TitleBlockID = "title_block"
	TitleInputID = "title_input"

	CategoryBlockID  = "category_block"
	CategorySelectID = "category_select"

	ContextBlockID = "context_block"
	ContextInputID = "context_input"

	DecisionBlockID = "decision_block"
	DecisionInputID = "decision_input"

	ConsequencesBlockID = "consequences_block"
	ConsequencesInputID = "consequences_input"

	LogDecisionCallbackID = "log_decision"
)

type Decision struct {
	Title        string
	SlackHandle  string
	TeamID       string
	UserID       string
	Category     string
	Date         string
	Context      string
	Decision     string
	Consequences string
}

var decisionTemplate = `# {{.Title}}

Author: [@{{.SlackHandle}}](slack://user?team={{.TeamID}}&id={{.UserID}})

Category: ` + "`{{.Category}}`" + `

Date: {{.Date}}

## Context

{{.Context}}

## Decision

{{.Decision}}

## Consequences

{{.Consequences}}
`
