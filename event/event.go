package event

import (
	"context"
	"fmt"

	"github.com/flw-cn/slack-bot/types"
)

type Type string

const (
	EvJoinChannel   Type = "JoinedChannel"
	EvMessage       Type = "Message"
	EvDirectMessage Type = "DirectMessage"
	EvDirectMention Type = "DirectMention"
	EvMentionedMe   Type = "MentionedMe"
	EvFileMessage   Type = "FileMessage"
	EvUnhandled     Type = "UnhandledEvent"
)

type Event struct {
	Ctx  context.Context
	Type Type
	Data fmt.Stringer
}

func NewEvent(ctx context.Context, t Type, data fmt.Stringer) *Event {
	return &Event{
		Ctx:  ctx,
		Type: t,
		Data: data,
	}
}

type MessageType string

const (
	DirectMessage  MessageType = "DirectMessage"
	DirectMention  MessageType = "DirectMention"
	MentionedMe    MessageType = "MentionedMe"
	ChannelMessage MessageType = "ChannelMessage"
)

type Message struct {
	Type      MessageType
	Channel   types.Channel
	User      types.User
	Mentioned []types.User
	Text      string
}

func (m Message) String() string {
	var str string

	if len(m.Mentioned) > 0 {
		str = fmt.Sprintf("[%s] %s %s said to %s: %s", m.Type, m.Channel, m.User, m.Mentioned, m.Text)
	} else {
		str = fmt.Sprintf("[%s] %s %s said: %s", m.Type, m.Channel, m.User, m.Text)
	}

	return str
}

type FileInfo struct {
	ID        string
	Type      string
	Name      string
	Channel   types.Channel
	User      types.User
	Mentioned []types.User
	Comment   string
}

type File interface {
	fmt.Stringer
	Info() FileInfo
	Download() (string, func(), error)
}
