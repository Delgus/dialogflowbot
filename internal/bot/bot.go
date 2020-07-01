package bot

import (
	"context"
	"fmt"
	"io"

	dg "cloud.google.com/go/dialogflow/apiv2"
	easybot "github.com/delgus/easy-bot"
	"google.golang.org/api/option"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type DialogFlowBot struct {
	client    *dg.SessionsClient
	projectID string
}

// NewBot return new bot
func NewBot(credentialJSON, projectID string) (*DialogFlowBot, error) {
	// https://dialogflow.com/docs/reference/v2-auth-setup
	client, err := dg.NewSessionsClient(
		context.Background(),
		option.WithCredentialsJSON([]byte(credentialJSON)))
	if err != nil {
		return nil, err
	}
	return &DialogFlowBot{
		client:    client,
		projectID: projectID,
	}, nil
}

// Command implement interface Bot
func (dfb *DialogFlowBot) Command(command easybot.Command) (string, error) {
	// get dialogflow answer
	resp, err := dfb.client.DetectIntent(context.Background(), &dialogflow.DetectIntentRequest{
		Session: fmt.Sprintf("projects/%s/agent/sessions/%d'", dfb.projectID, command.Args.UserID),
		QueryInput: &dialogflow.QueryInput{
			Input: &dialogflow.QueryInput_Text{Text: &dialogflow.TextInput{
				Text:         command.Name,
				LanguageCode: "ru",
			}},
		},
	})
	if err != nil {
		return "", err
	}
	result := resp.GetQueryResult().FulfillmentText
	if result == "" {
		return "", io.EOF
	}

	return result, nil
}
