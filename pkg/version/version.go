// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package version

import (
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/gosuri/uitable"
)

// version包用于记录版本号信息, 版本号功能几乎所有go应用都会用到, 因此需要将version包提供给其他外部应用程序使用
// 根据目录规范, 将version包放在pkg目录下

var (
	// 实际使用中, gitVersion通常会通过 -ldflags参数在编译时赋值为实际版本.
	gitVersion = "v0.0.0-master+$Format:%h$"
	// 实际使用中, gitVersion通常会通过 -ldflags参数在编译时赋值为构建时的时间戳.
	buildDate = "1970-01-01T00:00:00Z"
	// gitCommit是git的SHA1值, git rev-prase HEAD 命令的输出.
	gitCommit = "$Format:%H$"
	// gitTreeState代表构建时git仓库的状态, 可能值为clean, dirty.
	gitTreeState = ""
)

// Info包含了版本信息.
type Info struct {
	GitVersion   string `json:"gitVersion"`
	GitCommit    string `json:"gitCommit"`
	GitTreeState string `json:"gitTreeState"`
	BuildDate    string `json:"buildDate"`
	GoVersion    string `json:"goVersion"`
	Compiler     string `json:"compiler"`
	Platform     string `json:"platform"`
}

// String返回友好可读版本信息字符串.
func (info Info) String() string {
	return info.GitVersion
}

// ToJSON以json格式返回版本信息.
func (info Info) ToJSON() string {
	s, _ := json.Marshal(info)

	return string(s)
}

// Text将版本信息编码为utf-8格式的文本并返回.
func (info Info) Text() string {
	table := uitable.New()
	table.RightAlign(0)
	table.MaxColWidth = 80
	table.Separator = " "
	table.AddRow("gitVersion:", info.GitVersion)
	table.AddRow("gitCommit:", info.GitCommit)
	table.AddRow("gitTreeState:", info.GitTreeState)
	table.AddRow("buildDate:", info.BuildDate)
	table.AddRow("goVersion:", info.GoVersion)
	table.AddRow("compiler:", info.Compiler)
	table.AddRow("platform:", info.Platform)

	return table.String()
}

// Get返回详尽的代码版本库信息, 用来表明二进制文件由哪个版本的代码构建.
func Get() Info {
	return Info{
		// 以下变量通常由 -ldflags进行设置
		// GoVersion, Compiler, Platform可以使用runtime包来动态获取
		GitVersion:   gitVersion,
		GitCommit:    gitCommit,
		GitTreeState: gitTreeState,
		BuildDate:    buildDate,
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
