package context

import (
	"context"

	"github.com/LouisHatton/menu-link-up/internal/users"
)

func AddUserToContext(ctx context.Context, user users.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

func GetUserFromContext(ctx context.Context) users.User {
	return ctx.Value(userContextKey).(users.User)
}
