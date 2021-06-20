package vk

import (
	"context"
	"net/http"
	"net/url"

	"dialogflowbot/providers/common"

	"github.com/SevereCloud/vksdk/v2/api"
	"github.com/SevereCloud/vksdk/v2/api/params"
	"github.com/SevereCloud/vksdk/v2/callback"
	"github.com/SevereCloud/vksdk/v2/events"
)

type Client struct {
	api        *api.VK
	webhookURL *url.URL
	confirmKey string
}

func NewClient(accessToken string, webhookURL *url.URL, confirmKey string) *Client {
	return &Client{
		api:        api.NewVK(accessToken),
		webhookURL: webhookURL,
		confirmKey: confirmKey,
	}
}

func (c *Client) GetMessages() <-chan common.Message {
	ch := make(chan common.Message)
	cb := callback.NewCallback()
	cb.ConfirmationKey = c.confirmKey

	cb.MessageNew(func(ctx context.Context, obj events.MessageNewObject) {
		ch <- common.Message{
			Provider: common.VKProvider,
			ChatID:   obj.Message.FromID,
			Content:  obj.Message.Text,
		}
	})

	http.HandleFunc(c.webhookURL.Path, cb.HandleFunc)

	return ch
}

func (c *Client) SendMessage(msg common.Message) error {
	b := params.NewMessagesSendBuilder()
	b.PeerID(msg.ChatID)
	b.RandomID(0)
	b.DontParseLinks(false)
	b.Message(msg.Content)

	_, err := c.api.MessagesSend(b.Params)
	return err
}
