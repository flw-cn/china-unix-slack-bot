package main

import (
	"context"
	"strings"

	slackbot "github.com/flw-cn/go-slackbot"
	"github.com/flw-cn/slack"
	"github.com/pityonline/china-unix-slack-bot/service"
)

func QueryIP(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	id := bot.BotUserID()
	id = "<@" + id + ">"
	text := strings.Replace(evt.Text, id, "", -1)
	text = strings.Trim(text, " ")
	parts := strings.Fields(text)

	api := "http://freeapi.ipip.net/"
	var m string
	if len(parts) != 2 {
		m = "Usage: ip <ip address>"
	} else {
		ip := parts[1]
		m = service.IPQuery(api, ip)
	}
	bot.Reply(evt, m, slackbot.WithTyping)
}
