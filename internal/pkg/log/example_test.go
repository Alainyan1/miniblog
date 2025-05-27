// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package log_test

import (
	"miniblog/internal/pkg/log"
	"testing"
	"time"
)

func TestLogger(t *testing.T) {
	opts := &log.Options{
		Level:             "debug",
		Format:            "json",
		DisableCaller:     false,
		DisableStacktrace: false,
		OutputPaths:       []string{"stdout"},
	}

	// 初始化全局日志
	log.Init(opts)

	log.Debugw("This is a debug message", "key1", "value1")
	log.Infow("This is a info message", "key2", 123)
	log.Warnw("This is a warning message", "current_time", time.Now())
	log.Errorw("This is an error message", "error", "something went wrong", "current_time", time.Now())

	// Panicw 和 Fatalw 会中断程序运行, 因此在测试中应小心使用
	// 在单独的环境中运行该测试
	// log.Panicw("This is a panic message", "reason", "unexpected situation")
	// log.Fatalw("This is a fatal message", "reason", "critical failure")

	// 确保日志缓冲区被刷新
	log.Sync()
}
