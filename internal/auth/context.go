package auth

import (
	"context"
	"errors"
	"log"
)

var UserCtxKey = &contextKey{"user"}
var RespWriter = &contextKey{"respWriter"}

type contextKey struct {
	name string
}

// RequireAuth etc.
func RequireAuth(ctx context.Context) (string, error) {
	op := "auth.RequireAuth()"
	log.Printf("%s", op)

	user, ok := ctx.Value(UserCtxKey).(string)
	if !ok || user == "" {
		return "", errors.New("未授权，请先登录并携带有效 Token")
	}
	return user, nil
}
