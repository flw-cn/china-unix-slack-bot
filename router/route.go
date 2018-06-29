package router

import (
	"context"

	"github.com/flw-cn/slack-bot/event"
)

// Route represents a route item
//
// A Router can contained several routes, in order.
type Route struct {
	matcher   Matcher
	handler   Handler
	subrouter *Router
}

// Hear adds a RegexpMatcher for the message text
func (r *Route) Hear(regex string) *Route {
	r.matcher = NewRegexpMatcher(regex)
	return r
}

// Messages sets the types of Messages we want to handle
func (r *Route) Messages(types ...event.Type) *Route {
	r.matcher = NewTypesMatcher(types)
	return r
}

func (r *Route) When(matcher Matcher) *Route {
	r.matcher = matcher
	return r
}

// Call sets a handler for the route
func (r *Route) Call(handler EventHandler) *Route {
	r.handler = Handler(handler)
	return r
}

func (r *Route) Hook() *Route {
	r.handler = nil
	return r
}

// Subrouter creates a subrouter
func (r *Route) Subrouter() *Router {
	r.subrouter = &Router{}
	return r.subrouter
}

// Match matches
func (r *Route) Match(ctx context.Context, data interface{}) (Handler, context.Context) {
	matched, newCtx := r.matcher.Match(ctx, data)
	if !matched {
		return nil, ctx
	}

	if r.handler != nil {
		return r.handler, newCtx
	}

	if r.subrouter != nil {
		return r.subrouter.Match(newCtx, data)
	}

	return nil, ctx
}
