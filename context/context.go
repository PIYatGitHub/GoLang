package context

import (
	"context"

	"models"
)

const (
	userKey privateKey = "user"
)

type privateKey string

//WithUser returns the user according to the context
func WithUser(ctx context.Context, user *models.User) context.Context {
	return context.WithValue(ctx, userKey, user)
}

//User returns the user according to the context
func User(ctx context.Context) *models.User {
	if temp := ctx.Value(userKey); temp != nil {
		if user, ok := temp.(*models.User); ok {
			return user
		}
	}
	return nil
}
