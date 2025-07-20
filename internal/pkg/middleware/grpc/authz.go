package grpc

import (
	"context"
	"miniblog/internal/pkg/contextx"
	"miniblog/internal/pkg/errno"
	"miniblog/internal/pkg/log"

	"google.golang.org/grpc"
)

type Authorize interface {
	Authorize(subject, object, action string) (bool, error)
}

func AuthzInterceptor(authorize Authorize) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		subject := contextx.UserID(ctx)
		object := info.FullMethod
		action := "CALL"

		// 记录授权上下文信息
		log.Debugw("Build authorize context", "subject", subject, "object", object, "action", action)

		// 调用授权接口进行认证
		if allowed, err := authorize.Authorize(subject, object, action); err != nil || !allowed {
			return nil, errno.ErrPermissionDenied.WithMessage(
				"access denied: subject=%s, object=%s, action=%s, reason=%v",
				subject,
				object,
				action,
				err,
			)
		}

		return handler(ctx, req)
	}
}
