package greeter

import (
	"context"

	"github.com/flw-cn/slack-bot/event"
	"github.com/flw-cn/slack-bot/plugin"
)

type Config struct {
	plugin.BaseConfig
}

type Greeter struct {
	*plugin.Base
	config Config
}

func New(config Config) *Greeter {
	if config.ID == "" {
		config.ID = "greeter"
	} else {
		config.ID = "greeter-" + config.ID
	}

	g := &Greeter{
		Base:   plugin.NewBase(config.BaseConfig),
		config: config,
	}

	return g
}

func (g *Greeter) Handle(ctx context.Context, data interface{}) {
	fe := ctx.Value(plugin.CtxKeyFrontend).(plugin.Frontend)
	msg := data.(*event.Message)
	resp := "Hello, " + msg.User.Name + "!"
	fe.SendTextMessage(msg.Channel, resp)
}
