package bot

import (
	"context"
	"errors"
	"fmt"

	"github.com/flw-cn/slack-bot/event"
	"github.com/flw-cn/slack-bot/plugin"
	"github.com/flw-cn/slack-bot/router"
)

type Config struct {
	plugin.BaseConfig
}

// Bot is a bot
type Bot struct {
	*plugin.Base
	router.Router
	config Config

	initialized bool
	eventChan   chan *event.Event
	frontends   []plugin.Frontend
	backends    []plugin.Backend
}

func New(config Config) *Bot {
	b := &Bot{
		Base:      plugin.NewBase(config.BaseConfig),
		eventChan: make(chan *event.Event, 1024),
		config:    config,
	}

	return b
}

func (b *Bot) SetDebug(debug bool) {
	b.Base.SetDebug(debug)

	for _, be := range b.backends {
		be.SetDebug(debug)
	}

	for _, fe := range b.frontends {
		fe.SetDebug(debug)
	}
}

func (b *Bot) LoadBackend(be ...plugin.Backend) error {
	if b.initialized {
		return errors.New("OOPS! bot already initialized. Please read the document.")
	}

	b.backends = append(b.backends, be...)
	return nil
}

func (b *Bot) LoadFrontend(fe ...plugin.Frontend) error {
	if b.initialized {
		return errors.New("OOPS! bot already initialized. Please read the document.")
	}

	b.frontends = append(b.frontends, fe...)
	return nil
}

func (b *Bot) Init() error {
	debug := b.GetDebug()
	for _, be := range b.backends {
		be.SetLogger(b.Logger)
		be.SetDebug(debug)
		err := be.Init()
		if err != nil {
			return fmt.Errorf("Can't initialize backend %s: %v", be.ID(), err)
		}
	}

	for _, fe := range b.frontends {
		fe.SetLogger(b.Logger)
		fe.SetDebug(debug)
		err := fe.Init()
		if err != nil {
			return fmt.Errorf("Can't initialize frontend %s: %v", fe.ID(), err)
		}
	}

	b.initialized = true

	return nil
}

func (b *Bot) Start() error {
	if !b.initialized {
		return errors.New("The bot is not initialized yet. please call bot.Init() first.")
	}

	for _, fe := range b.frontends {
		err := fe.Start()
		if err != nil {
			return err
		}

		ch := fe.IncomingEvents()
		go func(lfe plugin.Frontend) {
			for e := range ch {
				e.Ctx = context.WithValue(e.Ctx, plugin.CtxKeyFrontend, lfe)
				b.eventChan <- e
			}
		}(fe)
	}

	go b.dispatch()

	return nil
}

func (b *Bot) dispatch() {
	for e := range b.eventChan {
		fe := e.Ctx.Value(plugin.CtxKeyFrontend).(plugin.Frontend)
		b.Logger.Printf("[%s] %s", fe.ID(), e.Data)
		handler, ctx := b.Match(e.Ctx, e.Data)
		if handler != nil {
			handler.Handle(ctx, e.Data)
		}
	}
}

func (b *Bot) Stop() error {
	for _, be := range b.backends {
		be.Stop()
	}
	for _, fe := range b.frontends {
		fe.Stop()
	}

	return nil
}

func (b *Bot) Mount(r *router.Route, be plugin.Backend) {
	r.Call(router.EventHandler(be.Handle))
}
