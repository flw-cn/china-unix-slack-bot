package router

import (
	"context"
)

type Handler interface {
	Handle(context.Context, interface{})
}

// EventHandler is a event handler
type EventHandler func(ctx context.Context, data interface{})

func (h EventHandler) Handle(ctx context.Context, data interface{}) {
	h(ctx, data)
}
