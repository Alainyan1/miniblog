// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package validation

import (
	"context"
	"miniblog/internal/pkg/errno"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	apiv1 "miniblog/pkg/api/apiserver/v1"

	genericvalidation "github.com/onexstack/onexstack/pkg/validation"
)

func (v *Validator) ValidatePostRules() genericvalidation.Rules {
	// 定义各字段的校验逻辑, 通过一个 map 实现模块化和简化
	// 对于ValidateAllFields和ValidateSelectedFields函数, 如果结构体中某个字段不存在对应的Rules, 回跳过该字段的校验
	return genericvalidation.Rules{
		"PostID": func(value any) error {
			if value.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("postID cannot be empty")
			}
			return nil
		},
		"Title": func(value any) error {
			if value.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("title cannot be empty")
			}
			return nil
		},
		"Content": func(value any) error {
			if value.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("content cannot be empty")
			}
			return nil
		},
	}
}

// ValidateAllFields对请求参数中的所有字段进行校验, 每个字段的校验规则在ValidatePostRules中设置
// ValidateCreatePostRequest 校验 CreatePostRequest 结构体的有效性.
func (v *Validator) ValidateCreatePostRequest(ctx context.Context, rq *apiv1.CreatePostRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidatePostRules())
}

// ValidateUpdatePostRequest 校验更新用户请求.
func (v *Validator) ValidateUpdatePostRequest(ctx context.Context, rq *apiv1.UpdatePostRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidatePostRules())
}

// ValidateDeletePostRequest 校验 DeletePostRequest 结构体的有效性.
func (v *Validator) ValidateDeletePostRequest(ctx context.Context, rq *apiv1.DeletePostRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidatePostRules())
}

// ValidateGetPostRequest 校验 GetPostRequest 结构体的有效性.
func (v *Validator) ValidateGetPostRequest(ctx context.Context, rq *apiv1.GetPostRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidatePostRules())
}

// ValidateListPostRequest 校验 ListPostRequest 结构体的有效性.
// 调用了ValidateSelectedFields函数, 只会校验传入的Offset, Limit
// 如果指定了不存在的字段, 则会跳过该字段的检验
// 对于其他字段, 可以自行实现校验逻辑, 可以根据需求选择哪些使用通用的校验规则, 哪些使用自行实现的校验规则
func (v *Validator) ValidateListPostRequest(ctx context.Context, rq *apiv1.ListPostRequest) error {
	if err := validation.Validate(rq.GetTitle(), validation.Length(5, 100), is.URL); err != nil {
		return errno.ErrInvalidArgument.WithMessage(err.Error())
	}
	return genericvalidation.ValidateSelectedFields(rq, v.ValidatePostRules(), "Offset", "Limit")
}
