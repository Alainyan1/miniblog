// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// 测试原生 SQL 插入
	dsn := "miniblog:miniblog1234@tcp(127.0.0.1:3306)/miniblog?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// 测试1：使用 DEFAULT 值
	fmt.Println("=== Test 1: Using DEFAULT ===")
	_, err = db.Exec(`
        INSERT INTO user (userID, username, password, nickname, email, phone) 
        VALUES (?, ?, ?, ?, ?, ?)
    `, "test1", "testuser1", "pass", "nick", "test1@example.com", "13800000001")

	if err != nil {
		fmt.Printf("❌ DEFAULT insert failed: %v\n", err)
	} else {
		fmt.Printf("✅ DEFAULT insert succeeded\n")

		// 查询结果
		var createdAt time.Time
		err = db.QueryRow("SELECT createdAt FROM user WHERE userID = ?", "test1").Scan(&createdAt)
		if err == nil {
			fmt.Printf("   CreatedAt: %v\n", createdAt)
		}

		// 清理
		db.Exec("DELETE FROM user WHERE userID = ?", "test1")
	}

	// 测试2：显式设置时间
	fmt.Println("\n=== Test 2: Explicit time ===")
	now := time.Now()
	_, err = db.Exec(`
        INSERT INTO user (userID, username, password, nickname, email, phone, createdAt, updatedAt) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `, "test2", "testuser2", "pass", "nick", "test2@example.com", "13800000002", now, now)

	if err != nil {
		fmt.Printf("❌ Explicit insert failed: %v\n", err)
	} else {
		fmt.Printf("✅ Explicit insert succeeded\n")
		db.Exec("DELETE FROM user WHERE userID = ?", "test2")
	}

	// 测试3：零值时间
	fmt.Println("\n=== Test 3: Zero time ===")
	zeroTime := time.Time{}
	fmt.Printf("Zero time: %v (IsZero: %v)\n", zeroTime, zeroTime.IsZero())

	_, err = db.Exec(`
        INSERT INTO user (userID, username, password, nickname, email, phone, createdAt, updatedAt) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `, "test3", "testuser3", "pass", "nick", "test3@example.com", "13800000003", zeroTime, zeroTime)

	if err != nil {
		fmt.Printf("❌ Zero time insert failed: %v\n", err)
	} else {
		fmt.Printf("✅ Zero time insert succeeded\n")
		db.Exec("DELETE FROM user WHERE userID = ?", "test3")
	}
}
