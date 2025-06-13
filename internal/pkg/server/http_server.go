// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package server

import (
	"context"
	"errors"
	"net/http"

	"miniblog/internal/pkg/log"

	genericoptions "github.com/onexstack/onexstack/pkg/options"
)

// 代表有一个http服务器
type HTTPServer struct {
	srv *http.Server
}

// 创建一个新的服务器实例
func NewHTTPServer(httpOptions *genericoptions.HTTPOptions, handler http.Handler) *HTTPServer {
	return &HTTPServer{
		srv: &http.Server{
			Addr:    httpOptions.Addr,
			Handler: handler,
		},
	}
}

// 启动服务器
func (s *HTTPServer) RunOrDie() {
	log.Infow("Start to listen the incoming requests", "protocol", protocolName(s.srv), "addr", s.srv.Addr)
	if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalw("Failed to server HTTP(s) server", "err", err)
	}
}

// 优雅关停
func (s *HTTPServer) GracefulStop(ctx context.Context) {
	log.Infow("Gracefully stop HTTP(s) server")
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Errorw("HTTP(s) server forced to shutdown", "err", err)
	}
}
