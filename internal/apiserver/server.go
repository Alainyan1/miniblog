// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package apiserver

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/viper"
)

// 基于初始化配置创建运行时配置

// 存储应用相关配置
type Config struct {
	ServerMode string
	JWTKey     string
	Expiration time.Duration
}

// 定义一个联合服务器, 根据ServerMode决定要启动的服务器类型
type UnionServer struct {
	cfg *Config
}

// 根据配置创建联合服务器
func (cfg *Config) NewUnionServer() (*UnionServer, error) {
	return &UnionServer{cfg: cfg}, nil
}

// Run运行应用
func (s *UnionServer) Run() error {
	fmt.Printf("ServerMode from ServerOptions: %s\n", s.cfg.JWTKey)
	fmt.Printf("ServerMode from Viper: %s\n\n", viper.GetString("jwt-key"))

	jsonData, _ := json.MarshalIndent(s.cfg, "", " ")
	fmt.Println(jsonData)

	select {}
	return nil
}
