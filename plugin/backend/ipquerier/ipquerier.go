package ipquerier

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
}

type IPQuerier struct {
	*plugin.Base
	config Config
}

func New(config Config) *IPQuerier {
	if config.ID == "" {
		config.ID = "ipQuerier"
	} else {
		config.ID = "ipQuerier-" + config.ID
	}

	q := &IPQuerier{
		Base:   plugin.NewBase(config.BaseConfig),
		config: config,
	}

	return q
}

func (p *IPQuerier) Handle(ctx context.Context, data interface{}) {
	fe := ctx.Value(plugin.CtxKeyFrontend).(plugin.Frontend)
	msg := data.(*event.Message)
	dict := ctx.Value(plugin.CtxKeyMatchedNames).(map[string]string)
	resp := ""

	resp, err := queryIP(dict["IP"])
	if err != nil {
		resp = fmt.Sprintf("API Error: %s", err)
		p.Logger.Print(resp)
	}

	fe.SendTextMessage(msg.Channel, resp)
}

func queryIP(ip string) (string, error) {
	url := "http://freeapi.ipip.net/" + ip
	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var result []string
	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	resp := fmt.Sprintf("IP: %s\nLocation: %v", ip, result)

	return resp, nil
}
