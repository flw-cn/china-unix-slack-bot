package router

import (
	"context"
	"regexp"

	"github.com/flw-cn/slack-bot/event"
	"github.com/flw-cn/slack-bot/plugin"
	"github.com/flw-cn/slack-bot/util"
)

// Matcher type for matching message routes
type Matcher interface {
	Match(context.Context, interface{}) (bool, context.Context)
}

// RegexpMatcher is a regexp matcher
type RegexpMatcher struct {
	regex *regexp.Regexp
}

// NewRegexpMatcher builds a new RegexpMatcher
//
// regex must be a valid regular expression
func NewRegexpMatcher(regex string) *RegexpMatcher {
	re := regexp.MustCompile(regex)
	return &RegexpMatcher{regex: re}
}

// Match matches an event
func (rm *RegexpMatcher) Match(ctx context.Context, data interface{}) (bool, context.Context) {
	msg, ok := data.(*event.Message)
	if !ok {
		return false, ctx
	}

	ok, dict := util.NamedRegexpParse(msg.Text, rm.regex)
	if !ok {
		return false, ctx
	}

	ctx = context.WithValue(ctx, plugin.CtxKeyMatchedNames, dict)
	return true, ctx
}

// TypesMatcher is a matcher based of event type
type TypesMatcher struct {
	types map[event.Type]bool
}

// NewTypesMatcher builds a new TypesMatcher
func NewTypesMatcher(types []event.Type) *TypesMatcher {
	dict := map[event.Type]bool{}
	for _, t := range types {
		dict[t] = true
	}
	return &TypesMatcher{types: dict}
}

// Match matches an event
func (tm *TypesMatcher) Match(ctx context.Context, data interface{}) (bool, context.Context) {
	switch msg := data.(type) {
	case *event.Message:
		switch msg.Type {
		case "DirectMessage":
			if _, ok := tm.types[event.EvDirectMessage]; ok {
				return true, ctx
			}
		case "DirectMention":
			if _, ok := tm.types[event.EvDirectMention]; ok {
				return true, ctx
			}
		case "MentionedMe":
			if _, ok := tm.types[event.EvMentionedMe]; ok {
				return true, ctx
			}
		case "ChannelMessage":
			if _, ok := tm.types[event.EvMessage]; ok {
				return true, ctx
			}
		}
	}

	return false, ctx
}
