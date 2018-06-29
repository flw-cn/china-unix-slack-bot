package plugin

type ContextKey string

const (
	CtxKeyUserName     ContextKey = "UserName"
	CtxKeyFrontend     ContextKey = "Frontend"
	CtxKeyEvent        ContextKey = "Event"
	CtxKeyMatchedNames ContextKey = "MatchedNames"
)
