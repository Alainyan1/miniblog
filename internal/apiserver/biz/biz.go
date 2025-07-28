// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package biz

// Biz层依赖Store层, 主要用来实现系统中REST资源的各类业务操作, 例如用户资源的增删改查等.
import (
	postv1 "miniblog/internal/apiserver/biz/v1/post"
	userv1 "miniblog/internal/apiserver/biz/v1/user"
	"miniblog/internal/apiserver/store"

	"github.com/onexstack/onexstack/pkg/authz"

	"github.com/google/wire"
)

// ProviderSet 是一个 Wire 的 Provider 集合, 用于声明依赖注入的规则.
// 包含 NewBiz 构造函数, 用于生成 biz 实例.
// wire.Bind 用于将接口 IBiz 与具体实现 *biz 绑定,
// 这样依赖 IBiz 的地方会自动注入 *biz 实例.
var ProviderSet = wire.NewSet(NewBiz, wire.Bind(new(IBiz), new(*biz)))

// 定义了业务层需要实现的方法.
type IBiz interface {
	// 获取用户业务接口
	UserV1() userv1.UserBiz
	// 获取帖子业务接口
	PostV1() postv1.PostBiz
}

type biz struct {
	store store.IStore
	authz *authz.Authz
}

var _ IBiz = (*biz)(nil)

func NewBiz(store store.IStore, authz *authz.Authz) *biz {
	return &biz{store: store, authz: authz}
}

func (b *biz) UserV1() userv1.UserBiz {
	return userv1.New(b.store, b.authz)
}

func (b *biz) PostV1() postv1.PostBiz {
	return postv1.New(b.store)
}
