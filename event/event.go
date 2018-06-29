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
	Data fmt.Stringer
}

func NewEvent(ctx context.Context, data fmt.Stringer) *Event {
	return &Event{
		Ctx:  ctx,
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
	if len(m.Mentioned) > 0 {
		return fmt.Sprintf("[%s] %s %s said to %s: %s", m.Type, m.Channel, m.User, m.Mentioned, m.Text)
	} else {
		return fmt.Sprintf("[%s] %s %s said: %s", m.Type, m.Channel, m.User, m.Text)
	}
}
