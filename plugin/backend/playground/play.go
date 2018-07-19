package playground

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
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

	switch msg := data.(type) {
	case *event.Message:
		p.PlayGoCode(ctx, fe, msg)
	case event.File:
		p.PlayGoFile(ctx, fe, msg)
	}
}

func (p *Playground) PlayGoCode(ctx context.Context, fe plugin.Frontend, msg *event.Message) {
	dict := ctx.Value(plugin.CtxKeyMatchedNames).(map[string]string)

	lang := dict["lang"]
	code := dict["code"]

	lang = strings.ToLower(lang)
	if lang != "go" && lang != "golang" {
		return
	}

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
		output = strings.TrimSpace(output)
		m = fmt.Sprintf("运行结果出来了！\n```\n%s\n```", output)
		fe.SendTextMessage(msg.Channel, m)
	}
}

func (p *Playground) PlayGoFile(ctx context.Context, fe plugin.Frontend, file event.File) {
	fileInfo := file.Info()
	codeFile, cleaner, err := file.Download()
	if err != nil {
		p.Debugf("ERROR: %v", err)
		return
	}

	defer cleaner()

	m := fmt.Sprintf("你这段代码能运行吗？让我试着运行一下！")
	fe.SendTextMessage(fileInfo.Channel, m)

	output, err := docker.PlayFile(fileInfo.Type, codeFile)
	if err != nil {
		m = fmt.Sprintf("运行结果出错啦。\nError: %s\n%s", err, output)
		if e, ok := err.(*exec.ExitError); ok {
			status := e.ProcessState.Sys().(syscall.WaitStatus)
			if status.Signaled() && status.Signal() == syscall.SIGKILL {
				m = m + "\n看起来像是运行时间太长超时了，要不你再检查一下？"
			}
		}
		fe.SendTextMessage(fileInfo.Channel, m)
		return
	}

	if output != "" {
		output = strings.TrimSpace(output)
		m = fmt.Sprintf("运行结果出来了！\n```\n%s\n```", output)
		fe.SendTextMessage(fileInfo.Channel, m)
	}
}
