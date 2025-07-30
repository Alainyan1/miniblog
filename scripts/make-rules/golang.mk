
# ==============================================================================
# 用来进行编译的 Makefile
#

GO := go

# 定义 GO_BUILD_FLAGS 变量, 追加链接器标志(linker flags)
GO_BUILD_FLAGS += -ldflags "$(GO_LDFLAGS)"

ifeq ($(GOOS),windows)
	GO_OUT_EXT := .exe
endif

ifeq ($(ROOT_PACKAGE),)
	$(error the variable ROOT_PACKAGE must be set prior to including golang.mk)
endif

# 如果 GOBIN 未定义, 设置为 $(GOPATH)/bin, 用于安装 Go 工具或二进制文件
GOPATH := $(shell go env GOPATH)
ifeq ($(origin GOBIN), undefined)
	GOBIN := $(GOPATH)/bin
endif

# COMMANDS: 查找 $(PROJ_ROOT_DIR)/cmd/ 目录下的所有子目录, 排除.md
COMMANDS ?= $(filter-out %.md, $(wildcard $(PROJ_ROOT_DIR)/cmd/*))
# BINS: 提取 COMMANDS 中每个子目录的名称(notdir 移除路径前缀)
BINS ?= $(foreach cmd,${COMMANDS},$(notdir $(cmd)))

ifeq ($(COMMANDS),)
  $(error Could not determine COMMANDS, set PROJ_ROOT_DIR or run in source dir)
endif
ifeq ($(BINS),)
  $(error Could not determine BINS, set PROJ_ROOT_DIR or run in source dir)
endif

go.build.verify:
	@if ! which go &>/dev/null; then echo "Cannot found go compile tool. Please install go tool first."; exit 1; fi

# go.build.% 是一个模式规则, 用于编译特定平台和命令的二进制文件, 匹配格式为linux_amd64.server, 平台为linux, 命令为server
# CGO_ENABLED=1: 启用 CGO(可能用于依赖 C 库的代码)
go.build.%: ## 编译 Go 源码.
	$(eval COMMAND := $(word 2,$(subst ., ,$*)))
	$(eval PLATFORM := $(word 1,$(subst ., ,$*)))
	$(eval OS := $(word 1,$(subst _, ,$(PLATFORM))))
	$(eval ARCH := $(word 2,$(subst _, ,$(PLATFORM))))
	@echo "===========> Building binary $(COMMAND) $(VERSION) for $(OS) $(ARCH)"
	@mkdir -p $(OUTPUT_DIR)/platforms/$(OS)/$(ARCH)
	@CGO_ENABLED=1 GOOS=$(OS) GOARCH=$(ARCH) $(GO) build $(GO_BUILD_FLAGS) \
		-o $(OUTPUT_DIR)/platforms/$(OS)/$(ARCH)/$(COMMAND)$(GO_OUT_EXT) \
		$(ROOT_PACKAGE)/cmd/$(COMMAND)

go.build: go.build.verify $(addprefix go.build., $(addprefix $(PLATFORM)., $(BINS))) # 根据指定的平台编译源码.

go.format: tools.verify.goimports ## 格式化 Go 源码.
	@echo "===========> Running formaters to format codes"
	@$(FIND) -type f -name '*.go' | $(XARGS) gofmt -s -w
	@$(FIND) -type f -name '*.go' | $(XARGS) goimports -w -local $(ROOT_PACKAGE)
	@$(GO) mod edit -fmt

go.tidy: ## 自动添加/移除依赖包.
	@echo "===========> Running 'go mod tidy'..."
	@$(GO) mod tidy

go.test: ## 执行单元测试.
	@echo "===========> Running unit tests"
	@mkdir -p $(OUTPUT_DIR)
	@$(GO) test -race -cover \
		-coverprofile=$(OUTPUT_DIR)/coverage.out \
		-timeout=10m -shuffle=on -short \
		-v `go list ./...|egrep -v 'tools|vendor|third_party'`

go.cover: go.test ## 执行单元测试，并校验覆盖率阈值.
	@echo "===========> Running code coverage tests"
	@$(GO) tool cover -func=$(OUTPUT_DIR)/coverage.out | awk -v target=$(COVERAGE) -f $(PROJ_ROOT_DIR)/scripts/coverage.awk

go.lint: tools.verify.golangci-lint ## 执行静态代码检查.
	@echo "===========> Running golangci to lint source codes"
	@golangci-lint run -c $(PROJ_ROOT_DIR)/.golangci.yaml $(PROJ_ROOT_DIR)/...

# 伪目标（防止文件与目标名称冲突）
.PHONY: go.build.verify go.build.% go.build go.format go.tidy go.test go.cover go.lint
