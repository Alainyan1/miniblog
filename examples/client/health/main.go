package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"time"

	"log"

	apiv1 "miniblog/pkg/api/apiserver/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	// grpc服务的地址
	addr = flag.String("addr", "localhost:6666", "The grpc server address to connect to")
	// 限制列出用户的数量
	limit = flag.Int64("limit", 10, "Limit to list users")
)

func main() {
	// 解析命令行参数
	flag.Parse()

	// 建立与 gRPC 服务器的连接
	// grpc.NewClient 用于建立客户端与 gRPC 服务端的连接
	// grpc.WithTransportCredentials(insecure.NewCredentials()) 表示使用不安全的传输（即不使用 TLS）
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Fatalf("Failed to connect to grpc server: %v", err)
	}
	defer conn.Close()

	// 创建miniblog客户端
	// 使用连接创建一个MiniBlog的grpc客户端实例
	client := apiv1.NewMiniBlogClient(conn)

	// 设置上下文, 带有3秒的超时时间
	// context.WithTimeout 用于设置调用的超时时间，防止请求无限等待
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// 调用miniblog到Healthz方法, 发起grpc请求, 检查服务健康状况
	resp, err := client.Healthz(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to call healthzz: %v", err)
	}

	// 使用json.Marshal将返回的响应数据转换为JSON格式
	// 输出json数据到终端
	jsonData, _ := json.Marshal(resp)
	fmt.Println(string(jsonData))
}
