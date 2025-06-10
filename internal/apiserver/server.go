// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package apiserver

import (
	"context"
	"errors"
	"miniblog/internal/pkg/log"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	genericoptions "github.com/onexstack/onexstack/pkg/options"

	handler "miniblog/internal/apiserver/handler/grpc"

	apiv1 "miniblog/pkg/api/apiserver/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	// GRPCServerMode 定义 gRPC 服务模式.
	// 使用 gRPC 框架启动一个 gRPC 服务器.
	GRPCServerMode = "grpc"
	// GRPCServerMode 定义 gRPC + HTTP 服务模式.
	// 使用 gRPC 框架启动一个 gRPC 服务器 + HTTP 反向代理服务器.
	GRPCGatewayServerMode = "grpc-gateway"
	// GinServerMode 定义 Gin 服务模式.
	// 使用 Gin Web 框架启动一个 HTTP 服务器.
	GinServerMode = "gin"
)

// 基于初始化配置创建运行时配置

// 存储应用相关配置
type Config struct {
	ServerMode  string
	JWTKey      string
	Expiration  time.Duration
	HTTPOptions *genericoptions.HTTPOptions
	GRPCOptions *genericoptions.GRPCOptions
}

// UnionServer 定义一个联合服务器. 根据 ServerMode 决定要启动的服务器类型.
//
// 联合服务器分为以下 2 大类:
//  1. Gin 服务器: 由 Gin 框架创建的标准的 REST 服务器. 根据是否开启 TLS, 来判断启动 HTTP 或者 HTTPS
//  2. GRPC 服务器：由 gRPC 框架创建的标准 RPC 服务器
//  3. HTTP 反向代理服务器: 由 grpc-gateway 框架创建的 HTTP 反向代理服务器
//     根据是否开启 TLS, 来判断启动 HTTP 或者 HTTPS
//
// HTTP 反向代理服务器依赖 gRPC 服务器, 所以在开启 HTTP 反向代理服务器时, 会先启动 gRPC 服务器.
type UnionServer struct {
	cfg *Config
	srv *grpc.Server
	lis net.Listener
}

// // 根据配置创建联合服务器
// func (cfg *Config) NewUnionServer() (*UnionServer, error) {
// 	return &UnionServer{cfg: cfg}, nil
// }

// // Run运行应用
// func (s *UnionServer) Run() error {
// 	// 打印配置内容
// 	// fmt.Printf("ServerMode from ServerOptions: %s\n", s.cfg.JWTKey)
// 	// fmt.Printf("ServerMode from Viper: %s\n\n", viper.GetString("jwt-key"))

// 	// log包打印
// 	log.Infow("ServerMode from ServerOptions", "jwt-key", s.cfg.JWTKey)
// 	log.Infow("ServerMode from Viper", "jwt-key", viper.GetString("jwt-key"))

// 	// jsonData, _ := json.MarshalIndent(s.cfg, "", " ")
// 	// fmt.Println(jsonData)

// 	// 空的select{} 语句会永久阻塞当前goroutine
// 	// 在服务器应用中, 这种模式通常用于保持main goroutine不退出,
// 	// 但前提是有其他goroutine在运行当没有其他活跃的goroutine时,
// 	// 就会触发"all goroutines are asleep - deadlock"错误
// 	// select {}
// 	return nil
// }

// NewUnionServer 根据配置创建联合服务器.
func (cfg *Config) NewUnionServer() (*UnionServer, error) {
	lis, err := net.Listen("tcp", cfg.GRPCOptions.Addr)
	if err != nil {
		log.Fatalw("Failed to listen", "err", err)
		return nil, err
	}
	// 创建一个gRPC服务实例grcsrv
	grpcsrv := grpc.NewServer()
	// 调用apiv1.RegisterMiniBlogServer方法将miniblog服务的处理器注册到gRPC服务器中
	// handler.NewHandler()返回一个服务器处理实例, 该实例实现了MiniBlog服务的业务逻辑
	apiv1.RegisterMiniBlogServer(grpcsrv, handler.NewHandler())

	return &UnionServer{cfg: cfg, srv: grpcsrv, lis: lis}, nil
}

func (s *UnionServer) Run() error {
	log.Infow("Start to listening the incoming requests on grpc address", "addr", s.cfg.GRPCOptions.Addr)

	// 在协程中启动grpc服务
	// 要先于http服务器启动, 否则http服务器无法转发请求到grpc服务器
	go s.srv.Serve(s.lis)

	// insecure.NewCredentials()用于创建不安全的传输凭据Transport Credentials
	// 因为http请求转发到grpc客户端时内部转发行为, 所以这里不用进行通信加密和身份验证
	dialOptions := []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials())}

	// 通过NewClient创建grpc客户端连接 conn
	conn, err := grpc.NewClient(s.cfg.GRPCOptions.Addr, dialOptions...)
	if err != nil {
		return err
	}

	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			// 设置序列号protobuf数据时, 枚举类型的字段以数字格式输出
			// 否则, 默认会以字符串格式输出, 跟枚举类型定义不一致
			UseEnumNumbers: true,
		},
	}))

	// 注册http路由, 并将grpc服务的方法注册为http rest接口, 并将http请求转换为grpc接口请求, 发送到grpc客户端连接 conn 中
	if err := apiv1.RegisterMiniBlogHandler(context.Background(), gwmux, conn); err != nil {
		return err
	}

	log.Infow("Start to Listen the incoming requests", "protocol", "http", "addr", s.cfg.HTTPOptions.Addr)
	// 创建http服务实例httpsrv
	httpsrv := &http.Server{
		Addr:    s.cfg.HTTPOptions.Addr,
		Handler: gwmux,
	}

	// 调用httpsrv.ListenAndServe启动http服务
	if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	return nil
}
