// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package model

import (
	"miniblog/internal/pkg/rid"

	"gorm.io/gorm"
)

// 实现钩子函数, 用于在模型操作前后执行特定逻辑

// 在创建数据库记录后生成postID
func (m *PostM) AfterCreate(tx *gorm.DB) error {
	// 基于数据库生成的自增ID生成一个形如user-uvalgf的唯一id, 并调用tx.Save()方法将ID更新到表记录中
	m.PostID = rid.PostID.New(uint64(m.ID))
	return tx.Save(m).Error
}

// 在创建数据库记录后生成userID
func (m *UserM) AfterCreate(tx *gorm.DB) error {
	m.UserID = rid.UserID.New(uint64(m.ID))
	return tx.Save(m).Error
}
