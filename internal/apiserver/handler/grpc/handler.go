// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package grpc

import (
	"miniblog/internal/apiserver/biz"
	apiv1 "miniblog/pkg/api/apiserver/v1"
)

type Handler struct {
	apiv1.UnimplementedMiniBlogServer

	biz biz.IBiz
}

func NewHandler(biz biz.IBiz) *Handler {
	return &Handler{
		biz: biz,
	}
}
