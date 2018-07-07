package ruyi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/flw-cn/slack-bot/event"
	"github.com/flw-cn/slack-bot/plugin"
)

type Config struct {
	plugin.BaseConfig
	URL    string `flag:"|http://api.ruyi.ai/v1/message|API {URL}"`
	AppKey string `flag:"||must provide your {AppKey} here"`
}

type Ruyi struct {
	*plugin.Base
	config Config
}

func New(config Config) *Ruyi {
	if config.ID == "" {
		config.ID = "ruyi"
	} else {
		config.ID = "ruyi-" + config.ID
	}

	g := &Ruyi{
		Base:   plugin.NewBase(config.BaseConfig),
		config: config,
	}

	return g
}

func (o *Ruyi) Handle(ctx context.Context, data interface{}) {
	fe := ctx.Value(plugin.CtxKeyFrontend).(plugin.Frontend)
	msg := data.(*event.Message)

	m, err := o.callAPI(msg.User.String(), msg.Text)
	if err != nil {
		o.Logger.Printf("Error: %v", err)
		return
	}

	fe.SendTextMessage(msg.Channel, m)
}

func (o *Ruyi) callAPI(user, words string) (string, error) {
	req := &Request{
		AppKey: o.config.AppKey,
		UserID: user,
		Q:      words,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBuffer(body)
	res, err := http.Post(o.config.URL, "application/json", buf)
	if err != nil {
		return "", err
	}

	body, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	o.Debugf("Result: %s", body)

	result := &Response{}
	err = json.Unmarshal(body, result)
	if err != nil {
		fmt.Printf("Error: %#v\n", err)
		return "", err
	}

	o.Debugf("Result: %#v", result)
	if result.Result.Intents == nil {
		return "", errors.New("result is empty")
	}

	return result.Result.Intents[0].Result["text"].(string), nil
}
