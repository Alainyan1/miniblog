// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package store

// 统一的错误处理和日志: 如果业务层直接操作数据库, 错误处理和日志记录会分散, 难以统一管理.
import (
	"context"
	"sync"

	"github.com/onexstack/onexstack/pkg/store/where" // 封装查询逻辑
	"gorm.io/gorm"
)

var (
	once sync.Once
	// 全局变量, 方便其他包调用已经初始化好的datastore实例.
	S *datastore
)

// 定义Store层需要实现的方法.
type IStore interface {
	// 返回Store层的*gorm.DB实例
	DB(ctx context.Context, wheres ...where.Where) *gorm.DB
	TX(ctx context.Context, fn func(ctx context.Context) error) error

	User() UserStore
	Post() PostStore
}

// 用于在context.Context中存储事务的上下文键.
type transactionKey struct{}

// datastore是IStore的具体实现.
type datastore struct {
	core *gorm.DB
	// 可以根据需要添加其他数据库实例
	// fake *gorm.DB
}

var _ IStore = (*datastore)(nil)

// 工厂函数, 创建一个IStore类型的实例.
func NewStore(db *gorm.DB) *datastore {
	// 单例模式保证全局共享一个数据库连接池, 减少资源开销, 同时方便其他模块直接访问 store.S
	once.Do(func() {
		S = &datastore{db}
	})
	return S
}

// 如果未传入任何条件, 则返回上下文中的数据库实例(事务实例或核心数据库实例).
func (store *datastore) DB(ctx context.Context, wheres ...where.Where) *gorm.DB {
	db := store.core
	// 从上下文中提取事务实例
	if tx, ok := ctx.Value(transactionKey{}).(*gorm.DB); ok {
		db = tx
	}
	// 遍历所有传入的条件并逐一叠加到数据库查询对象上
	for _, whr := range wheres {
		db = whr.Where(db)
	}
	return db
}

// 4. 如果fn返回错误, 事务会自动会滚, 否则事务提交.
func (store *datastore) TX(ctx context.Context, fn func(ctx context.Context) error) error {
	return store.core.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			ctx := context.WithValue(ctx, transactionKey{}, tx)
			return fn(ctx)
		},
	)
}

// 返回一个实现了UserStore接口的实例.
func (store *datastore) User() UserStore {
	return newUserStore(store)
}

// 返回一个实现了PostStore接口的实例.
func (store *datastore) Post() PostStore {
	return newPostStore(store)
}

// 使用示例
// 初始化 Store
// db, _ := gorm.Open(mysql.Open("dsn"), &gorm.Config{})
// store := store.NewStore(db)

// // 查询用户
// ctx := context.Background()
// userDB := store.DB(ctx, where.Eq("id", 1), where.Eq("status", "active"))
// var user User
// userDB.First(&user)

// // 事务操作
// err := store.TX(ctx, func(ctx context Comma separated list of conditions for filtering the query

// System: 结果) error {
// 	// 在事务中更新用户状态
// 	userDB := store.DB(ctx)
// 	return userDB.Model(&user).Update("status", "inactive").Error
// })
// if err != nil {
// 	fmt.Println("事务失败:", err)
// }
