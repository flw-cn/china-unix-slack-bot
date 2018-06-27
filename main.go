package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/flw-cn/go-fortune"
	"github.com/flw-cn/go-slackbot"
	"github.com/flw-cn/go-smartConfig"
	"github.com/flw-cn/playground/docker"
	"github.com/flw-cn/slack"
	"github.com/pityonline/china-unix-slack-bot/service"
	log "github.com/sirupsen/logrus"
)

type Config struct {
	Debug bool   `flag:"d|false|debug mode, default to 'false'"`
	Token string `flag:"t||must provide your {SLACK_TOKEN} here"`
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
}

var Dumper = spew.ConfigState{
	Indent:                  " ",
	DisablePointerAddresses: true,
	DisableCapacities:       true,
	SortKeys:                true,
}

var config Config

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
	toMe.MessageHandler(Fortune)

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

func QueryIP(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	id := bot.BotUserID()
	id = "<@" + id + ">"
	text := strings.Replace(evt.Text, id, "", -1)
	strings.Trim(text, " ")
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
	file, err := slackDownloadFile(bot.Client, evt.File.ID)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return
	}

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

func slackDownloadFile(api *slack.Client, fileID string) (string, error) {
	file, _, _, err := api.GetFileInfo(fileID, 0, 0)
	if err != nil {
		return "", err
	}

	url := file.URLPrivateDownload

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer "+config.Token)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	tmpdir, err := ioutil.TempDir("", "slack-bot-downloaded-files-")
	if err != nil {
		return "", err
	}

	tmpFile := tmpdir + "/file"
	err = ioutil.WriteFile(tmpFile, content, 0666)
	if err != nil {
		return "", err
	}

	return tmpFile, nil
}
