package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"

	slackbot "github.com/flw-cn/go-slackbot"
	"github.com/flw-cn/slack"
)

func Whatever(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	r := rand.Intn(100)
	if r > 95 {
		Fortune(ctx, bot, evt)
		return
	}

	text := evt.Text
	text = slackbot.StripDirectMention(text)

	m, err := tuling123(text)
	if err != nil {
		Fortune(ctx, bot, evt)
		return
	}

	bot.Reply(evt, m, slackbot.WithTyping)
}

func tuling123(words string) (string, error) {
	url := "http://www.tuling123.com/openapi/api"
	key := os.Getenv("TULING_API_KEY")
	url = fmt.Sprintf("%s?key=%s&info=%s", url, key, words)

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	text := struct {
		Text string
	}{}

	err = json.Unmarshal(body, &text)
	if err != nil {
		fmt.Printf("stderr: %#v\n", err)
		return "", err
	}

	return text.Text, nil
}
