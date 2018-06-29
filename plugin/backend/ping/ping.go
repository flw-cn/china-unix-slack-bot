package ping

import (
	"context"

	"github.com/flw-cn/slack-bot/event"
	"github.com/flw-cn/slack-bot/plugin"
)

type Config struct {
	plugin.BaseConfig
}

type Ping struct {
	*plugin.Base
	config Config
}

func New(config Config) *Ping {
	if config.ID == "" {
		config.ID = "greeter"
	} else {
		config.ID = "greeter-" + config.ID
	}

	p := &Ping{
		Base:   plugin.NewBase(config.BaseConfig),
		config: config,
	}

	return p
}

func (p *Ping) Handle(ctx context.Context, data interface{}) {
	fe := ctx.Value(plugin.CtxKeyFrontend).(plugin.Frontend)
	msg := data.(*event.Message)
	resp := "pong"
	fe.SendTextMessage(msg.Channel, resp)
}
