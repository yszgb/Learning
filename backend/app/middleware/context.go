package middleware

import (
	"context"
	"net/http"

	"learning.com/learning/models"
)

type contextKey string

const UserContextKey contextKey = "user"

// 从请求上下文中获取用户信息
func (m *Middleware) ContextSetUser(r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user) // 将用户信息存入上下文
	// WithValue: 为给定的请求创建一个新的上下文，该上下文包含了键值对
	return r.WithContext(ctx)
}

func (m *Middleware) ContextGetUser(r *http.Request) *models.User {
	user, ok := r.Context().Value(UserContextKey).(*models.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}