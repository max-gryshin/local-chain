package pkg

import (
	"context"

	"github.com/hashicorp/raft"
)

type contextKeyType int

const serverIDKey contextKeyType = iota

func ContextWithServerID(ctx context.Context, serverID raft.ServerID) context.Context {
	return context.WithValue(ctx, serverIDKey, serverID)
}

func ServerIDFromContext(ctx context.Context) raft.ServerID {
	if ctx == nil {
		return ""
	}
	if span, ok := ctx.Value(serverIDKey).(raft.ServerID); ok {
		return span
	}
	return ""
}
