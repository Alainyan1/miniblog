// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package biz

// Biz层依赖Store层, 主要用来实现系统中REST资源的各类业务操作, 例如用户资源的增删改查等
import (
	postv1 "miniblog/internal/apiserver/biz/v1/post"
	userv1 "miniblog/internal/apiserver/biz/v1/user"
	"miniblog/internal/apiserver/store"
	"miniblog/pkg/auth"
)

// 定义了业务层需要实现的方法
type IBiz interface {
	// 获取用户业务接口
	UserV1() userv1.UserBiz
	// 获取帖子业务接口
	PostV1() postv1.PostBiz
}

type biz struct {
	store store.IStore
	authz *auth.Authz
}

var _ IBiz = (*biz)(nil)

func NewBiz(store store.IStore, authz *auth.Authz) *biz {
	return &biz{store: store, authz: authz}
}

func (b *biz) UserV1() userv1.UserBiz {
	return userv1.New(b.store, b.authz)
}

func (b *biz) PostV1() postv1.PostBiz {
	return postv1.New(b.store)
}
