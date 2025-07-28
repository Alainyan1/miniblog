// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package gin

import (
	"context"
	"miniblog/internal/apiserver/model"
	"miniblog/internal/pkg/contextx"
	"miniblog/internal/pkg/errno"
	"miniblog/internal/pkg/log"

	"github.com/onexstack/onexstack/pkg/token"

	"github.com/gin-gonic/gin"
	"github.com/onexstack/onexstack/pkg/core"
)

type UserRetriever interface {
	// GetUser 根据用户ID获取用户信息
	GetUser(ctx context.Context, userID string) (*model.UserM, error)
}

func AuthnMiddleware(retriever UserRetriever) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, err := token.ParseRequest(ctx)
		if err != nil {
			core.WriteResponse(ctx, nil, errno.ErrTokenInvalid.WithMessage(err.Error()))
			// 如果授权失败, abort会阻止调用待处理的处理程序
			ctx.Abort()
			return
		}

		log.Debugw("Token parsing successful", "userID", userID)

		user, err := retriever.GetUser(ctx, userID)
		if err != nil {
			core.WriteResponse(ctx, nil, errno.ErrUnauthenticated.WithMessage(err.Error()))
			ctx.Abort()
			return
		}

		c := contextx.WithUserID(ctx.Request.Context(), userID)
		c = contextx.WithUsername(c, user.Username)
		ctx.Request = ctx.Request.WithContext(c)

		ctx.Next()
	}
}
