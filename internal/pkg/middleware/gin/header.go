// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package gin

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// 实现了跨域中间件

// gin中间件, 用于禁止客户端缓存HTTP请求的返回结果, 确保每次请求都重新从服务器获取数据
// no-cache 表示客户端在再次使用缓存之前, 必须向服务器验证资源是否过期
// no-store: 更严格, 禁止客户端和中间代理(如 CDN)存储任何响应内容
// max-age=0: 缓存的最大有效时间为 0 秒, 强制缓存立即过期
// must-revalidate: 当缓存过期后, 客户端必须重新验证资源
func NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate")
	c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	// Last-Modified 指定资源的最后修改时间, 设置为当前时间
	// 客户端可能用此头与服务器的 If-Modified-Since 比较, 判断资源是否需要重新请求
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	// 调用 c.Next(), 将请求传递给下一个中间件或路由处理程序, 允许 Gin 框架继续处理请求
	c.Next()
}

// gin中间件, 用于处理CORS请求
func Cors(c *gin.Context) {
	// 处理预检请求
	if c.Request.Method == http.MethodOptions {
		// 允许所有源(*)进行跨域请求
		c.Header("Access-Control-Allow-Origin", "*")
		// 指定允许的 HTTP 方法
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		// 指定允许的请求头, 覆盖常见的自定义头(如 authorization)和标准头(如 content-type)
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		// 指定服务器支持的 HTTP 方法
		c.Header("Allow", "HEAD, GET, POST, PUT, PATCH, DELETE, OPTIONS")
		// 设置响应内容类型为 JSON
		c.Header("Content-Type", "application/json")
		// 对于预检请求, 返回状态码 200 并终止请求处理, 不继续执行后续中间件或路由逻辑
		c.AbortWithStatus(http.StatusOK)
		return
	}
	c.Next() // 如果不是预检请求, 调用 c.Next(), 将请求传递给下一个中间件或路由处理程序
}

// gin中间件, 用于添加安全相关的http头
func Secure(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	// 禁止页面被嵌入到 <frame>, <iframe> 或 <object> 中, 防止**点击劫持(Clickjacking)**攻击
	// DENY 表示任何情况下都不允许嵌入
	c.Header("X_Frame-Options", "DENY")
	// 防止浏览器根据内容"嗅探" MIME 类型, 强制使用服务器指定的 Content-Type
	// 避免浏览器将非脚本文件(如 JSON)误解析为可执行脚本, 降低 MIME 类型混淆攻击风险
	c.Header("X-Content-Type-Options", "nosniff")
	// 启用浏览器的内置 XSS 过滤器, 1; mode=block 表示如果检测到潜在的 XSS 攻击, 浏览器将阻止页面加载
	// 现代浏览器已废弃此头
	c.Header("X-XSS-Protection", "1; mode=block")
	if c.Request.TLS != nil {
		// Strict-Transport-Security(HSTS)强制浏览器只通过 HTTPS 访问服务器, 防止中间人攻击
		// HSTS 策略有效期为 1 年
		c.Header("Strict-Transport-Security", "max-age=31536000")
	}
	c.Next()
}
