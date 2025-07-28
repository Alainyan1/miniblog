// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package rid_test

import (
	"fmt"

	"miniblog/internal/pkg/rid"
)

// 实例代码, 向开发者展示用法及输出
func ExampleResourceID_String() {
	// 定义一个资源标识符，例如用户资源
	userID := rid.UserID

	// 调用String方法，将ResourceID类型转换为字符串类型
	idString := userID.String()

	// 输出结果
	fmt.Println(idString)

	// Output:
	// user
}
