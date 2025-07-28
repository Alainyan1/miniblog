// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

//go:build wireinject
// +build wireinject

package apiserver

import (
	"miniblog/internal/apiserver/biz"
	"miniblog/internal/apiserver/pkg/validation"
	"miniblog/internal/apiserver/store"
	ginmw "miniblog/internal/pkg/middleware/gin"
	"miniblog/internal/pkg/server"

	"github.com/onexstack/onexstack/pkg/authz"

	"github.com/google/wire"
)

// 使用Wire实现依赖注入逻辑, 为InitializeWebServer函数生成依赖关系图, 自动注入需要的组件, 最终构建出一个完整的server.Server实例
// 参数*Config包含了创建server.Server类型实例的所有依赖项, wire.Build告诉wire如何按照依赖关系注入和构建server.Server对象
func InitializeWebServer(*Config) (server.Server, error) {
	wire.Build(
		wire.NewSet(NewWebServer, wire.FieldsOf(new(*Config), "ServerMode")),
		wire.Struct(new(ServerConfig), "*"),
		wire.NewSet(store.ProviderSet, biz.ProviderSet),
		ProvideDB,
		validation.ProviderSet,
		wire.NewSet(
			wire.Struct(new(UserRetriever), "*"),
			wire.Bind(new(ginmw.UserRetriever), new(*UserRetriever)),
		),
		authz.ProviderSet,
	)
	return nil, nil
}
