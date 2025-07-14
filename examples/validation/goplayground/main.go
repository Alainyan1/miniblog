// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package main

import (
	"fmt"

	"github.com/go-playground/validator"
)

type LoginRequest struct {
	Username string `validate:"required"`       // 必填字段, 结构体标签, 为结构体字段添加元数据
	Password string `validate:"required,min=6"` // 最小长度为6
	Email    string `validate:"required,email"` // 必填且必须是邮箱格式
}

func main() {
	validate := validator.New()

	req := LoginRequest{
		Username: "user",
		Password: "12345",
		Email:    "invalid-email",
	}

	err := validate.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Printf("Field '%s' failed validation, rule '%s'\n", err.Field(), err.Tag())
		}
	} else {
		fmt.Println("Validation passed!")
	}
}
