// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package gin

import (
	"miniblog/internal/pkg/contextx"
	"miniblog/internal/pkg/known"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// 实现gin请求的中间件

// gin中间件, 用于在每个http请求的上下文和响应中注入'x-request-id'键值对.
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取`x-request-id`, 若不存在, 生成新的uuid
		requestID := c.Request.Header.Get(known.XRequestID)

		if requestID == "" {
			requestID = uuid.New().String()
		}

		// 将RequestID保存到conetxt.Context中, 以便后续程序使用
		ctx := contextx.WithRequestID(c.Request.Context(), requestID)
		c.Request = c.Request.WithContext(ctx)

		// 将requestid保存到http返回头中, header到键为`x-request-id`
		c.Writer.Header().Set(known.XRequestID, requestID)

		c.Next()
	}
}
