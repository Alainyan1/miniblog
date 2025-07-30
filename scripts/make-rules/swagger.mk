# ==============================================================================
# Makefile helper functions for swagger
# 用于生成和提供 Swagger API 文档
# 

swagger.run: tools.verify.swagger
	@echo "===========> Generating swagger API docs"
	@# mixin将多个swagger json合并成一个文件
	@# -q: 静默模式, 减少命令输出
	@# --keep-spec-order:保留输入文件的顺序
	@# --format=yaml: 输出格式为 YAML
	@# --ignore-conflicts: 忽略合并时的冲突, 如重复的 API 路径
	@swagger mixin `find $(PROJ_ROOT_DIR)/api/openapi -name "*.swagger.json"` \
		-q                                                    \
		--keep-spec-order                                     \
		--format=yaml                                         \
		--ignore-conflicts                                    \
		-o $(PROJ_ROOT_DIR)/api/openapi/apiserver/v1/openapi.yaml
	@echo "Generated at: $(PROJ_ROOT_DIR)/api/openapi/apiserver/v1/openapi.yaml"

# swagger serve用于启动一个 Web 服务器来提供 Swagger API 文档
# -F=redoc 使用 ReDoc 渲染器(一种美观的 OpenAPI 文档渲染工具)显示文档
# --no-open: 防止自动打开浏览器
swagger.serve: tools.verify.swagger
	@swagger serve -F=redoc --no-open --port 65534 $(PROJ_ROOT_DIR)/api/openapi/apiserver/v1/openapi.yaml

# 伪目标（防止文件与目标名称冲突）
.PHONY: swagger.run swagger.serve