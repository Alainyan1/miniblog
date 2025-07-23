// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package main

// 连接musql数据库, 根据数据库表结构自动生成对应的go结构体(gorm模型), 生成查询接口和相关的curd操作代码.
import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/onexstack/onexstack/pkg/db"
	"github.com/spf13/pflag"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

// 帮助文本信息.
const helpText = `Usage: main [flags] arg [arg...]
This is a pflag example. 
Flags: 
`

// 自定义查询接口, 可以为生成的模型添加特定的查询功能.
type Querier interface {
	FilterWithNameAndRole(name string) ([]gen.T, error)
}

// 保存代码生成的配置.
type GenerateConfig struct {
	ModelPackagePath string                 // 生成模型的包路径
	GenerateFunc     func(g *gen.Generator) // 具体生成函数
}

// 预定义生成的配置.
var generateConfigs = map[string]GenerateConfig{
	"mb": {ModelPackagePath: "../../internal/apiserver/model", GenerateFunc: GenerateMiniBlogModels},
}

// 命令行参数.
var (
	addr       = pflag.StringP("addr", "a", "127.0.0.1:3306", "MySQL host address.")
	username   = pflag.StringP("username", "u", "miniblog", "Username to connect to the database.")
	password   = pflag.StringP("password", "p", "miniblog1234", "Password to use when connecting to the database.")
	database   = pflag.StringP("db", "d", "miniblog", "Database name to connect to.")
	modelPath  = pflag.String("model-pkg-path", "", "Generated model code's package name.")
	components = pflag.StringSlice("component", []string{"mb"}, "Generated model code's for specified component.")
	help       = pflag.BoolP("help", "h", false, "Show this help message.")
)

// 3. 处理组件: 遍历指定的组件, 为每个组件生成代码.
func main() {
	// 设置自定义的使用说明函数
	pflag.Usage = func() {
		fmt.Printf("%s", helpText)
		pflag.PrintDefaults()
	}
	pflag.Parse()

	// 如果有帮助标志, 则显示帮助标志并退出
	if *help {
		pflag.Usage()
		return
	}

	// 初始化数据库连接
	dbInstance, err := initializeDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// 处理组件并生成代码
	for _, component := range *components {
		processComponent(component, dbInstance)
	}
}

// 创建并返回一个数据库连接.
func initializeDatabase() (*gorm.DB, error) {
	// NewMySQL需要传入指针类型
	dpOptions := &db.MySQLOptions{
		Addr:     *addr,
		Username: *username,
		Password: *password,
		Database: *database,
	}

	return db.NewMySQL(dpOptions)
}

// 处理单个组件以生成代码.
func processComponent(component string, dbInstance *gorm.DB) {
	config, ok := generateConfigs[component]
	if !ok {
		log.Printf("Component '%s' not found in configuration. Skipping.", component)
		return
	}

	// 解析包模型
	modelPkgPath := resolveModelPackagePath(config.ModelPackagePath)

	// 创建生成器实例
	generator := createGenerator(modelPkgPath)
	generator.UseDB(dbInstance)

	// 应用自定义生成器选项
	applyGeneratorOptions(generator)

	// 使用指定的函数生成器
	config.GenerateFunc(generator)

	// 执行代码生成
	generator.Execute()
}

// 确定模型生成的包路径.
func resolveModelPackagePath(defaultPath string) string {
	// 如果命令行参数中指定了模型包路径, 则使用该路径
	if *modelPath != "" {
		return *modelPath
	}
	absPath, err := filepath.Abs(defaultPath)
	if err != nil {
		log.Printf("Error resolving path: %v", err)
		return defaultPath
	}
	return absPath
}

// 初始化并返回一个新的生成器实例.
func createGenerator(packagePath string) *gen.Generator {
	return gen.NewGenerator(gen.Config{
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext, // 使用默认查询模式, 包含查询接口, 不使用上下文
		ModelPkgPath:      packagePath,
		WithUnitTest:      true,  // 生成单元测试代码
		FieldNullable:     true,  // 对数据库中可空的字段, 使用指针类型
		FieldSignable:     false, // 禁止无符号属性以提高兼容性
		FieldWithIndexTag: false, // 不包含GORM的索引标签
		FieldWithTypeTag:  false, // 不包含GORM的类型标签
	})
}

// 设置自定义生成器选项.
func applyGeneratorOptions(g *gen.Generator) {
	// 为特定字段自定义gorm标签
	g.WithOpts(
		// gorm.io/gen默认生成时间gorm的标签为`gorm:"default:current_timestamp()"`, 由于不能被SQLite识别, 这里自定义为`default:current_timestamp`
		gen.FieldGORMTag("createdAt", func(tag field.GormTag) field.GormTag {
			tag.Set("default", "current_timestamp")
			return tag
		}),
		gen.FieldGORMTag("updatedAt", func(tag field.GormTag) field.GormTag {
			tag.Set("default", "current_timestamp")
			return tag
		}),
	)
}

func GenerateMiniBlogModels(g *gen.Generator) {
	g.GenerateModelAs(
		"user",                         // 生成用户模型, 数据库表名为"user"
		"UserM",                        // 生成的结构体为"userM"
		gen.FieldIgnore("placeholder"), // 忽略占位符字段
		gen.FieldGORMTag("username", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "idx_user_username")
			return tag
		}),
		// 为userID和phone字段添加唯一索引
		gen.FieldGORMTag("userID", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "idx_user_userID")
			return tag
		}),
		gen.FieldGORMTag("phone", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "idx_user_phone")
			return tag
		}),
	)
	// 生成post模型, 数据库表名为"post", 生成的结构体为"PostM"
	g.GenerateModelAs(
		"post",
		"PostM",
		gen.FieldIgnore("placeholder"),
		gen.FieldGORMTag("postID", func(tag field.GormTag) field.GormTag {
			tag.Set("uniqueIndex", "idx_post_postID")
			return tag
		}),
	)
	// 生成CasbinRule模型(权限管理), 数据库表名为"casbin_rule", 生成的结构体为"CasbinRuleM"
	g.GenerateModelAs(
		"casbin_rule",
		"CasbinRuleM",
		gen.FieldRename("ptype", "PType"), // 为了符合 Go 命名规范, 将字段名从ptype改为PType
		gen.FieldIgnore("placeholder"),
	)
}
