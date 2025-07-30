# ============================
# 全局Makefile变量

# makefile的shell切换为bash
SHELL := /bin/bash

# MAKEFILE_LIST是Makefile内置变量, 包含了所有被读取的Makefile文件名列表
# 当前Makefile文件的路径总是为列表的最后一个元素
COMMON_SELF_DIR := $(dir $(lastword $(MAKEFILE_LIST)))
# 项目根目录, 将其转为绝对路径
# PROJECT_ROOT_DIR := $(abspath $(shell cd $(COMMON_SELF_DIR))/ && pwd -P)
# 直接指定项目根目录
PROJECT_ROOT_DIR := /Users/paulfaltings/Desktop/goproject/miniblog
# 构建产物, 临时文件存放目录
OUTPUT_DIR := $(PROJECT_ROOT_DIR)/_output

ROOT_PACKAGE=github.com/alainyan/miniblog

# Protobuf文件路径
APIROOT = $(PROJECT_ROOT_DIR)/pkg/api


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

# 编译的操作系统可以是 linux/windows/darwin
PLATFORMS ?= darwin_amd64 windows_amd64 linux_amd64 linux_arm64

# 设置一个指定的操作系统
ifeq ($(origin PLATFORM), undefined)
	ifeq ($(origin GOOS), undefined)
# 设置GOOS(操作系统)
		GOOS := $(shell go env GOOS)
	endif
	ifeq ($(origin GOARCH), undefined)
# 设置GOARCH架构
		GOARCH := $(shell go env GOARCH)
	endif
# 构建PLATFORM 变量, 例如 linux_amd64、darwin_arm64
	PLATFORM := $(GOOS)_$(GOARCH)
# 构建镜像时，使用 darwin 作为默认的 OS
	IMAGE_PLAT := darwin_$(GOARCH)
else
	GOOS := $(word 1, $(subst _, ,$(PLATFORM)))
	GOARCH := $(word 2, $(subst _, ,$(PLATFORM)))
	IMAGE_PLAT := $(PLATFORM)
endif


# 设置单元测试覆盖率阈值. 这里的 1 只是示例用
# 生产环境建议阈值设置高一点, 例如60.
ifeq ($(origin COVERAGE),undefined)
COVERAGE := 1
endif

# Makefile 设置
# 默认情况下, make 在处理子目录时会打印进入和离开目录的信息, 这在复杂项目中可能导致输出冗长
# 如果用户没有定义 V, 例如通过 make V=1, 则启用 --no-print-directory, 减少输出信息, 使构建日志更简洁
ifndef V
MAKEFLAGS += --no-print-directory
endif

# Linux 命令设置
# FIND 用于在 Makefile 中查找项目中的文件, 但排除第三方库或依赖目录, 这些目录通常包含外部代码, 不需要参与某些构建或检查任务
FIND := find . ! -path './third_party/*' ! -path './vendor/*'
# XARGS 配合 FIND 使用, 确保在处理文件列表时, 如果 find 没有返回任何文件, xargs 不会错误执行后续命令
XARGS := xargs --no-run-if-empty