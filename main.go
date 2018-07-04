package main

import (
	"log"
	"os"

	"github.com/flw-cn/go-smartConfig"
	"github.com/flw-cn/slack-bot/bot"
	"github.com/flw-cn/slack-bot/event"
	"github.com/flw-cn/slack-bot/plugin/backend/fortune"
	"github.com/flw-cn/slack-bot/plugin/backend/greeter"
	"github.com/flw-cn/slack-bot/plugin/backend/ipQuerier"
	"github.com/flw-cn/slack-bot/plugin/backend/ping"
	"github.com/flw-cn/slack-bot/plugin/backend/playground"
	"github.com/flw-cn/slack-bot/plugin/backend/tuling"
	"github.com/flw-cn/slack-bot/plugin/frontend/slack"
)

type Config struct {
	Debug    bool `flag:"d|false|debug mode, default to 'false'"`
	Bot      bot.Config
	Frontend struct {
		Slack slack.Config
	}
	Backend struct {
		Greeter   greeter.Config
		IPQuerier ipQuerier.Config
		Play      playground.Config
		Fortune   fortune.Config
		Tuling    tuling.Config
		Ping      ping.Config
	}
}

func main() {
	var config Config
	smartConfig.LoadConfig("Bot", "v0.3.0", &config)

	logger := log.New(os.Stderr, "BOT ", log.LstdFlags)

	if config.Debug {
		logger.Printf("Running in debug mode...")
	}

	err := startBot(logger, config)
	if err != nil {
		logger.Printf("Error: %v", err)
		return
	}

	select {}
}

func startBot(logger *log.Logger, config Config) error {
	bot := bot.New(config.Bot)
	if config.Debug {
		bot.SetDebug(true)
	}

	slack := slack.New(config.Frontend.Slack)
	err := bot.LoadFrontend(slack)
	if err != nil {
		return err
	}

	greeter := greeter.New(config.Backend.Greeter)
	ipQuerier := ipQuerier.New(config.Backend.IPQuerier)
	play := playground.New(config.Backend.Play)
	fortune := fortune.New(config.Backend.Fortune)
	tuling := tuling.New(config.Backend.Tuling)
	ping := ping.New(config.Backend.Ping)
	err = bot.LoadBackend(greeter, ipQuerier, play, fortune, tuling, ping)
	if err != nil {
		return err
	}

	bot.SetLogger(logger)

	err = bot.Init()
	if err != nil {
		return err
	}

	err = bot.Start()
	if err != nil {
		return err
	}

	toMe := bot.On(event.EvDirectMention, event.EvMentionedMe, event.EvDirectMessage).Subrouter()
	hook := toMe.Hear(`(?i)^(hi|hello)$`).Hook()
	bot.Mount(hook, greeter)

	hook = toMe.Hear(`(?i)^\s*ping\s*$`).Hook()
	bot.Mount(hook, ping)

	hook = toMe.Hear(`(?i)^ip\s+(?P<IP>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})$`).Hook()
	bot.Mount(hook, ipQuerier)

	code := bot.On(event.EvMessage).Subrouter()
	hook = code.Hear("(?s)```\\s*?\\n?(?P<lang>\\w*)\\n\\s*(?P<code>.+)\\s*```").Hook()
	bot.Mount(hook, play)

	hook = toMe.Hear(`(?P<words>(唐?诗)|(宋?词))`).Hook()
	bot.Mount(hook, fortune)

	hook = toMe.Hear("").Hook()
	bot.Mount(hook, tuling)

	return nil
}
