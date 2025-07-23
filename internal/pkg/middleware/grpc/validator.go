// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package grpc

import (
	"context"

	"google.golang.org/grpc"
)

// 用于自定义验证的接口.
type RequestValidator interface {
	Validate(ctx context.Context, rq any) error
}

// grpc拦截器用于对请求进行验证.
func ValidatorInterceptor(validator RequestValidator) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, rq any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		// 调用自定义验证方法
		if err := validator.Validate(ctx, rq); err != nil {
			// 这里不用返回 errno.ErrInvalidArgument 类型的错误信息, 由 validator.Validate 返回.
			return nil, err
		}
		// 继续处理请求
		return handler(ctx, rq)
	}
}
