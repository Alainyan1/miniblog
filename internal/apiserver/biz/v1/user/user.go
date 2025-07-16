// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package user

import (
	"context"
	"miniblog/internal/apiserver/model"
	"miniblog/internal/apiserver/store"
	"miniblog/internal/pkg/contextx"
	"miniblog/internal/pkg/errno"
	"miniblog/internal/pkg/known"
	"miniblog/internal/pkg/log"
	apiv1 "miniblog/pkg/api/apiserver/v1"
	"miniblog/pkg/auth"
	"miniblog/pkg/token"
	"sync"
	"time"

	"miniblog/internal/apiserver/pkg/conversion"

	"github.com/jinzhu/copier"
	"github.com/onexstack/onexstack/pkg/store/where"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserBiz interface {
	Create(ctx context.Context, rq *apiv1.CreateUserRequest) (*apiv1.CreateUserResponse, error)
	Update(ctx context.Context, rq *apiv1.UpdateUserRequest) (*apiv1.UpdateUserResponse, error)
	Delete(ctx context.Context, rq *apiv1.DeleteUserRequest) (*apiv1.DeleteUserResponse, error)
	Get(ctx context.Context, rq *apiv1.GetUserRequest) (*apiv1.GetUserResponse, error)
	List(ctx context.Context, rq *apiv1.ListUserRequest) (*apiv1.ListUserResponse, error)

	UserExpansion
}

// 扩展接口实现了用户登录, Token刷新, 密码修改和差性能示例方法
type UserExpansion interface {
	Login(ctx context.Context, rq *apiv1.LoginRequest) (*apiv1.LoginResponse, error)
	RefreshToken(ctx context.Context, rq *apiv1.RefreshTokenRequest) (*apiv1.RefreshTokenResponse, error)
	ChangePassword(ctx context.Context, rq *apiv1.ChangePasswordRequest) (*apiv1.ChangePasswordResponse, error)
	ListWithBadPerformance(ctx context.Context, rq *apiv1.ListUserRequest) (*apiv1.ListUserResponse, error)
}

type userBiz struct {
	store store.IStore
}

var _ UserBiz = (*userBiz)(nil)

func New(store store.IStore) *userBiz {
	return &userBiz{store: store}
}

func (b *userBiz) Login(ctx context.Context, rq *apiv1.LoginRequest) (*apiv1.LoginResponse, error) {
	// 获取用户登录的所有信息
	whr := where.F("username", rq.GetUsername())

	userM, err := b.store.User().Get(ctx, whr)

	if err != nil {
		return nil, errno.ErrUserNotFound
	}

	// 对比传入的明文密码和数据库中已加密过的密码是否匹配
	// auth.Compare会将传入的明文密码加密然后和数据库中的密码对比
	if err := auth.Compare(userM.Password, rq.GetPassword()); err != nil {
		log.W(ctx).Errorw("Failed to compare password", "err", err)
		return nil, errno.ErrPasswordInvalid
	}

	// 实现Token签发逻辑, 在签发token时会在token的payload中保存用户id
	tokenStr, expireAt, err := token.Sign(userM.UserID)
	if err != nil {
		return nil, errno.ErrSignToken
	}

	return &apiv1.LoginResponse{Token: tokenStr, ExpireAt: timestamppb.New(expireAt)}, nil
}

// 刷新用户的身份验证令牌
func (b *userBiz) RefreshToken(ctx context.Context, rq *apiv1.RefreshTokenRequest) (*apiv1.RefreshTokenResponse, error) {
	// TODO: 实现Token签发逻辑
	return &apiv1.RefreshTokenResponse{Token: "<placeholder>", ExpireAt: timestamppb.New(time.Now().Add(2 * time.Hour))}, nil
}

// 修改用户密码
// userM结构体实现了BeforeCreate的钩子, 在用户创建记录前, 会将明文密码加密后保存
// 更新用户时, 不会调用BeforeUpdate钩子, 因此需要在修改密码时手动加密新密码
func (b *userBiz) ChangePassword(ctx context.Context, rq *apiv1.ChangePasswordRequest) (*apiv1.ChangePasswordResponse, error) {
	// where.T方法用于构造查询条件
	// T函数自动为查询添加租户的隔离条件
	// ???一个云端CRM(客户关系管理)系统可能为多个公司(租户)提供服务, 每个租户只能访问自己的客户数据
	userM, err := b.store.User().Get(ctx, where.T(ctx))

	if err != nil {
		return nil, err
	}

	// 验证旧密码是否正确
	if err := auth.Compare(userM.Password, rq.GetOldPassword()); err != nil {
		log.W(ctx).Errorw("Failed to compare password", "err", err)
		return nil, errno.ErrPasswordInvalid
	}

	// 如果旧密码验证通过, 则加密新密码
	userM.Password, _ = auth.Encrypt(rq.GetNewPassword())
	if err := b.store.User().Update(ctx, userM); err != nil {
		return nil, err
	}

	return &apiv1.ChangePasswordResponse{}, nil
}

func (b *userBiz) Create(ctx context.Context, rq *apiv1.CreateUserRequest) (*apiv1.CreateUserResponse, error) {
	var userM model.UserM
	// 使用copier的Copy函数给目标结构体变量userM赋值
	_ = copier.Copy(&userM, rq)

	// b.store.User().Create(ctx, &userM)将用户保存在数据库中
	if err := b.store.User().Create(ctx, &userM); err != nil {
		return nil, err
	}

	return &apiv1.CreateUserResponse{UserID: userM.UserID}, nil
}

func (b *userBiz) Update(ctx context.Context, rq *apiv1.UpdateUserRequest) (*apiv1.UpdateUserResponse, error) {
	userM, err := b.store.User().Get(ctx, where.T(ctx))
	if err != nil {
		return nil, err
	}

	// Username是*string类型
	if rq.Username != nil {
		userM.Username = rq.GetUsername()
	}
	if rq.Email != nil {
		userM.Email = rq.GetEmail()
	}
	if rq.Nickname != nil {
		userM.Nickname = rq.GetNickname()
	}
	if rq.Phone != nil {
		userM.Phone = rq.GetPhone()
	}

	if err := b.store.User().Update(ctx, userM); err != nil {
		return nil, err
	}

	return &apiv1.UpdateUserResponse{}, nil

}

func (b *userBiz) Delete(ctx context.Context, rq *apiv1.DeleteUserRequest) (*apiv1.DeleteUserResponse, error) {
	// 只有root用户可以删除用户
	// 这里不用where.T()因为where.T()会查询root自己
	// 因为where.T()会添加条件, 只会针对特定的数据进行查询
	if err := b.store.User().Delete(ctx, where.F("userID", rq.GetUserID())); err != nil {
		return nil, err
	}

	return &apiv1.DeleteUserResponse{}, nil
}

func (b *userBiz) Get(ctx context.Context, rq *apiv1.GetUserRequest) (*apiv1.GetUserResponse, error) {
	userM, err := b.store.User().Get(ctx, where.T(ctx))
	if err != nil {
		return nil, err
	}
	return &apiv1.GetUserResponse{User: conversion.UserModelToUserV1(userM)}, nil
}

func (b *userBiz) List(ctx context.Context, rq *apiv1.ListUserRequest) (*apiv1.ListUserResponse, error) {
	whr := where.P(int(rq.GetOffset()), int(rq.GetLimit()))
	if contextx.Username(ctx) != known.AdminUsername {
		whr.T(ctx)
	}
	count, userList, err := b.store.User().List(ctx, whr)
	if err != nil {
		return nil, err
	}

	// 线程安全的map存储并发查询结果
	var m sync.Map
	// 创建errgroup管理并发执行的goroutine
	// 普通的goroutine没有同步机制, 可能导致goroutine泄漏或资源竞争
	// 纯 WaitGroup和goroutine都需要额外的错误处理和📱机制
	eg, ctx := errgroup.WithContext(ctx)

	// 设置并发限制
	eg.SetLimit(known.MaxErrGroupConcurrency)

	// 并发统计用户的博客数
	for _, user := range userList {
		eg.Go(func() error {
			select {
			case <-ctx.Done():
				return nil
			default:
				count, _, err := b.store.Post().List(ctx, where.T(ctx))
				if err != nil {
					return err
				}
				// 讲存储层用户模型userM转换为API层的用户对象(apiv1.User)
				converted := conversion.UserModelToUserV1(user)
				converted.PostCount = count
				// 将转换后的用户对象存储到syncMap中, id为key
				m.Store(user.ID, converted)

				return nil
			}
		})
	}

	// 等待并发任务完成
	if err := eg.Wait(); err != nil {
		log.W(ctx).Errorw("Failed to wait all function calls returned", "err", err)
		return nil, err
	}

	// 构建Response
	users := make([]*apiv1.User, 0, len(userList))
	for _, item := range userList {
		user, _ := m.Load(item.ID)
		users = append(users, user.(*apiv1.User))
	}

	log.W(ctx).Debugw("Get users from backend storage", "count", len(users))

	return &apiv1.ListUserResponse{TotalCount: count, Users: users}, nil
}

// ListWithBadPerformance 是性能较差的实现方式(已废弃).
func (b *userBiz) ListWithBadPerformance(ctx context.Context, rq *apiv1.ListUserRequest) (*apiv1.ListUserResponse, error) {
	whr := where.P(int(rq.GetOffset()), int(rq.GetLimit()))
	if contextx.Username(ctx) != known.AdminUsername {
		whr.T(ctx)
	}

	count, userList, err := b.store.User().List(ctx, whr)
	if err != nil {
		return nil, err
	}

	users := make([]*apiv1.User, 0, len(userList))
	for _, user := range userList {
		count, _, err := b.store.Post().List(ctx, where.T(ctx))
		if err != nil {
			return nil, err
		}

		converted := conversion.UserModelToUserV1(user)
		converted.PostCount = count
		users = append(users, converted)
	}

	log.W(ctx).Debugw("Get users from backend storage", "count", len(users))

	return &apiv1.ListUserResponse{TotalCount: count, Users: users}, nil
}
