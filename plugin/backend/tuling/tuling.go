package tuling

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/flw-cn/slack-bot/event"
	"github.com/flw-cn/slack-bot/plugin"
)

type Config struct {
	plugin.BaseConfig
	URL   string `flag:"|http://www.tuling123.com/openapi/api|API {URL}"`
	Token string `flag:"||must provide your API {token} here"`
}

type Tuling struct {
	*plugin.Base
	config Config
}

func New(config Config) *Tuling {
	if config.ID == "" {
		config.ID = "tuling"
	} else {
		config.ID = "tuling-" + config.ID
	}

	t := &Tuling{
		Base:   plugin.NewBase(config.BaseConfig),
		config: config,
	}

	return t
}

func (t *Tuling) Handle(ctx context.Context, data interface{}) {
	fe := ctx.Value(plugin.CtxKeyFrontend).(plugin.Frontend)
	msg := data.(*event.Message)

	m, err := t.callAPI(msg.Text)
	if err != nil {
		t.Logger.Printf("Error: %v", err)
		return
	}

	fe.SendTextMessage(msg.Channel, m)
}

func (t *Tuling) callAPI(words string) (string, error) {
	url := fmt.Sprintf("%s?key=%s&info=%s", t.config.URL, t.config.Token, words)

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
