// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package validation

import (
	"context"
	"miniblog/internal/pkg/contextx"
	"miniblog/internal/pkg/errno"

	apiv1 "miniblog/pkg/api/apiserver/v1"

	genericvalidation "github.com/onexstack/onexstack/pkg/validation"
)

func (v *Validator) ValidateUserRules() genericvalidation.Rules {
	validatePassword := func() genericvalidation.ValidatorFunc {
		return func(value any) error {
			return isValidPassword(value.(string))
		}
	}

	return genericvalidation.Rules{
		"Password":    validatePassword(),
		"OldPassword": validatePassword(),
		"NewPassword": validatePassword(),
		"UserID": func(value any) error {
			if value.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("userID cannot be empty")
			}
			return nil
		},
		"Username": func(value any) error {
			if !isValidUsername(value.(string)) {
				return errno.ErrUsernameInvalid
			}
			return nil
		},
		"Nickname": func(value any) error {
			if len(value.(string)) >= 30 {
				return errno.ErrInvalidArgument.WithMessage("nickname must be less than 30 characters")
			}
			return nil
		},
		"Email": func(value any) error {
			return isValidEmail(value.(string))
		},
		"Phone": func(value any) error {
			return isValidPhone(value.(string))
		},
		"Limit": func(value any) error {
			if value.(int64) <= 0 {
				return errno.ErrInvalidArgument.WithMessage("limit must be greater than 0")
			}
			return nil
		},
		"Offset": func(value any) error {
			return nil
		},
	}
}

func (v *Validator) ValidateLoginRequest(ctx context.Context, rq *apiv1.LoginRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateUserRules())
}

func (v *Validator) ValidateChangePasswordRequest(ctx context.Context, rq *apiv1.ChangePasswordRequest) error {
	if rq.GetUserID() != contextx.UserID(ctx) {
		return errno.ErrPermissionDenied.WithMessage("The logged-in user `%s` does not match request user `%s`", contextx.UserID(ctx), rq.GetUserID())
	}
	return genericvalidation.ValidateAllFields(rq, v.ValidateUserRules())
}

func (v *Validator) ValidateCreateUserRequest(ctx context.Context, rq *apiv1.CreateUserRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateUserRules())
}

// ValidateUpdateUserRequest 校验更新用户请求.
func (v *Validator) ValidateUpdateUserRequest(ctx context.Context, rq *apiv1.UpdateUserRequest) error {
	if rq.GetUserID() != contextx.UserID(ctx) {
		return errno.ErrPermissionDenied.WithMessage("The logged-in user `%s` does not match request user `%s`", contextx.UserID(ctx), rq.GetUserID())
	}
	return genericvalidation.ValidateSelectedFields(rq, v.ValidateUserRules(), "UserID")
}

// ValidateDeleteUserRequest 校验 DeleteUserRequest 结构体的有效性.
func (v *Validator) ValidateDeleteUserRequest(ctx context.Context, rq *apiv1.DeleteUserRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateUserRules())
}

// ValidateGetUserRequest 校验 GetUserRequest 结构体的有效性.
func (v *Validator) ValidateGetUserRequest(ctx context.Context, rq *apiv1.GetUserRequest) error {
	if rq.GetUserID() != contextx.UserID(ctx) {
		return errno.ErrPermissionDenied.WithMessage("The logged-in user `%s` does not match request user `%s`", contextx.UserID(ctx), rq.GetUserID())
	}
	return genericvalidation.ValidateAllFields(rq, v.ValidateUserRules())
}

// ValidateListUserRequest 校验 ListUserRequest 结构体的有效性.
func (v *Validator) ValidateListUserRequest(ctx context.Context, rq *apiv1.ListUserRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateUserRules())
}
