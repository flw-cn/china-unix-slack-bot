package main

import (
	"context"
	"fmt"

	slackbot "github.com/flw-cn/go-slackbot"
	"github.com/flw-cn/playground/docker"
	"github.com/flw-cn/slack"
	log "github.com/sirupsen/logrus"
)

func PlayGo(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	codeInfo := slackbot.NamedCapturesFromContext(ctx)

	lang := codeInfo.Get("lang")
	code := codeInfo.Get("code")

	output, err := docker.PlayCode(lang, code)
	if err != nil {
		log.Printf("Error: %s\n%s", err, output)
		return
	}

	if output != "" {
		m := fmt.Sprintf("咦？有人发代码，让我试着运行一下！\n%s", output)
		bot.Reply(evt, m, slackbot.WithTyping)
	}
}

func PlayGoFile(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	lang := evt.File.Filetype
	file, cleaner, err := slackDownloadFile(bot.Client, evt.File.ID)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return
	}

	defer cleaner()

	output, err := docker.PlayFile(lang, file)
	if err != nil {
		log.Printf("Error: %s\n%s", err, output)
		return
	}

	if output != "" {
		m := fmt.Sprintf("你这段代码能运行吗？让我试着运行一下！\n%s", output)
		bot.ReplyInThread(evt, m, slackbot.WithTyping)
	}
}
