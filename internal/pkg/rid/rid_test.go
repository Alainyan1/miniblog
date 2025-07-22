// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package rid_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"miniblog/internal/pkg/rid"
)

// Mock Salt function used for testing
func Salt() string {
	return "staticSalt"
}

func TestResourceID_String(t *testing.T) {
	// 测试 UserID 转换为字符串
	userID := rid.UserID
	assert.Equal(t, "user", userID.String(), "UserID.String() should return 'user'")

	// 测试 PostID 转换为字符串
	postID := rid.PostID
	assert.Equal(t, "post", postID.String(), "PostID.String() should return 'post'")
}

func TestResourceID_New(t *testing.T) {
	// 测试生成的ID是否带有正确前缀
	userID := rid.UserID
	uniqueID := userID.New(1)

	assert.True(t, len(uniqueID) > 0, "Generated ID should not be empty")
	assert.Contains(t, uniqueID, "user-", "Generated ID should start with 'user-' prefix")

	// 生成另外一个唯一标识符，确保生成的值不同（唯一性）
	anotherID := userID.New(2)
	assert.NotEqual(t, uniqueID, anotherID, "Generated IDs should be unique")
}

func BenchmarkResourceID_New(b *testing.B) {
	// 性能测试, 必须以Benchmark开头, 运行测试函数需要指定参数-test.bench来运行
	b.ResetTimer() // 重置计时器, 进行一些耗时的准备工作
	// 测量 userID.New 方法的性能表现
	for i := 0; i < b.N; i++ {
		userID := rid.UserID
		_ = userID.New(uint64(i))
	}
}

// 模糊测试, 一种自动化测试技术, 通过不断生成输入数据来测试程序中的漏洞
func FuzzResourceID_New(f *testing.F) {
	// 添加预置测试数据, 使用f.Add方法告诉Fuzzing引擎所需要的数据类型和顺序
	f.Add(uint64(1))      // 添加一个种子值counter为1
	f.Add(uint64(123456)) // 添加一个较大的种子值

	f.Fuzz(func(t *testing.T, counter uint64) {
		// 测试UserID的New方法
		result := rid.UserID.New(counter)

		// 断言结果不为空
		assert.NotEmpty(t, result, "The generated unique identifier should not be empty")

		// 断言结果必须包含资源标识符前缀
		assert.Contains(t, result, rid.UserID.String()+"-", "The generated unique identifier should contain the correct prefix")

		// 断言前缀不会与uniqueStr部分重叠
		splitParts := strings.SplitN(result, "-", 2)
		assert.Equal(t, rid.UserID.String(), splitParts[0], "The prefix part of the result should correctly match the UserID")

		// 断言生成的ID具有固定长度（基于NewCode的配置）
		if len(splitParts) == 2 {
			assert.Equal(t, 6, len(splitParts[1]), "The unique identifier part should have a length of 6")
		} else {
			t.Errorf("The format of the generated unique identifier does not meet expectation")
		}
	})
}
