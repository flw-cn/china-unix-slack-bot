package main

import (
	"context"
	"os"

	"github.com/flw-cn/go-slackbot"
	"github.com/flw-cn/go-smartConfig"
	"github.com/flw-cn/slack"
	"github.com/pityonline/china-unix-slack-bot/service"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Debug bool   `flag:"d|false|debug mode, default to 'false'"`
	Token string `flag:"t||must provide your {SLACK_TOKEN} here"`
}

var config Config

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

func main() {

	smartConfig.LoadConfig("Slack Bot", "v0.2.0", &config)

	bot, _ := slackbot.NewBot(config.Token)

	if config.Debug {
		logger := log.New()
		logger.Formatter = &log.TextFormatter{}
		logger.Out = os.Stderr
		logger.SetLevel(log.DebugLevel)
		logger.Debug("Running in debug mode...")
		bot.SetLogger(logger)
		bot.Client.SetDebug(true)
	}

	toMe := bot.Messages(slackbot.DirectMessage, slackbot.Mention).Subrouter()
	toMe.Hear("(?i)(hi|hello).*").MessageHandler(Hello)
	toMe.Hear("(?i)(ping).*").MessageHandler(Ping)
	toMe.Hear("(?i)(ip) .*").MessageHandler(QueryIP)
	toMe.MessageHandler(Whatever)

	code := bot.Messages(slackbot.Message).Subrouter()
	code.Hear("(?s)```\\s*?(?P<lang>\\S*)\\n\\s*(?P<code>.+)\\s*```").MessageHandler(PlayGo)

	codeFile := bot.Messages(slackbot.FileShared).Subrouter()
	codeFile.MessageHandler(PlayGoFile)

	bot.Run(true, nil)
}

func Hello(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	m := service.Greet()
	bot.Reply(evt, m, slackbot.WithTyping)
}

func Ping(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	m := service.Ping()
	bot.Reply(evt, m, slackbot.WithTyping)
}
