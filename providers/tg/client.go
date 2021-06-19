package tg

import (
	"fmt"
	"net/url"

	"dialogflowbot/providers/common"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Client struct {
	api        *tg.BotAPI
	webhookURL *url.URL
}

func NewClient(accessToken string, webhookURL *url.URL) (*Client, error) {
	var client *Client
	var err error

	client.api, err = tg.NewBotAPI(accessToken)
	if err != nil {
		return nil, err
	}

	err = client.setWebhookIfNeed(webhookURL.String())
	if err != nil {
		return nil, err
	}

	client.webhookURL = webhookURL
	return client, err
}

func (c *Client) setWebhookIfNeed(webhookURL string) error {
	if webhookURL == "" {
		return fmt.Errorf("need set webhook url")
	}

	whInfo, err := c.api.GetWebhookInfo()
	if err != nil {
		return err
	}

	if whInfo.URL == webhookURL {
		return nil
	}

	if whInfo.IsSet() {
		_, err := c.api.RemoveWebhook()
		if err != nil {
			return err
		}
	}

	_, err = c.api.SetWebhook(tg.NewWebhook(webhookURL))
	return err
}

func (c *Client) GetMessages() <-chan common.Message {
	ch := make(chan common.Message)

	go func() {
		for update := range c.api.ListenForWebhook(c.webhookURL.Path) {
			if update.Message == nil || update.Message.From == nil {
				continue
			}

			ch <- common.Message{
				Provider: common.TGProvider,
				ChatID:   int(update.Message.Chat.ID),
				Content:  update.Message.Text,
			}
		}
	}()

	return ch
}

func (c *Client) SendMessage(msg common.Message) error {
	message := tg.NewMessage(int64(msg.ChatID), msg.Content)

	if _, err := c.api.Send(message); err != nil {
		return err
	}

	return nil
}
