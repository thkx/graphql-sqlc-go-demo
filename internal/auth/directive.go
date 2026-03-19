package auth

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/99designs/gqlgen/graphql"
	"github.com/thkx/graphql-sqlc-go-dome/internal/auth/jwt"
)

func AuthFieldMiddleware(ctx context.Context, obj any, next graphql.Resolver) (any, error) {
	fieldCtx := graphql.GetFieldContext(ctx)
	if fieldCtx == nil {
		return next(ctx)
	}

	if hasRequiresAuthDirective(fieldCtx.Field) {
		if err := validateAuth(ctx); err != nil {
			return nil, err
		}
	}

	return next(ctx)
}

func hasRequiresAuthDirective(field graphql.CollectedField) bool {
	// 遍历字段的所有指令，判断是否存在 @requiresAuth
	for _, v := range field.Directives {
		if v.Name == "requiresAuth" {
			return true
		}
	}
	return false
}

// validateAuth 统一校验上下文是否存在有效授权（复用原有 JWT 校验逻辑）
func validateAuth(ctx context.Context) error {
	username, ok := ctx.Value(RespWriter).(string)
	if ok && username != "" {
		return nil
	}

	if !ok {
		return errors.New("未授权：请先登录并携带有效 Token")
	}
	return nil
}

func validateAndGetUserID(cookie *http.Cookie) (string, error) {
	// op := "auth.validateAndGetUserID()"
	// fmt.Println(op)
	v, err := url.QueryUnescape(cookie.Value)
	if err != nil {
		return "", err
	}
	return jwt.ValidateToken(v)
}
