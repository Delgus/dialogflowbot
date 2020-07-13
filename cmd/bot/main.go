package main

import (
	"fmt"
	"net/http"

	"github.com/delgus/dialogflow-tg-bot/internal/bot"
	easybot "github.com/delgus/easy-bot"
	"github.com/delgus/easy-bot/clients/tg"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
)

type DialogFlowConfig struct {
	CredentialsJSON string `envconfig:"CREDENTIALS_JSON"`
	ProjectID       string `envconfig:"PROJECT_ID"`
}

type TGConfig struct {
	TGAccessToken string `envconfig:"TG_ACCESS_TOKEN"`
	TGWebhook     string `envconfig:"TG_WEBHOOK"`
}

type ServerConfig struct {
	Host string `envconfig:"HOST"`
	Port int    `envconfig:"PORT"`
}

type TGLoggerConfig struct {
	LogTGChatID      int64  `envconfig:"LOG_TG_CHAT_ID"`
	LogTGAccessToken string `envconfig:"LOG_TG_ACCESS_TOKEN_ID"`
}

type config struct {
	DialogFlowConfig
	TGConfig
	ServerConfig
	TGLoggerConfig
}

func main() {
	// configurations
	var cfg config
	err := envconfig.Process("", &cfg)
	if err != nil {
		logrus.Fatalf("can't load config: %v", err)
	}

	// dialogflow bot
	dfBot, err := bot.NewBot(cfg.CredentialsJSON, cfg.ProjectID)
	if err != nil {
		logrus.Fatalf("can't start dialog flow bot: %v", err)
	}

	// common options for apps
	opts := easybot.Options{
		EOFText:               `Не понимаю тебя`,
		InternalErrorText:     `Внутренняя ошибка. На нашей стороне ошибка, попробуйте позднее`,
		HelpText:              `Просто общайтесь со мной`,
		NotCorrectCommandText: "Неверная команда! \n",
	}

	// tg notifier
	tgNotifier, err := tg.NewNotifier(cfg.TGAccessToken)
	if err != nil {
		logrus.Fatalf("can't start notifier for bot: %v", err)
	}

	// tg listener
	tgListener, err := tg.NewListener(cfg.TGAccessToken, cfg.TGWebhook)
	if err != nil {
		logrus.Fatalf("can't start listener for bot: %v", err)
	}

	appLogger := logrus.StandardLogger()
	appLogger.SetLevel(logrus.DebugLevel)

	logrus.Debug(cfg.LogTGAccessToken, cfg.LogTGChatID)
	client, err := tgbotapi.NewBotAPI(cfg.LogTGAccessToken)
	if err != nil {
		logrus.Fatal(err)
	}
	msg := tgbotapi.MessageConfig{}
	msg.ChatID = cfg.LogTGChatID
	msg.Text = "Work!!!"
	res, err := client.Send(msg)
	if err != nil {
		logrus.Fatal(err, res)
	}
	app := &easybot.App{
		Notifier: tgNotifier,
		Bot:      dfBot,
		Listener: tgListener,
		Logger:   appLogger,
		Options:  opts,
	}

	go func() {
		if err := app.Run("/"); err != nil {
			logrus.Fatal("can not start tg bot app")
		}
	}()

	logrus.Info("start server for tg and app")
	addr := fmt.Sprintf(`%s:%d`, cfg.Host, cfg.Port)
	logrus.Fatal(http.ListenAndServe(addr, nil))
}
