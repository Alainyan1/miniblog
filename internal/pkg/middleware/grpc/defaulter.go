// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package grpc

import (
	"context"

	"google.golang.org/grpc"
)

// 对请求进行默认值设置
func DefaultInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// 调用Default方法, 若存在
		if defaulter, ok := req.(interface{ Default() }); ok {
			defaulter.Default()
		}

		// 继续处理请求
		return handler(ctx, req)
	}
}
