package mongo

import "context"

type contextKey struct {
	key string
}

var (
	ServerContextKey = &contextKey{"mongo-server"}
)

// ServerFromContext retrieves the server from the context.
func ServerFromContext(ctx context.Context) *Server {
	svr := ctx.Value(ServerContextKey)
	if svr == nil {
		return nil
	}

	return svr.(*Server)
}
