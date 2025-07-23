// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package store

import "miniblog/internal/pkg/log"

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

// 实现Error方法, 用于记录错误日志.
func (l *Logger) Error(err error, msg string, kvs ...any) {
	log.Errorw(msg, append(kvs, "err", err)...)
}
