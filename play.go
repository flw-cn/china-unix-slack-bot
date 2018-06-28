package main

import (
	"context"
	"fmt"
	"os/exec"
	"syscall"

	slackbot "github.com/flw-cn/go-slackbot"
	"github.com/flw-cn/playground/docker"
	"github.com/flw-cn/slack"
	log "github.com/sirupsen/logrus"
)

type PlaygroundConfig struct {
	Docker string `flag:"||docker path"`
	Host   string `flag:"||host path"`
}

func Init(config PlaygroundConfig) error {
	return docker.Boarding(config.Host, config.Docker)
}

func PlayGo(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	codeInfo := slackbot.NamedCapturesFromContext(ctx)

	lang := codeInfo.Get("lang")
	code := codeInfo.Get("code")

	m := fmt.Sprintf("咦？有人发代码，让我试着运行一下！")
	bot.Reply(evt, m, slackbot.WithTyping)

	output, err := docker.PlayCode(lang, code)
	if err != nil {
		m = fmt.Sprintf("运行结果出错啦。\nError: %s\n%s", err, output)
		if e, ok := err.(*exec.ExitError); ok {
			status := e.ProcessState.Sys().(syscall.WaitStatus)
			if status.Signaled() && status.Signal() == syscall.SIGKILL {
				m = m + "\n看起来像是运行时间太长超时了，要不你再检查一下？"
			}
		}
		bot.Reply(evt, m, slackbot.WithTyping)
		return
	}

	if output != "" {
		m = fmt.Sprintf("运行结果出来了！\n%s", output)
		bot.Reply(evt, m, slackbot.WithTyping)
	}
}

func PlayGoFile(ctx context.Context, bot *slackbot.Bot, evt *slack.MessageEvent) {
	lang := evt.File.Filetype
	file, cleaner, err := slackDownloadFile(config.Play.Docker, bot.Client, evt.File.ID)
	if err != nil {
		log.Printf("Error: %s\n", err)
		return
	}

	defer cleaner()

	m := fmt.Sprintf("你这段代码能运行吗？让我试着运行一下！")
	bot.ReplyInThread(evt, m, slackbot.WithTyping)

	output, err := docker.PlayFile(lang, file)
	if err != nil {
		log.Printf("Error: %s\n%s", err, output)
		return
	}

	if output != "" {
		m = fmt.Sprintf("运行结果出来了！\n%s", output)
		bot.ReplyInThread(evt, m, slackbot.WithTyping)
	}
}
