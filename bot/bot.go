package bot

import (
	"context"
	"fmt"
	"log"
	"net/http"

	dg "cloud.google.com/go/dialogflow/apiv2"
	"github.com/delgus/dialogflowbot/providers/common"
	"google.golang.org/api/option"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

const (
	internalErrorText = "Внутренняя ошибка. На нашей стороне ошибка, попробуйте позднее"
	noResult          = "Не понял тебя. Попробуй выразиться по другому"
)

type DialogFlowBot struct {
	client    *dg.SessionsClient
	providers map[common.ProviderType]common.Provider
	projectID string
}

// NewBot return new bot
func NewBot(credentialJSON, projectID string, providers map[common.ProviderType]common.Provider) (*DialogFlowBot, error) {
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
		providers: providers,
	}, nil
}

func (dfb *DialogFlowBot) Run() {
	commonCh := make(chan common.Message)

	for providerType := range dfb.providers {
		go func(providerType common.ProviderType) {
			for msg := range dfb.providers[providerType].GetMessages() {
				commonCh <- msg
			}
		}(providerType)
	}

	go func() {
		log.Fatal(http.ListenAndServe(":80", nil))
	}()

	for msg := range commonCh {
		dfb.answerJob(msg)
	}
}

func (dfb *DialogFlowBot) answerFromDialogflow(msg common.Message) common.Message {
	// get dialogflow answer
	resp, err := dfb.client.DetectIntent(context.Background(), &dialogflow.DetectIntentRequest{
		Session: fmt.Sprintf(
			"projects/%s/agent/sessions/%s%d",
			dfb.projectID, msg.Provider, msg.ChatID),
		QueryInput: &dialogflow.QueryInput{
			Input: &dialogflow.QueryInput_Text{Text: &dialogflow.TextInput{
				Text:         msg.Content,
				LanguageCode: "ru",
			}},
		},
	})
	if err != nil {
		log.Println(err)
		return common.Message{Provider: msg.Provider, ChatID: msg.ChatID, Content: internalErrorText}
	}
	result := resp.GetQueryResult().FulfillmentText
	if result == "" {
		return common.Message{Provider: msg.Provider, ChatID: msg.ChatID, Content: noResult}
	}

	return common.Message{Provider: msg.Provider, ChatID: msg.ChatID, Content: result}
}

func (dfb *DialogFlowBot) answerJob(msg common.Message) {
	response := dfb.answerFromDialogflow(msg)
	provider := dfb.providers[msg.Provider]

	if provider == nil {
		log.Fatal("can not find provider!", msg.Provider)
	}

	if err := provider.SendMessage(response); err != nil {
		log.Println(err)
	}
}
