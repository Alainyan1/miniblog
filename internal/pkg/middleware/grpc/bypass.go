// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package grpc

// bypass认证中间件会从请求头中获取用户ID数据并放通所有请求

import (
	"context"
	"miniblog/internal/pkg/contextx"
	"miniblog/internal/pkg/known"
	"miniblog/internal/pkg/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// grpc拦截器, 模拟所有请求都通过认证.
func AuthnBypasswInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// 从请求头获取用户ID
		userID := "user-000001" // 默认用户ID
		// rpc从metadata获取用户id, 类似http中的header?
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			// 获取key为x-user-id的值, x-user-id保存了UserID的值
			if values := md.Get(known.XUserID); len(values) > 0 {
				userID = values[0]
			}
		}

		log.Debugw("Simulated authentication successful", "userID", userID)

		// 讲默认信息存入上下文
		//nolint: staticcheck
		ctx = context.WithValue(ctx, known.XUserID, userID)

		// 为long和contextx提供上下文支持
		ctx = contextx.WithUserID(ctx, userID)

		// 继续处理请求
		return handler(ctx, req)
	}
}
