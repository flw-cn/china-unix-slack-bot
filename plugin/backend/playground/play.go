package playground

import (
	"context"
	"fmt"
	"os/exec"
	"syscall"

	"github.com/flw-cn/playground/docker"
	"github.com/flw-cn/slack-bot/event"
	"github.com/flw-cn/slack-bot/plugin"
)

type Config struct {
	plugin.BaseConfig
	Docker string `flag:"||the {path} used to store code inside Docker"`
	Host   string `flag:"||the real {path} used to store code outside Docker"`
}

type Playground struct {
	*plugin.Base
	config Config
}

func New(config Config) *Playground {
	if config.ID == "" {
		config.ID = "playground"
	} else {
		config.ID = "playground-" + config.ID
	}

	p := &Playground{
		Base:   plugin.NewBase(config.BaseConfig),
		config: config,
	}

	return p
}

func (p *Playground) Init() error {
	err := p.Base.Init()
	if err != nil {
		return err
	}

	return docker.Boarding(p.config.Host, p.config.Docker)
}

func (p *Playground) Handle(ctx context.Context, data interface{}) {
	fe := ctx.Value(plugin.CtxKeyFrontend).(plugin.Frontend)
	msg := data.(*event.Message)
	dict := ctx.Value(plugin.CtxKeyMatchedNames).(map[string]string)

	lang := dict["lang"]
	code := dict["code"]

	m := fmt.Sprintf("咦？有人发代码，让我试着运行一下！")
	fe.SendTextMessage(msg.Channel, m)

	output, err := docker.PlayCode(lang, code)
	if err != nil {
		m = fmt.Sprintf("运行结果出错啦。\nError: %s\n%s", err, output)
		if e, ok := err.(*exec.ExitError); ok {
			status := e.ProcessState.Sys().(syscall.WaitStatus)
			if status.Signaled() && status.Signal() == syscall.SIGKILL {
				m = m + "\n看起来像是运行时间太长超时了，要不你再检查一下？"
			}
		}
		fe.SendTextMessage(msg.Channel, m)
		return
	}

	if output != "" {
		m = fmt.Sprintf("运行结果出来了！\n```\n%s\n```", output)
		fe.SendTextMessage(msg.Channel, m)
	}
}

/*
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
*/
