// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package log

import (
	"context"
	"miniblog/internal/pkg/contextx"
	"miniblog/internal/pkg/known"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 因为该包封装了一些定制化的逻辑, 不适合对外暴露, 不适合放在/pkg目录下
// 但日志包又是项目内的共享包, 所以存放在internal/pkg下

// 通常将日志接口命名为Logger.
type Logger interface {
	// 记录调试级别的日志, 通常用于开发阶段, 包含详细的调试信息
	Debugw(msg string, kvs ...any)

	// 记录信息级别的日志, 表示系统正常运行状态
	Infow(msg string, kvs ...any)

	// 警告级别日志, 可能存在问题但不影响系统正常运行
	Warnw(msg string, kvs ...any)

	// 错误级别日志, 系统运行中出现的错误, 需要开发人员介入处理
	Errorw(msg string, kvs ...any)

	// 严重错误级别的日志, 系统无法继续运行, 记录日志后会发生panic
	Panicw(msg string, kvs ...any)

	// 致命错误级别日志, 系统无法继续运行, 记录日志后会直接推出程序
	Fatalw(msg string, kvs ...any)

	// 用于刷新日志缓冲区, 确保日志被完整写入目标存储
	Sync()
}

// 定义为不可导出类型, 有利于日志类型的封装和维护, 屏蔽实现细节.
type zapLogger struct {
	z *zap.Logger
}

// 这是一个类型断言的语法, 用于告诉编译器检查 *zapLogger 类型是否满足某个条件(在这里是是否实现 Logger 接口).
var _ Logger = (*zapLogger)(nil)

var (
	mu sync.Mutex
	// std 定义了默认的全局 Logger.
	std = New(NewOptions())
)

// 初始化全局的日志对象.
func Init(opts *Options) {
	// 因为会给全局变量std赋值, 对std变量加锁
	mu.Lock()
	defer mu.Unlock()

	std = New(opts)
}

// 一个日志包通常有两类zapLogger对象, 全局对象和局部对象
// 全局对象便于通过类似log.Infow()方式直接调用, 局部对象方便传入不同参数以创建自定义的Logger
// 所以需要实现New和Init两种初始化参数

// 如果Options参数为空, 则会使用默认的Options配置.
func New(opts *Options) *zapLogger {
	if opts == nil {
		opts = NewOptions()
	}

	// 将Options中的日志级别(string)转换为 zapcore.Level 类型
	var zapLevel zapcore.Level
	if err := zapLevel.UnmarshalText([]byte(opts.Level)); err != nil {
		// 如果指定了非法日志, 默认使用Info级别
		zapLevel = zapcore.InfoLevel
	}

	// 创建encode配置, 用于控制日志的输出格式
	encoderConfig := zap.NewProductionEncoderConfig()
	// 自定义 MessageKey 为 message, message 语义更明确
	encoderConfig.MessageKey = "message"
	// 自定义 TimeKey 为 timestamp, timestamp 语义更明确
	encoderConfig.TimeKey = "timestamp"
	// 指定时间序列化函数, 将时间序列化为 `2006-01-02 15:04:05.000` 格式, 更易读
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
	}
	// 指定 time.Duration 序列化函数, 将 time.Duration 序列化为经过的毫秒数的浮点数
	encoderConfig.EncodeDuration = func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendFloat64(float64(d) / float64(time.Millisecond))
	}

	// 创建zap.Logger需要的配置
	// 用于配置zap日志记录器的行为
	cfg := &zap.Config{
		DisableCaller:     opts.DisableCaller,
		DisableStacktrace: opts.DisableStacktrace,
		Level:             zap.NewAtomicLevelAt(zapLevel),
		Encoding:          opts.Format,
		EncoderConfig:     encoderConfig,
		OutputPaths:       opts.OutputPaths,
		ErrorOutputPaths:  []string{"stderr"},
	}

	// 使用cfg创建 *zap.Logger对象
	// 告诉zap在日志级别达到或超过PanicLevel时, 自动附加调用栈信息
	// zap被封装在其他包中, 调用栈会包含额外的层级
	// 由于log包对zap包实现了封装, 在调用栈中需要跳过的调用深度应该加2
	// 使用zap.AddCallerSkip(2)告诉zap在捕获调用者信息时, 向上跳过2层调用栈, 从而记录更上层的调用者信息
	z, err := cfg.Build(zap.AddStacktrace(zapcore.PanicLevel), zap.AddCallerSkip(2))
	if err != nil {
		panic(err)
	}

	// 将标准库的log输出重定向到zap.Logger
	zap.RedirectStdLog(z)

	return &zapLogger{z: z}
}

// 实现分层封装
// 全局 Debugw 函数:
// 全局 Debugw 函数（func Debugw(msg string, kvs ...any)）是为整个应用程序提供一个统一的, 简洁的日志记录入口
// 它的作用类似于一个"门面", 让开发者可以在代码的任何地方直接调用 Debugw("message", "key", "value"),无需关心底层的日志实现细节(如 zap 库或 zapLogger 结构体)
// 它隐藏了 std(全局日志实例)的具体实现, 降低了调用者的使用复杂度
// zapLogger 的 Debugw 方法:
// zapLogger 结构体的 Debugw 方法（func (l *zapLogger) Debugw(...)）是为特定的日志实例提供接口, 允许在不同模块或上下文中使用独立的日志记录器
// 它是对 zap.Logger 的直接封装, 定义了如何使用底层的 zap 库(例如通过 l.z.Sugar().Debugw)来记录日志
// 这个方法允许 zapLogger 在未来添加自定义逻辑(例如附加默认字段, 过滤日志等)
// 方便为不同模块提供独立的日志配置

// 调用的是底层zap.Logger的Sync方法, 将缓存中的日志刷新到磁盘文件中, 主程序需要在推出前调用Sync.
func Sync() {
	std.Sync()
}

func (l *zapLogger) Sync() {
	_ = l.z.Sync()
}

// 全局Debugw函数, 提供一个简化的全局接口, 方便在代码中记录debug级别的日志, 而无需直接访问底层的zap日志记录器.
func Debugw(msg string, kvs ...any) {
	std.Debugw(msg, kvs...)
}

// Sugar() 是 zap 提供的一个更简洁的 API, 允许使用变长参数(如键值对)来记录结构化日志, 而无需手动构造 zap.Field.
func (l *zapLogger) Debugw(msg string, kvs ...any) {
	l.z.Sugar().Debugw(msg, kvs...)
}

func Infow(msg string, kvs ...any) {
	std.Infow(msg, kvs...)
}

func (l *zapLogger) Infow(msg string, kvs ...any) {
	l.z.Sugar().Infow(msg, kvs...)
}

func Warnw(msg string, kvs ...any) {
	std.Warnw(msg, kvs...)
}

func (l *zapLogger) Warnw(msg string, kvs ...any) {
	l.z.Sugar().Warnw(msg, kvs...)
}

func Errorw(msg string, kvs ...any) {
	std.Errorw(msg, kvs...)
}

func (l *zapLogger) Errorw(msg string, kvs ...any) {
	l.z.Sugar().Errorw(msg, kvs...)
}

func Panicw(msg string, kvs ...any) {
	std.Panicw(msg, kvs...)
}

func (l *zapLogger) Panicw(msg string, kvs ...any) {
	l.z.Sugar().Panicw(msg, kvs...)
}

func Fatalw(msg string, kvs ...any) {
	std.Fatalw(msg, kvs...)
}

func (l *zapLogger) Fatalw(msg string, kvs ...any) {
	l.z.Sugar().Fatalw(msg, kvs...)
}

// W解析传入的context, 尝试提取关注的键值, 并添加到日志中.
func W(ctx context.Context) Logger {
	return std.W(ctx)
}

// W是WithContext的简称, 缩短函数名可以减小日志打印代码行的宽度, 减小日志行代码的折行概率.
func (l *zapLogger) W(ctx context.Context) Logger {
	lc := l.clone()

	// 定义一个映射, 关联context提取函数和日志字段名
	contextExtractors := map[string]func(context.Context) string{
		known.XRequestID: contextx.RequestID,
		known.XUserID:    contextx.UserID,
	}

	// 变量映射, 从context中提取值并添加到日志中
	for fieldName, extractor := range contextExtractors {
		if val := extractor(ctx); val != "" {
			lc.z = lc.z.With(zap.String(fieldName, val))
		}
	}

	return lc
}

// 由于log包会被多个请求并发调用, 为了防止请求ID被污染, 每个请求都会对log包深拷贝一个*zapLogger对象, 然后再添加请求id.
func (l *zapLogger) clone() *zapLogger {
	newLogger := *l
	return &newLogger
}
