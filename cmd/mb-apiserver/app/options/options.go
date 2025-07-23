// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package options

import (
	"errors"
	"fmt"
	"time"

	genericoptions "github.com/onexstack/onexstack/pkg/options"
	stringsutil "github.com/onexstack/onexstack/pkg/util/strings"
	"github.com/spf13/pflag"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"

	"miniblog/internal/apiserver" // 控制面依赖数据面
)

// 支持服务器模式集合.
var availableServerModes = sets.New(
	// "grpc",
	// "grpc-gateway",
	// "gin",
	apiserver.GinServerMode,
	apiserver.GRPCServerMode,
	apiserver.GRPCGatewayServerMode,
)

// 服务器配置选项.
type ServerOptions struct {
	// ServerMode定义服务器模式: gRPC, Gin HTTP, HTTP Reverse Proxy
	ServerMode string `json:"server-mode" mapstructure:"server-mode"`

	// JWTKey定义JWT密钥
	JWTKey string `json:"jwt-key" mapstructure:"jwt-key"`

	// Expiration定义JWT token过期时间
	Expiration time.Duration `json:"expiration" mapstructure:"expiration"`

	// HTTPOptions包含http配置选项
	HTTPOptions *genericoptions.HTTPOptions `json:"http" mapstructure:"http"`

	// GRPCOptions包含gRPC配置选项
	GRPCOptions *genericoptions.GRPCOptions `json:"grpc" mapstructure:"grpc"`

	// MySQLOptions 包含 MySQL 配置选项
	MySQLOptions *genericoptions.MySQLOptions `json:"mysql" mapstructure:"mysql"`

	// TLSOptions 包含 TLS 配置选项.
	TLSOptions *genericoptions.TLSOptions `json:"tls" mapstructure:"tls"`
}

// 创建带有默认值的ServerOptions实例.
func NewServerOptions() *ServerOptions {
	opts := &ServerOptions{
		ServerMode:   apiserver.GRPCGatewayServerMode,
		JWTKey:       "Rtg8BPKNEf2mB4mgvKONGPZZQSaJWNLijxR42qRgq0iBb5",
		Expiration:   2 * time.Hour,
		TLSOptions:   genericoptions.NewTLSOptions(),
		HTTPOptions:  genericoptions.NewHTTPOptions(),
		GRPCOptions:  genericoptions.NewGRPCOptions(),
		MySQLOptions: genericoptions.NewMySQLOptions(),
	}
	opts.HTTPOptions.Addr = ":5555"
	opts.GRPCOptions.Addr = ":6666"

	return opts
}

// 从而支持通过命令行选项来给ServerOption结构体中的字段赋值.
func (o *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&o.ServerMode, "server-mode", o.ServerMode, fmt.Sprintf("Server mode, available options: %v", availableServerModes.UnsortedList()))
	fs.StringVar(&o.JWTKey, "jwt-key", o.JWTKey, "JWT signing key. Must be at least 6 characters long.")
	// 绑定 JWT Token 的过期时间选项到命令行标志。
	// 参数名称为 `--expiration`, 默认值为 o.Expiration
	fs.DurationVar(&o.Expiration, "expiration", o.Expiration, "The expiration duration of JWT tokens.")
	o.TLSOptions.AddFlags(fs)
	o.HTTPOptions.AddFlags(fs)
	o.GRPCOptions.AddFlags(fs)
	o.MySQLOptions.AddFlags(fs)
}

// 检验ServerOptions中的选项是否合法.
func (o *ServerOptions) Validate() error {
	errs := []error{}

	// 校验ServerMode是否有效
	if !availableServerModes.Has(o.ServerMode) {
		errs = append(errs, fmt.Errorf("invalid server mode: must be one of %v", availableServerModes.UnsortedList()))
	}

	// 校验JWTKey长度
	if len(o.JWTKey) < 6 {
		errs = append(errs, errors.New("JWTKey must be at least 6 characters long"))
	}

	// 校验子选项
	errs = append(errs, o.TLSOptions.Validate()...)
	errs = append(errs, o.HTTPOptions.Validate()...)
	errs = append(errs, o.MySQLOptions.Validate()...)

	// 如果是grpc或grpc-gateway模式, 校验grpc配置
	if stringsutil.StringIn(o.ServerMode, []string{apiserver.GRPCServerMode, apiserver.GRPCGatewayServerMode}) {
		errs = append(errs, o.GRPCOptions.Validate()...)
	}

	// 合并错误并返回
	return utilerrors.NewAggregate(errs)
}

// Config方法基于ServerOptions创建新的apiserver.Config.
func (o *ServerOptions) Config() (*apiserver.Config, error) {
	return &apiserver.Config{
		ServerMode:   o.ServerMode,
		JWTKey:       o.JWTKey,
		Expiration:   o.Expiration,
		TLSOptions:   o.TLSOptions,
		HTTPOptions:  o.HTTPOptions,
		GRPCOptions:  o.GRPCOptions,
		MySQLOptions: o.MySQLOptions,
	}, nil
}
