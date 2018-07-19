package bot

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	api "github.com/flw-cn/slack"
	"github.com/flw-cn/slack-bot/event"
	"github.com/flw-cn/slack-bot/plugin"
	"github.com/flw-cn/slack-bot/plugin/backend/greeter"
	"github.com/flw-cn/slack-bot/plugin/frontend/slack"
	slacktest "github.com/lusis/slack-test"
	"github.com/stretchr/testify/assert"
)

var slackToken = "<YOU SLACK TOKEN>"

func init() {
	slackToken = os.Getenv("SLACK_TOKEN")
}

func setupTestEnv(t *testing.T) (*slacktest.Server, *Bot) {
	s := setupTestServer()
	bot := setupBot(t)
	return s, bot
}

func setupTestServer() *slacktest.Server {
	s := slacktest.NewTestServer()
	s.SetBotName("foobot")
	api.SLACK_API = "http://" + s.ServerAddr + "/"
	go s.Start()

	return s
}

func setupBot(t *testing.T) *Bot {
	botConfig := Config{}
	bot := New(botConfig)

	slackConfig := slack.Config{
		UseRTMStart: true,
	}
	slack := slack.New(slackConfig)

	greeterConfig := greeter.Config{}
	greeter := greeter.New(greeterConfig)

	err := bot.LoadFrontend(slack)
	assert.NoError(t, err, "must no error")
	err = bot.LoadBackend(greeter)
	assert.NoError(t, err, "must no error")
	err = bot.Init()
	assert.NoError(t, err, "must no error")
	err = bot.Start()
	assert.NoError(t, err, "must no error")

	return bot
}

type SawMessage struct {
	ID      int
	Channel string
	Text    string
	Type    string
}

const (
	maxWait = 1 * time.Second
)

func TestChannelMessage(t *testing.T) {
	s, bot := setupTestEnv(t)
	defer s.Stop()
	defer bot.Stop()

	bot.On(event.EvMessage).Call(func(ctx context.Context, data interface{}) {
		fe := ctx.Value(plugin.CtxKeyFrontend).(plugin.Frontend)
		msg := data.(*event.Message)
		fe.SendTextMessage(msg.Channel, "Your said: "+msg.Text)
	})

	s.SendMessageToChannel("#random", "hello")
	select {
	case m := <-s.SeenFeed:
		msg := &SawMessage{}
		err := json.Unmarshal([]byte(m), msg)
		assert.NoError(t, err)
		assert.Equal(t, "Your said: hello", msg.Text, "but bot didn't respond")
	case <-time.After(maxWait):
		assert.FailNow(t, "did not get channel message in time")
	}
}

func TestDirectMessage(t *testing.T) {
	s, bot := setupTestEnv(t)
	defer s.Stop()
	defer bot.Stop()

	bot.On(event.EvDirectMessage).Call(func(ctx context.Context, data interface{}) {
		fe := ctx.Value(plugin.CtxKeyFrontend).(plugin.Frontend)
		msg := data.(*event.Message)
		fe.SendTextMessage(msg.Channel, "Your said: "+msg.Text)
	})

	s.SendDirectMessageToBot("hello")
	select {
	case m := <-s.SeenFeed:
		msg := &SawMessage{}
		err := json.Unmarshal([]byte(m), msg)
		assert.NoError(t, err)
		assert.Equal(t, "Your said: hello", msg.Text, "but bot didn't respond")
	case <-time.After(maxWait):
		assert.FailNow(t, "did not get channel message in time")
	}
}

func TestRegexpMatcher(t *testing.T) {
	s, bot := setupTestEnv(t)
	defer s.Stop()
	defer bot.Stop()

	bot.Hear(`^ip\s+(?P<IP>\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3})$`).Call(func(ctx context.Context, data interface{}) {
		fe := ctx.Value(plugin.CtxKeyFrontend).(plugin.Frontend)
		msg := data.(*event.Message)
		v := ctx.Value(plugin.CtxKeyMatchedNames)
		if v == nil {
			return
		}
		dict := v.(map[string]string)
		fe.SendTextMessage(msg.Channel, "Someone wants to query IP "+dict["IP"])
	})

	s.SendMessageToChannel("#random", "ip 123.45.67.89")

	select {
	case m := <-s.SeenFeed:
		msg := &SawMessage{}
		err := json.Unmarshal([]byte(m), msg)
		assert.NoError(t, err)
		assert.Equal(t, "Someone wants to query IP 123.45.67.89", msg.Text, "but bot didn't respond")
	case <-time.After(maxWait):
		assert.FailNow(t, "did not get channel message in time")
	}
}

func TestSubrouter(t *testing.T) {
}
