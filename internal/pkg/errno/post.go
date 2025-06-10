// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package errno

import (
	"miniblog/internal/pkg/errorsx"
	"net/http"
)

var ErrPostNotFound = &errorsx.ErrorX{Code: http.StatusNotFound, Reason: "NotFound.PostNotFound", Message: "Post not found."}
