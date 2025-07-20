package auth

import (
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	adapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

const (
	// 默认的cabin 访问控制模型
	defaultAclModel = `[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act, eft

[role_definition]
g = _, _

[policy_effect]
e = !some(where (p.eft == deny))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && r.act == p.act`
)

// 授权器, 提供授权功能
type Authz struct {
	*casbin.SyncedEnforcer // casbin的同步授权器
}

// 函数选项类型, 用于自定义NewAuthz的行为
type Option func(*authzConfig)

// 授权器的配置结构
type authzConfig struct {
	aclModel           string        // casbin的模型字符串
	autoLoadPolicyTime time.Duration // 自动加载策略的时间间隔
}

// 返回一个默认配置
func defaultAuthzConfig() *authzConfig {
	return &authzConfig{
		// 默认使用内置的acl模型
		aclModel: defaultAclModel,
		// 自动加载更新策略时间间隔为5s
		autoLoadPolicyTime: 5 * time.Second,
	}
}

// 允许通过选项自定义ACL模型
func DefaultOptions() []Option {
	return []Option{
		// 使用默认的ACL模型
		WithAclModel(defaultAclModel),
		// 设置自动加载策略的时间间隔
		WithAutoLoadPolicyTime(10 * time.Second),
	}
}

// 允许通过选项自定义 ACL 模型.
func WithAclModel(model string) Option {
	return func(ac *authzConfig) {
		ac.aclModel = model
	}
}

// 允许通过选项自定义自动加载策略的时间间隔.
func WithAutoLoadPolicyTime(interval time.Duration) Option {
	return func(ac *authzConfig) {
		ac.autoLoadPolicyTime = interval
	}
}

// 创建一个使用casbin完成授权的授权器, 通过函数选项模式支持自定义配置
func NewAuthz(db *gorm.DB, opts ...Option) (*Authz, error) {
	// 初始化默认配置
	cfg := defaultAuthzConfig()

	// 应用所有选项
	for _, opt := range opts {
		opt(cfg)
	}

	// 初始化gorm适配器并用于casbin授权器
	adapter, err := adapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	// 从配置中加载casbin模型
	m, _ := model.NewModelFromString(cfg.aclModel)

	// 初始化授权器
	enforcer, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	// 从数据库加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, err
	}

	// 启动自动加载策略, 使用配置的时间间隔
	enforcer.StartAutoLoadPolicy(cfg.autoLoadPolicyTime)

	// 返回新的授权器实例
	return &Authz{enforcer}, nil
}

// 用于进行授权
func (a *Authz) Authorize(sub, obj, act string) (bool, error) {
	// 调用 Enforce 方法进行授权检查某个主体sub是否可以对某个资源obj进行指定操作act
	return a.Enforce(sub, obj, act)
}
