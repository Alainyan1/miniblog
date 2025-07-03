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

func UserModelToUserV1(userModel *model.UserM) *apiv1.User {
	var protoUser apiv1.User

	_ = core.CopyWithConverters(&protoUser, userModel)

	return &protoUser
}

func UserV1ToUserModel(protoUser *apiv1.User) *model.UserM {
	var userModel model.UserM

	_ = core.CopyWithConverters(&userModel, protoUser)

	return &userModel
}
