package context

type ContextKey struct {
	name string
}

var UserIDKey = &ContextKey{"userID"}
