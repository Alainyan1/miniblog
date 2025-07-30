# ==============================================================================
# 工具相关的 Makefile, 用于辅助开发工作
# 主要用于安装和验证 Go 项目中常用的开发工具

# ?= 表示如果 TOOLS 未在环境变量或命令行中定义, 则使用默认值
TOOLS ?= golangci-lint goimports protoc-plugins swagger addlicense protoc-go-inject-tag protolint

# 检查所有工具是否已经安装
tools.verify: $(addprefix tools.verify., $(TOOLS))

# 强制安装所有工具
tools.install: $(addprefix tools.install., $(TOOLS))

# 通用安装规则
tools.install.%:
	@echo "===========> Installing $*"
	@$(MAKE) install.$*

# 通用验证规则
tools.verify.%:
	@if ! which $* &>/dev/null; then $(MAKE) tools.install.$*; fi

# 安装具体的工具
install.golangci-lint:
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.63.2
	@golangci-lint completion bash > $(HOME)/.golangci-lint.bash
	@if ! grep -q .golangci-lint.bash $(HOME)/.bashrc; then echo "source \$$HOME/.golangci-lint.bash" >> $(HOME)/.bashrc; fi

install.goimports:
	@$(GO) install golang.org/x/tools/cmd/goimports@latest

install.protoc-plugins:
	@$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@v1.35.2
	@$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.5.1
	@$(GO) install github.com/onexstack/protoc-gen-defaults@v0.0.2
	@$(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.24.0
	@$(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@v2.24.0

install.swagger:
	@$(GO) install github.com/go-swagger/go-swagger/cmd/swagger@latest

install.addlicense:
	@$(GO) install github.com/marmotedu/addlicense@latest

install.protoc-go-inject-tag:
	@$(GO) install github.com/favadi/protoc-go-inject-tag@latest

install.protolint:
	@$(GO) install github.com/yoheimuta/protolint/cmd/protolint@latest

# 伪目标（防止文件与目标名称冲突）
.PHONY: tools.verify tools.install tools.install.% tools.verify.% install.golangci-lint \
	install.goimports install.protoc-plugins install.swagger \
	install.addlicense install.protoc-go-inject-tag protolint