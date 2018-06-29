package fortune

import (
	"context"
	"math/rand"
	"strings"

	"github.com/flw-cn/go-fortune"
	"github.com/flw-cn/slack-bot/event"
	"github.com/flw-cn/slack-bot/plugin"
)

type Config struct {
	plugin.BaseConfig
}

type Fortune struct {
	*plugin.Base
	config Config
}

func New(config Config) *Fortune {
	if config.ID == "" {
		config.ID = "fortune"
	} else {
		config.ID = "fortune-" + config.ID
	}

	f := &Fortune{
		Base:   plugin.NewBase(config.BaseConfig),
		config: config,
	}

	return f
}

func (f *Fortune) Handle(ctx context.Context, data interface{}) {
	fe := ctx.Value(plugin.CtxKeyFrontend).(plugin.Frontend)
	msg := data.(*event.Message)
	dict := ctx.Value(plugin.CtxKeyMatchedNames).(map[string]string)

	var m, o string
	var err error

	words := dict["words"]

	if words == "唐诗" {
		o, err = fortune.Draw(fortune.Category("tang300", 100))
	} else if words == "宋词" {
		o, err = fortune.Draw(fortune.Category("song100", 100))
	} else if strings.ContainsAny(words, "诗词") {
		o, err = fortune.Draw(
			fortune.Category("tang300", 50),
			fortune.Category("song100", 50),
		)
	} else {
		r := rand.Intn(100)
		if r < 5 {
			m = "中国话太难背了，咱还是说母语吧！\n"
			o, err = fortune.Draw(
				fortune.Category("literature", 40),
				fortune.Category("riddles", 30),
				fortune.Category("fortunes", 30),
			)
		} else {
			m = "你想听唐诗还是宋词？"
		}
	}

	if err != nil {
		f.Logger.Printf("error: %v", err)
		return
	}

	if o != "" {
		m = m + o
	}

	if m != "" {
		fe.SendTextMessage(msg.Channel, m)
	}
}
