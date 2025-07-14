// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package validation

import (
	"miniblog/internal/apiserver/store"
	"miniblog/internal/pkg/errno"
	"regexp"
)

// 验证逻辑的实现结构体
type Validator struct {
	// 有些复杂的验证逻辑, 可能需要直接查询数据库
	// 这里只是一个举例, 如果验证时, 有其他依赖的客户端/服务/资源等
	// 都可以一并注入进来
	store store.IStore
}

func New(store store.IStore) *Validator {
	return &Validator{store: store}
}

func isValidUsername(username string) bool {
	var (
		lengthRegex = `^.{3,20}$`       // 长度在3到20个字符
		validRegex  = `^[A-Za-z0-9_]+$` // 仅包含字母, 数字和下划线
	)

	// 校验长度
	if matched, _ := regexp.MatchString(lengthRegex, username); !matched {
		return false
	}

	// 校验合法性
	if matched, _ := regexp.MatchString(validRegex, username); !matched {
		return false
	}

	return true
}

func isValidPassword(password string) error {
	// 不为空
	if password == "" {
		return errno.ErrInvalidArgument.WithMessage("password cannot be empty")
	}

	// 长度大于6
	if len(password) < 6 {
		return errno.ErrInvalidArgument.WithMessage("password must be at least 6 characters long")
	}

	letterPattern := regexp.MustCompile(`[A-Za-z]`)
	if !letterPattern.MatchString(password) {
		return errno.ErrInvalidArgument.WithMessage("password must contain at least one letter")
	}

	numberPattern := regexp.MustCompile(`\d`)
	if !numberPattern.MatchString(password) {
		return errno.ErrInvalidArgument.WithMessage("password must contain at least one number")
	}

	return nil
}

// isValidEmail 判断电子邮件是否合法.
func isValidEmail(email string) error {
	// 检查电子邮件地址格式
	if email == "" {
		return errno.ErrInvalidArgument.WithMessage("email cannot be empty")
	}

	// 使用正则表达式校验电子邮件格式
	emailPattern := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailPattern.MatchString(email) {
		return errno.ErrInvalidArgument.WithMessage("invalid email format")
	}

	return nil
}

// isValidPhone 判断手机号码是否合法.
func isValidPhone(phone string) error {
	// 检查手机号码格式
	if phone == "" {
		return errno.ErrInvalidArgument.WithMessage("phone cannot be empty")
	}

	// 使用正则表达式校验手机号码格式（假设是中国手机号，11位数字）
	phonePattern := regexp.MustCompile(`^1[3-9]\d{9}$`)
	if !phonePattern.MatchString(phone) {
		return errno.ErrInvalidArgument.WithMessage("invalid phone format")
	}

	return nil
}
