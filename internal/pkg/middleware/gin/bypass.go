// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package gin

import (
	"miniblog/internal/pkg/contextx"
	"miniblog/internal/pkg/known"
	"miniblog/internal/pkg/log"

	"github.com/gin-gonic/gin"
)

func AuthnBypasswMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// http从Header中提取用户ID, 假设请求头名称为 "X-User-ID"
		userID := "user-000001" // 默认用户ID
		if val := c.GetHeader(known.XUserID); val != "" {
			userID = val
		}

		log.Debugw("Simulated authentication successful", "userID", userID)

		// 将用户ID和用户名注入到上下文中
		ctx := contextx.WithUserID(c.Request.Context(), userID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
