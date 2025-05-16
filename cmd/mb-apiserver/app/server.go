// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package app

import (
	"miniblog/cmd/mb-apiserver/app/options"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configFile string

// 创建一个 *cobra.Command 对象, 用于启动应用程序
func NewMiniBlogCommand() *cobra.Command {
	// 默认命令行选项
	opts := options.NewServerOptions()

	cmd := &cobra.Command{
		// 指定命令名字, 该名字会出现在帮助信息中
		Use: "mb-apiserver",
		// 命令简短描述
		Short: "A mini blog show best practices for develop a full-featured Go project",

		// 命令详细描述
		Long: `A mini blog show best practices for develop a full-featured Go project.

The project features include:
• Utilization of a clean architecture;
• Use of many commonly used Go packages: gorm, casbin, govalidator, jwt, gin, 
  cobra, viper, pflag, zap, pprof, grpc, protobuf, grpc-gateway, etc.;
• A standardized directory structure following the project-layout convention;
• Authentication (JWT) and authorization features (casbin);
• Independently designed log and error packages;
• Management of the project using a high-quality Makefile;
• Static code analysis;
• Includes unit tests, performance tests, fuzz tests, and mock tests;
• Rich web functionalities (tracing, graceful shutdown, middleware, CORS, 
  recovery from panics, etc.);
• Implementation of HTTP, HTTPS, and gRPC servers;
• Implementation of JSON and Protobuf data exchange formats;
• The project adheres to numerous development standards: 
  code standards, versioning standards, API standards, logging standards, 
  error handling standards, commit standards, etc.;
• Access to MySQL with programming implementation;
• Implemented business functionalities: user management and blog management;
• RESTful API design standards;
• OpenAPI 3.0/Swagger 2.0 API documentation;
• High-quality code.`,
		// 命令出错时, 不打印帮助信息。设置为 true 可以确保命令出错时一眼就能看到错误信息
		SilenceUsage: true,
		// 指定调用cmd.Execute()时, 执行的Run函数
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(opts)
		},
		// 设置命令运行时的参数检查, 不需要指定命令行参数。例如: ./miniblog param1 param2
		Args: cobra.NoArgs,
	}

	// 通过cobra.OnInitialize注册一个回调函数, 该函数在每次运行任意命令时调用
	// 初始化配置函数, 在每个命令运行时调用
	// 确保在程序运行时，将--config命令行选项指定的配置文件内容加载到viper中
	cobra.OnInitialize(OnInitialize)

	// cobra 支持持久性标志(PersistentFlag)，该标志可用于它所分配的命令以及该命令下的每个子命令
	// 推荐使用配置文件来配置应用，便于管理配置项
	// 通过PersistentFlags().StringVarP调用, 给应用添加了--config/-c命令行选项, 用来指定配置文件路径
	// 将配置文件保存到configFile变量中, configFile默认值由filePath函数生成
	cmd.PersistentFlags().StringVarP(&configFile, "config", "c", filePath(), "Path to the miniblog configuration file.")

	// 将 ServerOptions 中的选项绑定到cmd.PersistentFlags标志集中
	opts.AddFlags(cmd.PersistentFlags())
	return cmd
}

// 主运行逻辑, 复制初始化日志, 解析配置, 校验选项并启动服务器
func run(opts *options.ServerOptions) error {
	// 将 viper 中的配置解析到 opts.
	if err := viper.Unmarshal(opts); err != nil {
		return err
	}

	// 校验命令行选项
	if err := opts.Validate(); err != nil {
		return err
	}

	// 获取应用配置
	// 将命令行选项和应用应用配置分开, 更加灵活处理2种不同类型的配置
	cfg, err := opts.Config()
	if err != nil {
		return err
	}

	// 创建联合服务实例
	server, err := cfg.NewUnionServer()
	if err != nil {
		return err
	}

	// 启动服务器
	return server.Run()
}
