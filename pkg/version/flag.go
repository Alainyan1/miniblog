// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package version

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/pflag"
)

type versionValue int

const (
	// 未设置版本
	VersionNotSet versionValue = 0
	// 启用版本
	VersionEnabled versionValue = 1
	// 原始版本
	VersionRaw versionValue = 2
)

const (
	// 原始版本的字符串
	strRawVersion = "raw"
	// 版本标志的名称
	versionFlagName = "version"
)

// versionFlag定义了版本标志
var versionFlag = Version(versionFlagName, VersionNotSet, "Print version information and quit.")

func (v *versionValue) IsBoolFlag() bool {
	return true
}

func (v *versionValue) Get() interface{} {
	return v
}

// 实现了pflag.Value接口中的String方法
func (v *versionValue) String() string {
	if *v == VersionRaw {
		return strRawVersion // 返回原始版本字符串
	}

	return strconv.FormatBool(bool(*v == VersionEnabled))
}

// 实现了pflag.Value接口中的Set方法
func (v *versionValue) Set(s string) error {
	if s == strRawVersion {
		*v = VersionRaw
		return nil
	}
	boolVal, err := strconv.ParseBool(s)
	if boolVal {
		*v = VersionEnabled
	} else {
		*v = VersionNotSet
	}

	return err
}

// 实现了pflag.Value接口中的Type方法
func (v *versionValue) Type() string {
	return "version"
}

// 定义了一个具有指定名称和用法的标志
func VersionVar(p *versionValue, name string, value versionValue, usage string) {
	*p = value
	pflag.Var(p, name, usage)

	// `--version` 等价于 `--verison=true`
	pflag.Lookup(name).NoOptDefVal = "true"
}

// Version包装了VersionVar函数
func Version(name string, value versionValue, usage string) *versionValue {
	p := new(versionValue)
	VersionVar(p, name, value, usage)
	return p
}

// 给mb-apiserver 命令添加-v/--version 命令行选项
func AddFlags(fs *pflag.FlagSet) {
	fs.AddFlag(pflag.Lookup(versionFlagName))
}

// 用来指定当mb-apiserver命令执行并传入-v/--version命令行选项时, 应用会打印版本信息并推出
// 当执行mb-apiserver --version 命令时, --version命令行选项的值被赋给version包的versionFlag变量
// PrintAndExitIfRequested函数被执行, 根据versionFlag的值获取版本信息
// 例: _output/mb-apiserver --version=raw
func PrintAndExitIfRequested() {
	// 检查版本标志的值并打印相应的信息
	if *versionFlag == VersionRaw {
		fmt.Printf("%s\n", Get().Text())
		os.Exit(0)
	} else if *versionFlag == VersionEnabled {
		fmt.Printf("%s\n", Get().String())
		os.Exit(0)
	}
}
