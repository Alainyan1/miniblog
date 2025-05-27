// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 用于测试自定义的Logger
var mockLogger *zapLogger

// 初始化测试环境
func TestMain(m *testing.M) {
	opts := &Options{
		Level:             "debug",
		DisableCaller:     false,
		DisableStacktrace: false,
		Format:            "json",
		OutputPaths:       []string{"stdout"},
	}
	Init(opts)
	mockLogger = std
	m.Run()
}

// 测试日志记录方法
func TestLoggerMethods(t *testing.T) {
	assert.NotPanics(t, func() {
		Debugw("debug message", "key1", "value1")
		Infow("info message", "key2", "value2")
		Warnw("warn message", "key3", "value3")
		Errorw("error message", "key4", "value4")
	}, "Log methods should not cause a panic in this test")
	assert.Panics(t, func() {
		Panicw("panic message", "key5", "value5")
	}, "Panicw should cause a panic and exit the program")
}

// 测试初始化
func TestLoggerInitialization(t *testing.T) {
	opts := NewOptions()
	logger := New(opts)

	assert.NotNil(t, logger, "logger should not be nil after initialization")
	assert.IsType(t, &zapLogger{}, "logger should be of type *zapLogger")
}

// 测试日志同步
func TestSync(t *testing.T) {
	assert.NotPanics(t, func() {
		Sync()
	}, "Sync should not panic")
}
