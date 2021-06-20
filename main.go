package main

import (
	"log"
	"net/url"

	"github.com/delgus/dialogflowbot/bot"
	"github.com/delgus/dialogflowbot/providers/common"
	"github.com/delgus/dialogflowbot/providers/tg"
	"github.com/delgus/dialogflowbot/providers/vk"
	"github.com/delgus/dialogflowbot/wsapp/chat"
	"github.com/kelseyhightower/envconfig"
)

type config struct {
	CredentialsJSON string `envconfig:"CREDENTIALS_JSON"`
	ProjectID       string `envconfig:"PROJECT_ID"`

	TGAccessToken string `envconfig:"TG_ACCESS_TOKEN"`
	TGWebhook     string `envconfig:"TG_WEBHOOK"`

	VKAccessToken string `envconfig:"VK_ACCESS_TOKEN"`
	VKWebhook     string `envconfig:"VK_WEBHOOK"`
	VKConfirmKey  string `envconfig:"VK_CONFIRM_KEY"`

	WSURL string `envconfig:"WS_URL"`
}

func main() {
	// configurations
	var cfg config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal("can't load config", err)
	}

	// websocket app
	chat.Run()

	// providers:

	// tg
	tgWebhook, err := url.Parse(cfg.TGWebhook)
	if err != nil {
		log.Fatal(err)
	}

	tgProvider, err := tg.NewClient(cfg.TGAccessToken, tgWebhook)
	if err != nil {
		log.Fatal(err)
	}

	// vk
	vkWebhook, err := url.Parse(cfg.VKWebhook)
	if err != nil {
		log.Fatal(err)
	}

	vkProvider := vk.NewClient(cfg.VKAccessToken, vkWebhook, cfg.VKConfirmKey)
	if err != nil {
		log.Fatal(err)
	}

	// dialogflow bot
	dfBot, err := bot.NewBot(cfg.CredentialsJSON, cfg.ProjectID, map[common.ProviderType]common.Provider{
		common.TGProvider: tgProvider,
		common.VKProvider: vkProvider,
	})
	if err != nil {
		log.Fatal("can't start dialog flow bot", err)
	}

	dfBot.Run()
}
