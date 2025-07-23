// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package contextx

import "context"

type (
	// 用户id的上下文键.
	userIDKey struct{}

	// 用户name的上下文键.
	usernameKey struct{}

	// 访问令牌的上下文键.
	accessTokenKey struct{}
	// 请求id的上下文键.
	requestIDKey struct{}
)

// 将用户ID存放到上下文中.
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

// 从上下文中提取用户ID.
func UserID(ctx context.Context) string {
	// .(string)是一个类型断言, 尝试将interface{}转为指定类型
	userID, _ := ctx.Value(userIDKey{}).(string)
	return userID
}

// 将用户名放到上下文中.
func WithUsername(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, usernameKey{}, username)
}

// 从上下文中提取提取用户名.
func Username(ctx context.Context) string {
	username, _ := ctx.Value(usernameKey{}).(string)
	return username
}

// 将accessToken放到上下文中.
func WithAccessToken(ctx context.Context, accessToken string) context.Context {
	return context.WithValue(ctx, accessTokenKey{}, accessToken)
}

// 从上下文中读取accessToken.
func AccessToken(ctx context.Context) string {
	accessToken, _ := ctx.Value(accessTokenKey{}).(string)
	return accessToken
}

// 将请求ID存放到上下文中.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, requestID)
}

// 从上下文中提取请求ID.
func RequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDKey{}).(string)
	return requestID
}
