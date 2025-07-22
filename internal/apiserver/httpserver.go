// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package apiserver

import (
	"context"
	handler "miniblog/internal/apiserver/handler/http"
	"miniblog/internal/pkg/server"
	"net/http"

	mw "miniblog/internal/pkg/middleware/gin"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

type ginServer struct {
	srv server.Server
}

// 确保实现server.Server接口
var _ server.Server = (*ginServer)(nil)

// 初始化一个新的Gin服务器实例
func (c *ServerConfig) NewGinServer() server.Server {
	// 创建gin引擎
	engine := gin.New()

	// 注册全局中间件, 用于恢复 panic, 设置 HTTP 头, 添加请求 ID 等
	engine.Use(gin.Recovery(), mw.NoCache, mw.Cors, mw.Secure, mw.RequestIDMiddleware()) // 注册REST API路由

	c.InstallRESTAPI(engine)

	httpsrv := server.NewHTTPServer(c.cfg.HTTPOptions, c.cfg.TLSOptions, engine)

	return &ginServer{srv: httpsrv}
}

// 注册API路由, 路由的路径和http方法遵守REST规范
func (c *ServerConfig) InstallRESTAPI(engine *gin.Engine) {
	// 注册业务无关的API接口, 例如版本信息等
	InstallGenericAPI(engine)

	// 创建业务处理器
	handler := handler.NewHandler(c.biz, c.val)

	// 注册健康检查接口
	engine.GET("/healthz", handler.Healthz)

	// 注册用户登录和令牌刷新接口
	engine.POST("login", handler.Login)
	engine.PUT("/refresh-token", mw.AuthnMiddleware(c.retriever), handler.RefreshToken)

	// 中间件切片, 用于在请求处理前后执行逻辑, 如JWT认证
	authMiddlewares := []gin.HandlerFunc{mw.AuthnMiddleware(c.retriever), mw.AuthzMiddleware(c.authz)}

	// 注册v1版本API路由分组
	v1 := engine.Group("/v1")
	{
		userv1 := v1.Group("/users")
		{
			// 创建用户不用进行认证和授权
			userv1.POST("", handler.CreateUser)
			// 其余需要进行认证和授权
			userv1.Use(authMiddlewares...)
			userv1.PUT(":userID/change-password", handler.ChangePassword) // 修改用户密码
			userv1.PUT(":userID", handler.UpdateUser)                     // 更新用户信息
			userv1.DELETE(":userID", handler.DeleteUser)                  // 删除用户
			userv1.GET(":userID", handler.GetUser)                        // 查询用户详情
			userv1.GET("", handler.ListUser)                              // 查询用户列表.
		}

		postv1 := v1.Group("/posts", authMiddlewares...)
		{
			postv1.POST("", handler.CreatePost)       // 创建博客
			postv1.PUT(":postID", handler.UpdatePost) // 更新博客
			postv1.DELETE("", handler.DeletePost)     // 删除博客
			postv1.GET(":postID", handler.GetPost)    // 查询博客详情
			postv1.GET("", handler.ListPost)          // 查询博客列表
		}
	}
}

// 注册业务无关的路由
func InstallGenericAPI(engine *gin.Engine) {
	// 注册pprof路由, 用来提供性能调试和优化的API接口
	pprof.Register(engine)

	// 注册404路由
	engine.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, "Page not found")
	})
}

func (s *ginServer) RunOrDie() {
	s.srv.RunOrDie()
}

func (s *ginServer) GracefulStop(ctx context.Context) {
	s.srv.GracefulStop(ctx)
}
