// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package contextx

import "context"

type (
	// 用户id的上下文键
	userIDKey struct{}

	// 请求id的上下文键
	requestIDKey struct{}
)

// 将用户ID存放到上下文中
func WithUserID(ctx context.Context, userID string) context.Context {
	return context.WithValue(ctx, userIDKey{}, userID)
}

// 从上下文中提取用户ID
func UserID(ctx context.Context) string {
	// .(string)是一个类型断言, 尝试将interface{}转为指定类型
	userID, _ := ctx.Value(userIDKey{}).(string)
	return userID
}

// 将请求ID存放到上下文中
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDKey{}, requestID)
}

// 从上下文中提取请求ID
func RequestID(ctx context.Context) string {
	requestID, _ := ctx.Value(requestIDKey{}).(string)
	return requestID
}
