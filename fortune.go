package main

import (
	"context"
	"math/rand"

	fortune "github.com/flw-cn/go-fortune"
	slackbot "github.com/flw-cn/go-slackbot"
	"github.com/flw-cn/slack"
	log "github.com/sirupsen/logrus"
)

func Fortune(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	var err error

	r := rand.Intn(100)
	m := "听不懂你在说什么，不如我给你念首诗吧！"
	o := ""
	if r < 10 {
		m = m + "\n" + "你想听唐诗还是宋词？"
	} else if r > 95 {
		m = m + "\n" + "中国话太难背了，咱还是说母语吧！"
		o, err = fortune.Draw(
			fortune.Category("literature", 40),
			fortune.Category("riddles", 30),
			fortune.Category("fortunes", 30),
		)
		if err != nil {
			log.Printf("error: %#v", err)
		}
	} else {
		o, err = fortune.Draw(
			fortune.Category("tang300", 50),
			fortune.Category("song100", 50),
		)
		if err != nil {
			log.Printf("error: %#v", err)
		}
	}

	if o != "" {
		m = m + "\n" + o
	}

	bot.Reply(evt, m, slackbot.WithTyping)
}
