# ============================
# 全局Makefile变量

# MAKEFILE_LIST是Makefile内置变量, 包含了所有被读取的Makefile文件名列表
# 当前Makefile文件的路径总是为列表的最后一个元素
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
# 项目根目录, 将其转为绝对路径
# PROJECT_ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR))/ && pwd -P)
# 直接指定项目根目录
PROJECT_ROOT_DIR := /Users/paulfaltings/Desktop/goproject/miniblog
# 构建产物, 临时文件存放目录
OUTPUT_DIR := $(PROJECT_ROOT_DIR)/_output

# Protobuf文件路径
APIROOT = $(PROJECT_ROOT_DIR)/pkg/api

# makefile的shell切换为bash
SHELL := /bin/bash

# 设置单元测试覆盖率阈值. 这里的 1 只是示例用
# 生产环境建议阈值设置高一点, 例如60.
COVERAGE := 1


# ============================
# 定义版本相关变量
# 指定应用使用的version包, 会通过 `-ldflags -X` 向该包中指定的变量注入值
# 与go.mod文件中的module路径一致
VERSION_PACKAGE=github.com/onexstack/onexstack/pkg/version
# 定义VERSION语义化版本号
ifeq ($(origin VERSION), undefined)
  # 如果VERSION未定义, 则使用git命令获取版本号
  # --tags: 使用所有标签, 而不是仅使用带注释的标签
  # --always: 如果没有找到标签, 则使用提交id的所有作为替代
  # --match= <>: 仅考虑与指定模式匹配的标签, 例如'v*'仅匹配以v开头的标签, 确保符合语义化版本规范
  VERSION := $(shell git describe --tags --always --match='v*')
endif

# 检查代码仓库是否为dirty(默认为dirty)
GIT_TREE_STATE := "dirty"
ifeq (, $(shell git status --porcelain 2>/dev/null))
  # 如果没有未提交的更改, 则设置GIT_TREE_STATE为clean
  GIT_TREE_STATE := "clean"
endif

# 使用git rev-parse命令获取构建时的提交ID
GIT_COMMIT := $(shell git rev-parse HEAD)

# 使用 date -u 获取构建时间, -u选项表示使用UTC时间
GO_LDFLAGS += \
    -X $(VERSION_PACKAGE).gitVersion=$(VERSION) \
    -X $(VERSION_PACKAGE).gitCommit=$(GIT_COMMIT) \
    -X $(VERSION_PACKAGE).gitTreeState=$(GIT_TREE_STATE) \
    -X $(VERSION_PACKAGE).buildDate=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
# ============================
# 默认目标为all
.DEFAULT_GOAL := all

# Makefile all 伪目标, 执行`make`时默认执行all伪目标
# 该目标依赖于tidy format build add-copyright, 依次执行
.PHONY: all
all: tidy format build add-copyright

# =============================
# 其他伪目标
# 编译源码, 依赖tidy目标自动添加/移除依赖包
# 编译mb-apiserver, 生成的二进制文件存放在$(OUTPUT_DIR)目录下
.PHONY: build
build: tidy 
	@go build -v -ldflags "$(GO_LDFLAGS)" -o $(OUTPUT_DIR)/mb-apiserver $(PROJECT_ROOT_DIR)/cmd/mb-apiserver/main.go

# 格式化go源码
.PHONY: format
format:
	@gofmt -s -w ./

# 添加版权信息
.PHONY: add-copyright
add-copyright:
	@addlicense -v -f $(PROJECT_ROOT_DIR)/scripts/boilerplate.txt $(PROJECT_ROOT_DIR) --skip-dirs=third_party,vendor,$(OUTPUT_DIR)

# 自动添加/移除go依赖包
.PHONY: tidy
tidy:
	@go mod tidy

# 清理构建产物, 幂等删除, 在临时目录不存在时Makefile执行仍然成功
.PHONY: clean
clean:
	@-rm -vrf $(OUTPUT_DIR)

# protoc是protocol buffers文件的编译工具
# 通过插件机制实现对不同语言的支持, 例如使用--xxx_out时会首先查询是否存在内置的xxx插件
# 如果没有内置插件, 则会继续查询系统中是否存在名为xxx的可执行程序
# --go_out使用的插件为protoc-gen-go
# --proto_path或-l: 用于指定编译源码的搜索路径m 在构建.proto时, 会在这些路径下查找所需的Protobuf文件及其依赖
# --go_out: 用于生成与gRPC相关的go代码, 并配置生成文件的路径和文件结构, 主要参数包括plugins, paths, 分别表示生成go代码需要的插件及生成代码的位置
# paths参数支持import和source_relative两个选项
# import: 按照生成的go代码包的全路径创建目录结构
# source_relative: 生成文件应保持与输入文件相对路径一致
# --go-grpc_out: 功能与--go_out类似, 指定生成的*_grpc.pb.go文件的存放路径
.PHONY: protoc
protoc: #编译protobuf文件
	@echo "=========> Generate Protobuf Files"
	@mkdir -p $(PROJECT_ROOT_DIR)/api/openapi
	@# --grpc-gateway_out 用来在 pkg/api/apiserver/v1/ 目录下生成反向服务器代码 apiserver.pb.gw.go
	@# --openapiv2_out 用来在 api/openapi/apiserver/v1/ 目录下生成 Swagger V2 接口文档
	@protoc   \
		--proto_path=$(APIROOT) \
		--proto_path=$(PROJECT_ROOT_DIR)/third_party/protobuf \
		--go_out=paths=source_relative:$(APIROOT) \
		--go-grpc_out=paths=source_relative:$(APIROOT) \
		--grpc-gateway_out=allow_delete_body=true,paths=source_relative:$(APIROOT)	\
		--openapiv2_out=$(PROJECT_ROOT_DIR)/api/openapi		\
		--openapiv2_opt=allow_delete_body=true,logtostderr=true		\
		--defaults_out=paths=source_relative:$(APIROOT) \
		$(shell find $(APIROOT) -name *.proto)
	@find . -name '*.pb.go' -exec protoc-go-inject-tag -input={} \; 

.PHONY: debug-protoc
debug-protoc:
	@echo "PROJECT_ROOT_DIR: $(PROJECT_ROOT_DIR)" 
	@echo "APIROOT: $(APIROOT)"
	@echo "Proto files: $(shell find $(APIROOT) -name *.proto)"

.PHONY: ca
ca: #生成CA文件
	@mkdir -p $(OUTPUT_DIR)/cert
	@# 1. 生成根证书私钥(CA key), 指定私钥大小4096
	@openssl genrsa -out $(OUTPUT_DIR)/cert/ca.key 4096
	@# 2. 使用根私钥生成证书签名请求文件(CA CSR), 有效期为 10 年
	@# req命令创建证书的请求文件, nodes表示不加密私钥, new创建新的证书请求文件, key指定使用的私钥文件, out指定输出文件路径, subj设置或修改证书请求中的主体信息
	@openssl req -new -nodes -key $(OUTPUT_DIR)/cert/ca.key -sha256 -days 3650 -out $(OUTPUT_DIR)/cert/ca.csr \
		-subj "/C=CN/ST=Guangdong/L=Shenzhen/O=miniblog/OU=it/CN=127.0.0.1/emailAddress=alainyan@yahoo.com"
	@# 3. 使用根私钥签发根证书(CA CRT), 使其自签名, x509用于创建和修改509证书
	@# in指定输入文件, out指定输出文件, req指定输入文件为证书请求, signkey指定用于自签名的私钥文件
	@openssl x509 -req -days 365 -in $(OUTPUT_DIR)/cert/ca.csr -signkey $(OUTPUT_DIR)/cert/ca.key -out $(OUTPUT_DIR)/cert/ca.crt
	@# 4. 生成服务端私钥
	@openssl genrsa -out $(OUTPUT_DIR)/cert/server.key 2048
	@# 5. 使用服务端私钥生成服务端的证书签名请求(Server CSR)
	@# .crt, .cert代表Certificate, crt常见于unix, cer常见于win, .key存放私钥或公钥, .csr代表证书签名请求, 不是证书
	@openssl req -new -key $(OUTPUT_DIR)/cert/server.key -out $(OUTPUT_DIR)/cert/server.csr \
		-subj "/C=CN/ST=Guangdong/L=Shenzhen/O=serverdevops/OU=serverit/CN=localhost/emailAddress=alainyan@yahoo.com" \
		-addext "subjectAltName=DNS:localhost,IP:127.0.0.1"
	@#6. 使用根证书(CA)签发服务端证书(Server CRT)
	@# CA设置CA文件, 必须为PEM格式(文本格式), CAkey设置CA私钥文件, 必须为PEM格式, CAcreateserial创建序列号文件
	@openssl x509 -days 365 -sha256 -req -CA $(OUTPUT_DIR)/cert/ca.crt -CAkey $(OUTPUT_DIR)/cert/ca.key \
		-CAcreateserial -in $(OUTPUT_DIR)/cert/server.csr -out $(OUTPUT_DIR)/cert/server.crt -extensions v3_req \
		-extfile <(printf "[v3_req]\nsubjectAltName=DNS:localhost,IP:127.0.0.1")

.PHONY: test
test: # 执行单元测试
	@echo "==========> Running unit tests"
	@mkdir -p $(OUTPUT_DIR)
	@go test -race -cover \
		-coverprofile=$(OUTPUT_DIR)/coverage.out \
		-timeout=10m -shuffle=on -short \
		-v `go list ./...|egrep -v 'tools|vendor|third_party'`

.PHONY: cover
cover: test ## 执行单元测试并校验覆盖率
	@echo "==============> Running code coverage tests"
	@go tool cover -func=$(OUTPUT_DIR)/coverage.out | awk -v target=$(COVERAGE) -f $(PROJECT_ROOT_DIR)/scripts/coverage.awk

.PHONY: lint
lint: #执行静态代码检查
	@echo "=============> Running golangci to lint source codes"
	@golangci-lint run -c $(PROJECT_ROOT_DIR)/.golangci.yaml $(PROJECT_ROOT_DIR)/...