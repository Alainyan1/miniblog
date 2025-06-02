package main

import (
	"fmt"
	"miniblog/internal/pkg/errno"
	"miniblog/internal/pkg/errorsx"
)

func main() {
	// Code: 500
	// Reason: InternalError.DBConnection
	// Message: Message: "Something went wrong: DB connection failed"
	errx := errorsx.New(500, "InternalError.DBConnection", "Something went wrong: %s", "DB connection failed")

	// fmt.Println 会调用errx的Error方法, 输出:
	// error: code = 500 reason = InternalError.DBConnection message = Something went wrong: DB connection failed metadata = map[]
	fmt.Println(errx)

	// 给错误添加元数据
	errx.WithMetadata(map[string]string{
		"user_id":    "123",
		"request_id": "abc",
	})

	// 使用KV方法添加元数据
	errx.KV("trace_id", "123abc")

	// 使用WithMessage更新Message字段
	// 更新Message不会影响Code, Reason和Metadata
	errx.WithMessage("Updated message: %s", "retry failed")

	// 打印errx, 元数据和更新后的message会一并输出
	fmt.Println(errx)

	someerr := doSomething()

	fmt.Println(someerr)

	// 调用预定义错误errno.ErrUserNotFound的Is方法, 判断someerr是否属于该类型错误
	// Is方法只会比较err的Code和Reason字段, 这里返回true
	fmt.Println(errno.ErrUserNotFound.Is(someerr))

	// 返回false
	fmt.Println(errno.ErrAddRole.Is(someerr))

}

func doSomething() error {
	return errno.ErrUserNotFound.WithMessage("Can not find user")
}
