package context

import (
	"context"
)

func AddUserIdToContext(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, userIdContextKey, userId)
}

func GetUserIdFromContext(ctx context.Context) string {
	return ctx.Value(userIdContextKey).(string)
}
