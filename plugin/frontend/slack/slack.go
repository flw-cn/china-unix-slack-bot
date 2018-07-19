package slack

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	api "github.com/flw-cn/slack"
	"github.com/flw-cn/slack-bot/event"
	"github.com/flw-cn/slack-bot/plugin"
	"github.com/flw-cn/slack-bot/types"
	"github.com/flw-cn/slack-bot/util"
)

type Config struct {
	plugin.BaseConfig
	Token       string `flag:"||must provide your Slack API {token} here"`
	UseRTMStart bool   `flag:"|false|Use rtm.start(=true) or rtm.connect(=false)"`
}

type Slack struct {
	*plugin.Base
	config Config

	eventChan chan *event.Event
	done      chan bool
	me        types.User
	Client    *api.Client
	RTM       *api.RTM
}

func New(config Config) *Slack {
	if config.ID == "" {
		config.ID = "slack"
	} else {
		config.ID = "slack-" + config.ID
	}

	s := &Slack{
		Base:      plugin.NewBase(config.BaseConfig),
		config:    config,
		eventChan: make(chan *event.Event, 1024),
		done:      make(chan bool, 1),
	}

	return s
}

func (s *Slack) Init() error {
	err := s.Base.Init()
	if err != nil {
		return err
	}

	if s.Client == nil {
		s.Debugf("SLACK_TOKEN: %s", s.config.Token)
		s.Client = api.New(s.config.Token)
	}

	options := api.RTMOptions{}
	options.UseRTMStart = s.config.UseRTMStart
	s.RTM = s.Client.NewRTMWithOptions(&options)

	if s.config.BaseConfig.Debug {
		// TODO: api.SetLogger(s.logger)
		api.SetLogger(log.New(os.Stderr, "SLACK-API ", log.LstdFlags))
		s.Client.SetDebug(true)
		s.RTM.SetDebug(true)
	}

	return nil
}

func (s *Slack) Start() error {
	err := s.Base.Start()
	if err != nil {
		return err
	}

	go s.RTM.ManageConnection()
	go s.run()

	return nil
}

func (s *Slack) Stop() error {
	// TODO: shutdown slack API
	close(s.done)
	return nil
}

func (s *Slack) IncomingEvents() <-chan *event.Event {
	return s.eventChan
}

func (s *Slack) SendTextMessage(channel types.Channel, text string) {
	s.RTM.SendMessage(s.RTM.NewOutgoingMessage(text, channel.ID))
}

func (s *Slack) run() {
	s.Debug("Enter main loop.")
LOOP:
	for {
		select {
		case <-s.done:
			break LOOP
		case evt := <-s.RTM.IncomingEvents:
			ctx := context.Background()
			switch ev := evt.Data.(type) {
			case *api.ConnectedEvent:
				s.me = types.User{ID: ev.Info.User.ID, Name: ev.Info.User.Name}
				s.Logger.Printf("Connected as %s", s.me)

			case *api.ChannelJoinedEvent:
				// s.eventChan <- NewJoinChannelEvent(ctx, ev)
			case *api.GroupJoinedEvent:
				// s.eventChan <- NewJoinGroupEvent(ctx, ev)

			case *api.MessageEvent:
				s.eventChan <- s.NewTextMessageEvent(ctx, ev)
			case *api.FileSharedEvent:
				s.eventChan <- s.NewFileEvent(ctx, ev)

			case *api.AckMessage:
				if ev.RTMResponse.Ok {
					s.Logger.Print("Reply ok: ", ev.Text)
				} else {
					s.Logger.Print("Reply error: ", ev.RTMResponse.Error)
				}

			case *api.ConnectionErrorEvent:
				s.Logger.Printf("Connect error: %s", ev.Error())
			case *api.InvalidAuthEvent:
				s.Logger.Print("Invalid credentials")
				return
			case *api.RTMError:
				s.Logger.Print(ev.Error())
			case *api.AckErrorEvent:
				s.Logger.Print("AckError: ", ev.Error())

			case *api.LatencyReport:
				s.Logger.Print("LatencyReport: ", ev.Value)

			default:
				s.Debugf("Event(%s): %#v", evt.Type, evt.Data)
				// s.eventChan <- NewUnhandledEvent(ctx, &evt)
			}
		}
	}

	s.Logger.Printf("Leave run loop.")
}

type Event struct {
	context     context.Context
	chanJoined  *api.ChannelJoinedEvent
	groupJoined *api.GroupJoinedEvent
	unhandled   *api.RTMEvent
	// message     *api.MessageEvent
}

func (s *Slack) NewTextMessageEvent(ctx context.Context, ev *api.MessageEvent) *event.Event {
	name := ""
	u, err := s.Client.GetUserInfo(ev.User)
	if err == nil {
		name = u.Name
	}

	originText := strings.TrimSpace(ev.Text)
	list := whoMentioned(originText)
	stripedText := stripAllMention(originText)

	mentioned := make([]types.User, len(list))
	for i, id := range list {
		mentioned[i].ID = id
		u, err = s.Client.GetUserInfo(id)
		if err == nil {
			mentioned[i].Name = u.Name
		}
	}

	msg := &event.Message{
		Channel: types.Channel{
			ID:   ev.Channel,
			Name: s.resolveChannelName(ev.Channel),
		},
		User: types.User{
			ID:   ev.User,
			Name: name,
		},
		Mentioned: mentioned,
		Text:      stripedText,
	}

	var evType event.Type
	if isDirectMessage(ev.Channel) {
		evType = event.EvDirectMessage
		msg.Type = event.DirectMessage
	} else if isDirectMention(originText, s.me.ID) {
		evType = event.EvDirectMention
		msg.Type = event.DirectMention
	} else if isMentioned(originText, s.me.ID) {
		evType = event.EvMentionedMe
		msg.Type = event.MentionedMe
	} else {
		evType = event.EvMessage
		msg.Type = event.ChannelMessage
	}

	return event.NewEvent(ctx, evType, msg)
}

func (s *Slack) NewFileEvent(ctx context.Context, ev *api.FileSharedEvent) *event.Event {
	slackFile, _, _, err := s.Client.GetFileInfo(ev.FileID, 10, 1)
	if err != nil {
		return nil
	}

	name := ""
	u, err := s.Client.GetUserInfo(ev.File.User)
	if err == nil {
		name = u.Name
	}

	channel := ""
	if len(slackFile.Channels) > 0 {
		channel = slackFile.Channels[0]
	} else if len(slackFile.Groups) > 0 {
		channel = slackFile.Groups[0]
	}

	file := &File{
		fileInfo: event.FileInfo{
			ID:   slackFile.ID,
			Type: slackFile.Filetype,
			Name: slackFile.Name,
			Channel: types.Channel{
				ID:   channel,
				Name: s.resolveChannelName(channel),
			},
			User: types.User{
				ID:   slackFile.User,
				Name: name,
			},
			Comment: slackFile.InitialComment.Comment,
		},
		url:   slackFile.URLPrivateDownload,
		token: s.config.Token,
	}

	return event.NewEvent(ctx, event.EvFileMessage, file)
}

type File struct {
	fileInfo event.FileInfo
	url      string
	token    string
}

func (f *File) String() string {
	return fmt.Sprintf("%s uploaded a file: <%s|%s> into %s and commented: %s",
		f.fileInfo.User, f.fileInfo.ID, f.fileInfo.Name, f.fileInfo.Channel, f.fileInfo.Comment)
}

func (f *File) Info() event.FileInfo {
	return f.fileInfo
}

func (f *File) Download() (string, func(), error) {
	codeFile, cleaner, err := util.DownloadFile(f.url, "Bearer "+f.token)
	if err != nil {
		return "", nil, err
	}

	return codeFile, cleaner, nil
}

func (s *Slack) resolveChannelName(id string) string {
	name := ""
	if strings.HasPrefix(id, "C") {
		c, _ := s.Client.GetChannelInfo(id)
		if c != nil {
			name = c.Name
		}
	} else if strings.HasPrefix(id, "G") {
		g, _ := s.Client.GetGroupInfo(id)
		if g != nil {
			name = g.Name
		}
	}

	return name
}

func NewJoinChannelEvent(ctx context.Context, ev *api.ChannelJoinedEvent) *Event {
	return &Event{
		context:    ctx,
		chanJoined: ev,
	}
}

func NewJoinGroupEvent(ctx context.Context, ev *api.GroupJoinedEvent) *Event {
	return &Event{
		context:     ctx,
		groupJoined: ev,
	}
}

func NewUnhandledEvent(ctx context.Context, ev *api.RTMEvent) *Event {
	return &Event{
		context:   ctx,
		unhandled: ev,
	}
}
