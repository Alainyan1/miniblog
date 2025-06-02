package errno

import (
	"miniblog/internal/pkg/errorsx"
	"net/http"
)

var ErrPostNotFound = &errorsx.ErrorX{Code: http.StatusNotFound, Reason: "NotFound.PostNotFound", Message: "Post not found."}
