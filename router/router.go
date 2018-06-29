package router

import (
	"context"

	"github.com/flw-cn/slack-bot/event"
)

type Router struct {
	// Routes to be matched, in order.
	routes []*Route
	// Holds any error that occurred during the function chain that created this router.  Used by subrouters.
	// err error
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Match(ctx context.Context, data interface{}) (Handler, context.Context) {
	for _, route := range r.routes {
		if handler, ctx := route.Match(ctx, data); handler != nil {
			return handler, ctx
		}
	}

	return nil, ctx
}

// addRoute registers an empty route
func (r *Router) addRoute() *Route {
	route := &Route{}
	r.routes = append(r.routes, route)
	return route
}

// Hear add a new route to hear specified message
func (r *Router) Hear(regex string) *Route {
	return r.addRoute().Hear(regex)
}

// On add a new route to watch specified event
func (r *Router) On(types ...event.Type) *Route {
	return r.addRoute().Messages(types...)
}

// When adds a matcher
func (r *Router) When(m Matcher) *Route {
	return r.addRoute().When(m)
}
