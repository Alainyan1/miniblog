// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package errorsx

import (
	"errors"
	"fmt"
	"net/http"

	httpstatus "github.com/go-kratos/kratos/v2/transport/http/status"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/status"
)

// nolint
type ErrorX struct {
	// Code表示错误的http状态码, 用于与客户端进行交互时标识的错误类型
	Code int `json:"code,omitempty"`

	// Reason表示错误发生原因, 通常为业务错误码, 用于精准定位问题
	Reason string `json:"reason,omitempty"`

	// Message表示简短的错误信息, 通常可以直接暴露给用户
	Message string `json:"message,omitempty"`

	// Metadata用于存储与该错误相关的额外元信息, 可以包含上下文调试信息
	Metadata map[string]string `json:"metadata,omitempty"`
}

// 创建一个新错误.
func New(code int, reason string, format string, args ...any) *ErrorX {
	return &ErrorX{
		Code:    code,
		Reason:  reason,
		Message: fmt.Sprintf(format, args...),
	}
}

// 实现error接口中的Error方法.
func (err *ErrorX) Error() string {
	return fmt.Sprintf("error: code = %d reason = %s message = %s metadata = %v", err.Code, err.Reason, err.Message, err.Metadata)
}

// 设置错误的Message字段.
func (err *ErrorX) WithMessage(format string, args ...any) *ErrorX {
	err.Message = fmt.Sprintf(format, args...)
	return err
}

// 设置原数据.
func (err *ErrorX) WithMetadata(md map[string]string) *ErrorX {
	err.Metadata = md
	return err
}

// 使用kv对设置原数据.
func (err *ErrorX) KV(kvs ...string) *ErrorX {
	// 初始化原数据
	if err.Metadata == nil {
		err.Metadata = make(map[string]string)
	}

	// kv是成对出现的
	for i := 0; i < len(kvs); i += 2 {
		if i+1 < len(kvs) {
			err.Metadata[kvs[i]] = kvs[i+1]
		}
	}
	return err
}

// 返回grpc状态表示.
func (err *ErrorX) GRPCStatus() *status.Status {
	details := errdetails.ErrorInfo{Reason: err.Reason, Metadata: err.Metadata}
	s, _ := status.New(httpstatus.ToGRPCCode(err.Code), err.Message).WithDetails(&details)
	return s
}

// 设置请求id.
func (err *ErrorX) WithRequestID(requestID string) *ErrorX {
	return err.KV("X-Request-ID", requestID)
}

// 如果Code和Reason都相等, 返回true, 否则返回false.
func (err *ErrorX) Is(target error) bool {
	if errx := new(ErrorX); errors.As(target, &errx) {
		return errx.Code == err.Code && errx.Reason == err.Reason
	}

	return false
}

// 返回错误的http代码.
func Code(err error) int {
	if err == nil {
		return http.StatusOK
	}

	return FromError(err).Code
}

// Reason 返回特定错误的原因.
func Reason(err error) string {
	if err == nil {
		return ErrInternal.Reason
	}
	return FromError(err).Reason
}

// 将一个通用的error转为自定义的 *ErrorX类型.
func FromError(err error) *ErrorX {
	if err == nil {
		return nil
	}

	// 如果可以通过errors.As()转为errx, 直接返回
	if errx := new(ErrorX); errors.As(err, &errx) {
		return errx
	}

	// gRPC 的 status.FromError 方法尝试将 error 转换为 gRPC 错误的 status 对象.
	// 如果 err 不能转换为 gRPC 错误, 即不是 gRPC 的 status 错误
	// 则返回一个带有默认值的 ErrorX, 表示是一个未知类型的错误.
	gs, ok := status.FromError(err)
	if !ok {
		return New(ErrInternal.Code, ErrInternal.Reason, err.Error())
	}

	// 如果 err 是 gRPC 的错误类型, 会成功返回一个 gRPC status 对象 gs.
	// 使用 gRPC 状态中的错误代码和消息创建一个 ErrorX.
	ret := New(httpstatus.FromGRPCCode(gs.Code()), ErrInternal.Reason, gs.Message())

	// 遍历 gRPC 错误详情中的所有附加信息（Details）.
	for _, detail := range gs.Details() {
		if typed, ok := detail.(*errdetails.ErrorInfo); ok {
			ret.Reason = typed.GetReason()
			return ret.WithMetadata(typed.GetMetadata())
		}
	}

	return ret
}
