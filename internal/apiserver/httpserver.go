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

	// 注册REST API路由
	c.InstallRESTAPI(engine)

	httpsrv := server.NewHTTPServer(c.cfg.HTTPOptions, engine)

	return &ginServer{srv: httpsrv}
}

// 注册API路由, 路由的路径和http方法遵守REST规范
func (c *ServerConfig) InstallRESTAPI(engine *gin.Engine) {
	// 注册业务无关的API接口
	InstallGenericAPI(engine)

	// 创建业务处理器
	handler := handler.NewHandler()

	// 注册健康检查接口
	engine.GET("/healthz", handler.Healthz)
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
