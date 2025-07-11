// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package log

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOptions(t *testing.T) {
	opts := NewOptions()

	// 验证Options的默认值
	assert.NotNil(t, opts, "Options should not be nil")
	assert.Equal(t, false, opts.DisableCaller, "DisableCaller should be false by default")
	assert.Equal(t, false, opts.DisableStacktrace, "DisableStacktrace should be false by default")
	assert.Equal(t, "info", opts.Level, "Level should be 'info' by default")
	assert.Equal(t, "console", opts.Format, "Format should be 'console' by default")
	assert.Equal(t, []string{"stdout"}, opts.OutputPaths, "OutputPaths should be ['stdout'] by default")
}
