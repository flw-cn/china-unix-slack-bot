package router

import (
	"context"
	"regexp"
	"strings"

	"github.com/flw-cn/slack-bot/event"
	"github.com/flw-cn/slack-bot/plugin"
	"github.com/flw-cn/slack-bot/util"
)

// Matcher type for matching message routes
type Matcher interface {
	Match(context.Context, *event.Event) (bool, context.Context)
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
func (rm *RegexpMatcher) Match(ctx context.Context, ev *event.Event) (bool, context.Context) {
	msg, ok := ev.Data.(*event.Message)
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
func (tm *TypesMatcher) Match(ctx context.Context, ev *event.Event) (bool, context.Context) {
	_, ok := tm.types[ev.Type]
	return ok, ctx
}

// FileTypesMatcher is a matcher based of event type
type FileTypesMatcher struct {
	types map[string]bool
}

// NewFileTypesMatcher builds a new FileTypesMatcher
func NewFileTypesMatcher(types []string) *FileTypesMatcher {
	dict := map[string]bool{}
	for _, t := range types {
		t = strings.ToLower(t)
		dict[t] = true
	}
	return &FileTypesMatcher{types: dict}
}

// Match matches an event
func (ftm *FileTypesMatcher) Match(ctx context.Context, ev *event.Event) (bool, context.Context) {
	if ev.Type != event.EvFileMessage {
		return false, ctx
	}

	file, ok := ev.Data.(event.File)
	if !ok {
		return false, ctx
	}

	t := strings.ToLower(file.Info().Type)
	_, ok = ftm.types[t]
	return ok, ctx
}
