// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package apiserver

import (
	"context"
	"miniblog/internal/pkg/server"
)

type ginServer struct{}

var _ server.Server = (*ginServer)(nil)

func (c *ServerConfig) NewGinServer() server.Server {
	return &ginServer{}
}

func (s *ginServer) RunOrDie() {
	select {}
}

func (s *ginServer) GracefulStop(ctx context.Context) {}
