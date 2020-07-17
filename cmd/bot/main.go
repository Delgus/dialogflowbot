package main

import (
	"fmt"
	"net/http"

	"github.com/delgus/dialogflow-tg-bot/internal/bot"
	easybot "github.com/delgus/easy-bot"
	"github.com/delgus/easy-bot/clients/tg"
	tghook "github.com/delgus/tg-logrus-hook"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

type config struct {
	CredentialsJSON  string `envconfig:"CREDENTIALS_JSON"`
	ProjectID        string `envconfig:"PROJECT_ID"`
	TGAccessToken    string `envconfig:"TG_ACCESS_TOKEN"`
	TGWebhook        string `envconfig:"TG_WEBHOOK"`
	Host             string `envconfig:"HOST"`
	Port             int    `envconfig:"PORT"`
	LogLevel         string `envconfig:"LOG_LEVEL"`
	LogTGChatID      int64  `envconfig:"LOG_TG_CHAT_ID"`
	LogTGAccessToken string `envconfig:"LOG_TG_ACCESS_TOKEN"`
}

func main() {
	// configurations
	var cfg config
	err := envconfig.Process("", &cfg)
	if err != nil {
		logrus.Fatalf("can't load config: %v", err)
	}
	logrus.SetLevel(parseLogLevel(cfg.LogLevel))

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

	hook, err := tghook.NewHook(cfg.LogTGAccessToken, cfg.LogTGChatID, logrus.AllLevels)
	if err != nil {
		logrus.Errorf(`can not create tg hook for logging: %v`, err)
	} else {
		logrus.Info(`create tg hook for logging`)
		logrus.StandardLogger().AddHook(hook)
		logrus.Info(`add hook!!!`)
	}

	app := &easybot.App{
		Notifier: tgNotifier,
		Bot:      dfBot,
		Listener: tgListener,
		Logger:   logrus.StandardLogger(),
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

func parseLogLevel(level string) logrus.Level {
	switch level {
	case "trace":
		return logrus.TraceLevel
	case "debug":
		return logrus.DebugLevel
	case "info":
		return logrus.InfoLevel
	case "error":
		return logrus.ErrorLevel
	case "fatal":
		return logrus.FatalLevel
	case "panic":
		return logrus.PanicLevel
	default:
		return logrus.InfoLevel
	}
}
