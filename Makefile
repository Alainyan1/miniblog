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

# ============================
# 定义版本相关变量
# 指定应用使用的version包, 会通过 `-ldflags -X` 向该包中指定的变量注入值
# 与go.mod文件中的module路径一致
VERSION_PACKAGE=miniblog/pkg/version
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