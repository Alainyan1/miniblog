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

# ============================
# 定义版本相关变量
# 指定应用使用的version包, 会通过 `-ldflags -X` 向该包中指定的变量注入值 
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