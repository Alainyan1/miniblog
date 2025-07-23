// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package server

import (
	"context"
	"crypto/tls"
	"errors"
	"miniblog/internal/pkg/log"
	"net/http"

	genericoptions "github.com/onexstack/onexstack/pkg/options"
)

// 代表有一个http服务器.
type HTTPServer struct {
	srv *http.Server
}

// 创建一个新的服务器实例.
func NewHTTPServer(httpOptions *genericoptions.HTTPOptions, tlsOptions *genericoptions.TLSOptions, handler http.Handler) *HTTPServer {
	var tlsConfig *tls.Config
	if tlsOptions != nil && tlsOptions.UseTLS {
		tlsConfig = tlsOptions.MustTLSConfig()
	}
	return &HTTPServer{
		srv: &http.Server{
			Addr:      httpOptions.Addr,
			Handler:   handler,
			TLSConfig: tlsConfig,
		},
	}
}

// 启动服务器.
func (s *HTTPServer) RunOrDie() {
	log.Infow("Start to listen the incoming requests", "protocol", protocolName(s.srv), "addr", s.srv.Addr)
	serveFn := func() error { return s.srv.ListenAndServe() }

	if s.srv.TLSConfig != nil {
		serveFn = func() error { return s.srv.ListenAndServeTLS("", "") }
	}
	if err := serveFn(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalw("Failed to server HTTP(s) server", "err", err)
	}
}

// 优雅关停.
func (s *HTTPServer) GracefulStop(ctx context.Context) {
	log.Infow("Gracefully stop HTTP(s) server")
	// http.Server的shutdown工作流程如下:
	// 首先关闭所有已开启的监听器, 然后关闭所有的空闲连接, 等待所有活跃连接进入空闲状态后终止服务
	// 如果传入的ctx在服务完成终止前超时, 则shutdown方法会返回context相关的错误, 否则会返回由关闭服务监听器引发的其他错误
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Errorw("HTTP(s) server forced to shutdown", "err", err)
	}
}
