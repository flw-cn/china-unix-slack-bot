package plugin

import (
	"context"
	"io"

	"github.com/flw-cn/slack-bot/event"
	"github.com/flw-cn/slack-bot/types"
)

type Service interface {
	Runable
	Debugable
}

type Backend interface {
	Service
	Handle(context.Context, interface{})
}

type Frontend interface {
	Service
	IncomingEvents() <-chan *event.Event
	SendTextMessage(types.Channel, string)
}

type Runable interface {
	ID() string
	Init() error
	Start() error
	Stop() error
}

type Debugable interface {
	Logable
	SetDebug(bool)
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
}

type Logable interface {
	SetLogger(Logger)
	SetLogOutput(io.Writer)
}

type Logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	SetOutput(io.Writer)
}
