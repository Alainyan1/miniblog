// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package conversion

import (
	"miniblog/internal/apiserver/model"

	"github.com/onexstack/onexstack/pkg/core"

	apiv1 "miniblog/pkg/api/apiserver/v1"
)

// 将模型层的PostM转换为Protobuf层的Post(v1 博客对象)
func PostModelToPostV1(postModel *model.PostM) *apiv1.Post {
	var protoPost apiv1.Post
	_ = core.CopyWithConverters(&protoPost, postModel)
	return &protoPost
}

// 将 Protobuf 层的 Post(v1 博客对象)转换为模型层的 PostM(博客模型对象)
func PostV1ToPostModel(protoPost *apiv1.Post) *model.PostM {
	var postModel model.PostM
	_ = core.CopyWithConverters(&postModel, protoPost)
	return &postModel
}
