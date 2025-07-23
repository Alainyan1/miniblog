// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package server

import (
	"context"
	"net/http"
)

// 定义所以服务器类型的接口.
type Server interface {
	// RunOrDie 运行服务器, 如果运行失败会退出程序(OrDie的含义所在)
	RunOrDie()
	// 优雅关停, 关停服务器时需要处理 context 的超时时间
	GracefulStop(ctx context.Context)
}

// 从http.Server获取协议名.
func protocolName(server *http.Server) string {
	if server.TLSConfig != nil {
		return "https"
	}
	return "http"
}
