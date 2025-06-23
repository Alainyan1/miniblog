// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"

	apiv1 "miniblog/pkg/api/apiserver/v1"

	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

var (
	// 定义命令行参数
	// grpc服务地址
	addr = flag.String("addr", "localhost:6666", "The grpc server address to connect to.")
	// 限制列出用户的数量
	limit = flag.Int64("limit", 10, "Limit to list users.")
)

func main() {
	// 解析命令行参数
	flag.Parse()

	// 建立与grpc服务器的连接
	conn, err := grpc.NewClient(
		*addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(unaryClientInterceptor()))

	if err != nil {
		log.Fatalf("Failed to connect to grpc server: %v", err)
	}
	// 确保函数结束时关闭连接, 避免资源泄露
	defer conn.Close()

	// 创建miniblog客户端
	client := apiv1.NewMiniBlogClient(conn)

	// 设置上下文
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	// 取消上下文, 释放资源
	defer cancel()

	// 创建metadata传递请求元数据
	md := metadata.Pairs("customer-header", "value123")
	ctx = metadata.NewOutgoingContext(ctx, md)

	// 用于存储返回的header元数据
	var header metadata.MD
	// 调用healthz方法
	resp, err := client.Healthz(ctx, nil, grpc.Header(&header))
	if err != nil {
		log.Fatalf("Failed to call healthz: %v", err)
	}

	for key, val := range header {
		fmt.Printf("Response header (ket: %s, value: %s)\n", key, val)
	}

	// 返回的数据转换为json格式
	jsonData, _ := json.Marshal(resp)
	fmt.Println(string(jsonData))
}

func unaryClientInterceptor() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context, // 上下文对象, 包含请求的元数据和链路信息
		method string, // 调用的grpc方法名称
		req, reply interface{}, // 请求参数和响应结果(具体的结构体, 通常由.proto定义的服务方法生成)
		cc *grpc.ClientConn, // 客户端连接对象, 表示连接的grpc服务的
		invoker grpc.UnaryInvoker, // grpc的实际调用方法, 客户端拦截器需要调用他来完成请求的转发, 例如向grpc服务器发送请求并获取响应
		opts ...grpc.CallOption, // 可选调用项, 如超时时间, 拦截机制等
	) error {
		// 打印请求方法和参数
		log.Printf("[UnaryClientInterceptor] Invoking method: %s", method)

		// 添加自定义元数据
		md := metadata.Pairs("interceptor-header", "interceptor-value")
		// 创建的outgoing metadata会覆盖main函数创建的
		ctx = metadata.NewOutgoingContext(ctx, md)

		// 调用实际rpc方法
		err := invoker(ctx, method, req, reply, cc, opts...)

		if err != nil {
			log.Printf("[UnaryClientInterceptor] Method :%s, Error: %v", method, err)
			return err
		}
		return nil
	}
}
