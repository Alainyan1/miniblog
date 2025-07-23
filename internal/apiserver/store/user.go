// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package store

import (
	"context"
	"errors"
	"miniblog/internal/apiserver/model"
	"miniblog/internal/pkg/errno"
	"miniblog/internal/pkg/log"

	"github.com/onexstack/onexstack/pkg/store/where"
	"gorm.io/gorm"
)

// 2. 扩展方法: UserExpansion, 用于添加特定的业务逻辑方法.
type UserStore interface {
	Create(ctx context.Context, obj *model.UserM) error
	Update(ctx context.Context, obj *model.UserM) error
	Delete(ctx context.Context, opts *where.Options) error
	Get(ctx context.Context, opts *where.Options) (*model.UserM, error)
	List(ctx context.Context, opts *where.Options) (int64, []*model.UserM, error)

	UserExpansion
}

// 用户操作的附加方法.
type UserExpansion interface{}

type userStore struct {
	store *datastore
}

var _ UserStore = (*userStore)(nil)

func newUserStore(store *datastore) *userStore {
	return &userStore{store: store}
}

func (s *userStore) Create(ctx context.Context, obj *model.UserM) error {
	// 调用DB(ctx)方法尝试从context中获取事务实例, 如果没有则直返回*gorm.DB实例
	// 调用*gorm.DB的Create方法将对象插入到数据库中
	if err := s.store.DB(ctx).Create(&obj).Error; err != nil {
		log.Errorw("Failed to insert user into database", "err", err, "user", obj)
		// 返回了一个自定义的错误errno.ErrDBWrite
		// 直接返回gorm包的错误信息可能会暴露数据库细节, 这里使用errno包封装了错误
		return errno.ErrDBWrite.WithMessage(err.Error())
	}
	return nil
}

func (s *userStore) Update(ctx context.Context, obj *model.UserM) error {
	if err := s.store.DB(ctx).Save(obj).Error; err != nil {
		log.Errorw("Failed to update user in database", "err", err, "user", obj)
		return errno.ErrDBWrite.WithMessage(err.Error())
	}

	return nil
}

func (s *userStore) Delete(ctx context.Context, opts *where.Options) error {
	err := s.store.DB(ctx, opts).Delete(new(model.UserM)).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		log.Errorw("Failed to delete user from database", "err", err, "conditions", opts)
		return errno.ErrDBWrite.WithMessage(err.Error())
	}
	return nil
}

func (s *userStore) Get(ctx context.Context, opts *where.Options) (*model.UserM, error) {
	var obj model.UserM
	if err := s.store.DB(ctx, opts).First(&obj).Error; err != nil {
		log.Errorw("Failed to retrieve user from database", "err", err, "conditions", opts)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errno.ErrUserNotFound
		}
		return nil, errno.ErrDBRead.WithMessage(err.Error())
	}
	return &obj, nil
}

// nolint: nonamedreturns
func (s *userStore) List(ctx context.Context, opts *where.Options) (count int64, ret []*model.UserM, err error) {
	err = s.store.DB(ctx, opts).Order("id desc").Find(&ret).Offset(-1).Limit(-1).Count(&count).Error
	if err != nil {
		log.Errorw("Failed to list users from database", "err", err, "conditions", opts)
		err = errno.ErrDBRead.WithMessage(err.Error())
	}
	return
}
