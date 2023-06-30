package decision

import (
	"bytes"
	"fmt"
	"github.com/evnsio/decision/pkg/provider"
	"html/template"
	"strings"
	"sync"
	"time"

	"github.com/gosimple/slug"

	"github.com/slack-go/slack"
)

var (
	categoryOptions []*slack.OptionBlockObject
	categoryLock    sync.Mutex
)

func OpenDecisionModal(triggerID string, triggerChannel string) {
	titleLabel := slack.NewTextBlockObject(slack.PlainTextType, "Title", false, false)
	titlePlaceholderText := slack.NewTextBlockObject(slack.PlainTextType, "Give this decision a tl;dr title", false, false)
	titleInput := slack.NewPlainTextInputBlockElement(titlePlaceholderText, TitleInputID)
	titleInput.MaxLength = 60
	titleSection := slack.NewInputBlock(TitleBlockID, titleLabel, titleInput)

	categoryLabel := slack.NewTextBlockObject(slack.PlainTextType, "Category", false, false)
	categorySelect := slack.NewOptionsSelectBlockElement(slack.OptTypeExternal, nil, CategorySelectID)
	categorySelect.MinQueryLength = new(int)
	categorySection := slack.NewInputBlock(CategoryBlockID, categoryLabel, categorySelect)

	contextLabel := slack.NewTextBlockObject(slack.PlainTextType, "Context", false, false)
	contextPlaceholderText := slack.NewTextBlockObject(slack.PlainTextType, "Explain why this decision needs to be made. What forces are at play?", false, false)
	contextInput := slack.NewPlainTextInputBlockElement(contextPlaceholderText, ContextInputID)
	contextInput.Multiline = true
	contextSection := slack.NewInputBlock(ContextBlockID, contextLabel, contextInput)

	decisionLabel := slack.NewTextBlockObject(slack.PlainTextType, "Decision", false, false)
	decisionPlaceholderText := slack.NewTextBlockObject(slack.PlainTextType, "Document the decision you've made. Use active voice: \"We will...\"", false, false)
	decisionInput := slack.NewPlainTextInputBlockElement(decisionPlaceholderText, DecisionInputID)
	decisionInput.Multiline = true
	decisionSection := slack.NewInputBlock(DecisionBlockID, decisionLabel, decisionInput)

	consequencesLabel := slack.NewTextBlockObject(slack.PlainTextType, "Consequences", false, false)
	consequencesPlaceholderText := slack.NewTextBlockObject(slack.PlainTextType, "Describe the consequences, good and bad, after this decision has been made.", false, false)
	consequencesInput := slack.NewPlainTextInputBlockElement(consequencesPlaceholderText, ConsequencesInputID)
	consequencesInput.Multiline = true
	consequencesSection := slack.NewInputBlock(ConsequencesBlockID, consequencesLabel, consequencesInput)

	view := slack.ModalViewRequest{
		CallbackID:      LogDecisionCallbackID,
		Type:            slack.ViewType("modal"),
		Title:           slack.NewTextBlockObject(slack.PlainTextType, "Decision", false, false),
		Close:           slack.NewTextBlockObject(slack.PlainTextType, "Cancel", false, false),
		Submit:          slack.NewTextBlockObject(slack.PlainTextType, "Log Decision", false, false),
		PrivateMetadata: triggerChannel,
		Blocks: slack.Blocks{
			BlockSet: []slack.Block{
				titleSection,
				categorySection,
				contextSection,
				decisionSection,
				consequencesSection,
			},
		},
	}

	api := slack.New(Token)
	_, err := api.OpenView(triggerID, view)
	if err != nil {
		fmt.Printf("Error opening modal view: %v\n", err)
	}
}

func GetCategoryOptions(typeAheadValue *string) slack.OptionsResponse {
	// we only fetch the existing categories from github when the modal is first show
	// subsequent calls (sent as the user is typing) re-use the list from this fetch
	modalFirstOpened := typeAheadValue != nil && *typeAheadValue == ""
	if modalFirstOpened {
		categoryLock.Lock()
		defer categoryLock.Unlock()

		categoryOptions = make([]*slack.OptionBlockObject, 0)
		existingFolders, _ := provider.GetProvider().GetFolders()

		//existingFolders, _ := github.GetFolders()
		//existingFolders, _ := github.GetFolders()
		for _, folder := range existingFolders {
			categoryOptions = append(categoryOptions, slack.NewOptionBlockObject(
				strings.ToLower(folder),
				slack.NewTextBlockObject(slack.PlainTextType, folder, false, false)))
		}
	}

	// add the type ahead value as an option in the list
	var responseOptions = categoryOptions
	if typeAheadValue != nil && *typeAheadValue != "" {
		typeAheadOption := slack.NewOptionBlockObject(
			slug.Make(*typeAheadValue),
			slack.NewTextBlockObject(slack.PlainTextType, *typeAheadValue+" (Create new)", false, false))

		// only add it if it doesn't exist already
		if typeAheadOption != nil {
			typeAheadExists := false
			for _, category := range categoryOptions {
				if typeAheadOption.Value == category.Value {
					typeAheadExists = true
					break
				}
			}

			if !typeAheadExists {
				responseOptions = append([]*slack.OptionBlockObject{typeAheadOption}, categoryOptions...)
			}
		}
	}

	response := slack.OptionsResponse{
		Options: responseOptions,
	}

	return response
}

func HandleModalSubmission(payload *slack.InteractionCallback) {
	submissionValues := payload.View.State.Values

	sourceChannel := payload.View.PrivateMetadata

	title := submissionValues[TitleBlockID][TitleInputID].Value
	category := submissionValues[CategoryBlockID][CategorySelectID].SelectedOption.Value
	context := submissionValues[ContextBlockID][ContextInputID].Value
	decision := submissionValues[DecisionBlockID][DecisionInputID].Value
	consequences := submissionValues[ConsequencesBlockID][ConsequencesInputID].Value

	username := payload.User.Name
	if payload.User.Profile.DisplayName != "" {
		username = payload.User.Profile.DisplayName
	}

	decisionData := Decision{
		Title:        title,
		SlackHandle:  username,
		TeamID:       payload.Team.ID,
		UserID:       payload.User.ID,
		Category:     category,
		Date:         time.Now().Format("2006-01-02"),
		Context:      context,
		Decision:     decision,
		Consequences: consequences,
	}

	tmpl, err := template.New("decision").Parse(decisionTemplate)
	if err != nil {
		fmt.Printf("Failed to parse template: %v", err)
		return
	}

	var decisionBytes bytes.Buffer
	err = tmpl.Execute(&decisionBytes, decisionData)
	if err != nil {
		fmt.Printf("Failed to execute template: %v", err)
		return
	}

	dateNow := time.Now().Format("2006-01-02")
	fileName := category + "/" + dateNow + "-" + slug.Make(title) + ".md"
	commitMessage := title
	content := decisionBytes.Bytes()

	provider := provider.GetProvider()

	if CommitAsPRs {
		prURL, err := provider.RaisePullRequest(slug.Make(title), commitMessage, fileName, content)
		if err != nil {
			return
		}

		message := "✅ A pull request for \"" + title + "\" has been created <" + prURL + "|here>."
		sendDecisionLinkToUser(message, title, prURL, sourceChannel, payload.User.ID)
	} else {
		decisionURL, err := provider.CreateCommit(commitMessage, fileName, content)

		if err != nil {
			return
		}

		message := "✅ Your decision \"" + title + "\" has been committed <" + decisionURL + "|here>."
		sendDecisionLinkToUser(message, title, decisionURL, sourceChannel, payload.User.ID)
	}
}

func sendDecisionLinkToUser(message string, title string, fileURL string, channel string, user string) {
	// Return an ephemeral message to the user
	msgOption := slack.MsgOptionText(message, false)
	api := slack.New(Token)
	_, err := api.PostEphemeral(channel, user, msgOption)
	if err != nil {
		fmt.Printf("Failed to send message: %v (%v, %v)\n", err, channel, user)
		return
	}
}
