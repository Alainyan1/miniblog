// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package apiserver

import (
	"context"
	"miniblog/internal/apiserver/biz"
	"miniblog/internal/apiserver/model"
	"miniblog/internal/apiserver/pkg/validation"
	"miniblog/internal/apiserver/store"
	"miniblog/internal/pkg/contextx"
	"miniblog/internal/pkg/known"
	"miniblog/internal/pkg/log"
	"miniblog/pkg/auth"
	"miniblog/pkg/token"
	"os"
	"os/signal"

	"syscall"
	"time"

	mw "miniblog/internal/pkg/middleware/gin"
	"miniblog/internal/pkg/server"

	genericoptions "github.com/onexstack/onexstack/pkg/options"
	"github.com/onexstack/onexstack/pkg/store/where"
	"gorm.io/gorm"
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
	ServerMode   string
	JWTKey       string
	Expiration   time.Duration
	HTTPOptions  *genericoptions.HTTPOptions
	GRPCOptions  *genericoptions.GRPCOptions
	MySQLOptions *genericoptions.MySQLOptions
	TLSOptions   *genericoptions.TLSOptions
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
	srv server.Server
}

type ServerConfig struct {
	cfg       *Config
	biz       biz.IBiz
	val       *validation.Validator
	retriever mw.UserRetriever
	authz     *auth.Authz
}

// NewUnionServer 根据配置创建联合服务器.
func (cfg *Config) NewUnionServer() (*UnionServer, error) {

	// 注册租户解析函数, 通过上下文获取用户ID, nolint: gocritic告诉静态分析工具(如 golangci-lint)忽略特定代码行的某些检查规则
	//nolint: gocritic
	where.RegisterTenant("userID", func(ctx context.Context) string {
		return contextx.UserID(ctx)
	})

	// 初始化 token 包的签名密钥、认证 Key 及 Token 默认过期时间
	token.Init(cfg.JWTKey, known.XUserID, cfg.Expiration)

	// 创建服务配置
	serverConfig, err := cfg.NewServerConfig()
	if err != nil {
		return nil, err
	}

	log.Infow("Initializing federation server", "server-mode", cfg.ServerMode)

	// 根据服务模式创建对应的服务实例
	// 实际企业开发中, 可以根据需要只选择一种服务器模式.
	// 这里为了方便给你展示, 通过 cfg.ServerMode 同时支持了 Gin 和 GRPC 2 种服务器模式.
	// 默认为 gRPC 服务器模式.
	var srv server.Server
	switch cfg.ServerMode {
	case GinServerMode:
		srv, err = serverConfig.NewGinServer(), nil
	default:
		srv, err = serverConfig.NewGRPCServerOr()
	}

	if err != nil {
		return nil, err
	}

	return &UnionServer{srv: srv}, nil
}

func (s *UnionServer) Run() error {
	// 使用协程启动服务器, 保证主线程不会被阻塞
	go s.srv.RunOrDie()

	// 创建一个os.Signal类型的通道用于接收系统信号
	quit := make(chan os.Signal, 1)

	// 当执行 kill 命令时（不带参数），默认会发送 syscall.SIGTERM 信号
	// 使用 kill -2 命令会发送 syscall.SIGINT 信号（例如按 CTRL+C 触发）
	// 使用 kill -9 命令会发送 syscall.SIGKILL 信号, 但 SIGKILL 信号无法被捕获, 因此无需监听和处理
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	// 阻塞程序, 等待从quit channel中接受信号
	<-quit
	log.Infow("Shutting down server...")

	// 优雅关停, 通过context.WithTimeout创建上下文对象, 为优雅关停服务提供超时控制
	// 确保在一定时间内完成清理工作, 如果超过指定时间, 服务将被强制终止
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 先关闭依赖的服务, 在关闭被依赖的服务, ctx被传递给s.srv.GracefulStop方法, 用于通知服务相关的协程
	// 服务中的任务可以通过监听ctx.Done来检测是否需要终止
	s.srv.GracefulStop(ctx)
	log.Infow("Server exited")
	return nil
}

// 创建一个gorm.DB实例
func (cfg *Config) NewDB() (*gorm.DB, error) {
	return cfg.MySQLOptions.NewDB()
}

// 创建一个*ServerConfig实例
// 后续可以使用依赖注入的方式
func (cfg *Config) NewServerConfig() (*ServerConfig, error) {
	db, err := cfg.NewDB()
	if err != nil {
		return nil, err
	}

	store := store.NewStore(db)

	// 初始化权限认证模块
	authz, err := auth.NewAuthz(store.DB(context.TODO()))

	return &ServerConfig{
		cfg:       cfg,
		biz:       biz.NewBiz(store, authz),
		val:       validation.New(store),
		retriever: &UserRetriever{store: store},
		authz:     authz,
	}, nil
}

// UserRetriever 定义一个用户数据获取器. 用来获取用户信息.
type UserRetriever struct {
	store store.IStore
}

// GetUser 根据用户 ID 获取用户信息.
func (r *UserRetriever) GetUser(ctx context.Context, userID string) (*model.UserM, error) {
	return r.store.User().Get(ctx, where.F("userID", userID))
}

// func (s *UnionServer) Run() error {
// 	log.Infow("Start to listening the incoming requests on grpc address", "addr", s.cfg.GRPCOptions.Addr)

// 	// 在协程中启动grpc服务
// 	// 要先于http服务器启动, 否则http服务器无法转发请求到grpc服务器
// 	go s.srv.Serve(s.lis)

// 	// insecure.NewCredentials()用于创建不安全的传输凭据Transport Credentials
// 	// 因为http请求转发到grpc客户端时内部转发行为, 所以这里不用进行通信加密和身份验证
// 	dialOptions := []grpc.DialOption{grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials())}

// 	// 通过NewClient创建grpc客户端连接 conn
// 	conn, err := grpc.NewClient(s.cfg.GRPCOptions.Addr, dialOptions...)
// 	if err != nil {
// 		return err
// 	}

// 	gwmux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
// 		MarshalOptions: protojson.MarshalOptions{
// 			// 设置序列号protobuf数据时, 枚举类型的字段以数字格式输出
// 			// 否则, 默认会以字符串格式输出, 跟枚举类型定义不一致
// 			UseEnumNumbers: true,
// 		},
// 	}))

// 	// 注册http路由, 并将grpc服务的方法注册为http rest接口, 并将http请求转换为grpc接口请求, 发送到grpc客户端连接 conn 中
// 	if err := apiv1.RegisterMiniBlogHandler(context.Background(), gwmux, conn); err != nil {
// 		return err
// 	}

// 	log.Infow("Start to Listen the incoming requests", "protocol", "http", "addr", s.cfg.HTTPOptions.Addr)
// 	// 创建http服务实例httpsrv
// 	httpsrv := &http.Server{
// 		Addr:    s.cfg.HTTPOptions.Addr,
// 		Handler: gwmux,
// 	}

// 	// 调用httpsrv.ListenAndServe启动http服务
// 	if err := httpsrv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
// 		return err
// 	}
// 	return nil
// }
