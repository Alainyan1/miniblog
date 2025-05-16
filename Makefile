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

# 默认目标为all
.DEFAULT_GOAL := all

# Makefile all 伪目标, 执行`make`时默认执行all伪目标
# 该目标依赖于tidy format build add-copyright, 依次执行
.PHONY: all
all: tidy format build add-copyright

# 编译源码, 依赖tidy目标自动添加/移除依赖包
# 编译mb-apiserver, 生成的二进制文件存放在$(OUTPUT_DIR)目录下
.PHONY: build
build: tidy 
	@go build -v -o $(OUTPUT_DIR)/mb-apiserver $(PROJECT_ROOT_DIR)/cmd/mb-apiserver/main.go

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