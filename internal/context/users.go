package context

import (
	"context"
)

func AddUserIdToContext(ctx context.Context, userId string) context.Context {
	return context.WithValue(ctx, userIdContextKey, userId)
}

func GetUserIdFromContext(ctx context.Context) string {
	id := ctx.Value(userIdContextKey)
	if id == nil {
		panic("user id is not in context")
	}

	return id.(string)
}
