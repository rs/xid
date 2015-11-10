package xid

import "golang.org/x/net/context"

type key int

const idKey key = 0

// FromContext returns a uniq id associated to the request if any
func FromContext(ctx context.Context) (id ID, ok bool) {
	if ctx != nil {
		id, ok = ctx.Value(idKey).(ID)
	}
	return
}

// NewContext returns a copy of the parent context and associates it with passed xid.
func NewContext(ctx context.Context, id ID) context.Context {
	return context.WithValue(ctx, idKey, id)
}
