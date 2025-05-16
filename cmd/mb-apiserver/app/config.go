// Copyright 2024 alainyan <alainyan@yahoo.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/Alainyan1/miniblog.

package app

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	// 服务配置的默认目录
	defaultHomeDir = ".miniblog"

	// 默认配置文件名
	defaultConfigName = "mb-apiserver.yaml"
)

// onInitialize 设置需要读取的配置文件名、环境变量，并将其内容读取到 viper 中.
func OnInitialize() {
	if configFile != "" {
		// 从命令行指定的配置文件中读取
		viper.SetConfigFile(configFile)
	} else {
		// 使用默认配置文件路径和名称
		// 当configFile为空时, 会从当前目录.和$HOME/.miniblog目录下加载名为mb-apiserver.yaml的文件
		for _, dir := range searchDirs() {
			viper.AddConfigPath(dir)
		}

		// 设置配置文件格式为yaml
		viper.SetConfigType("yaml")

		// 配置文件名称(没有文件扩展名)
		viper.SetConfigName(defaultConfigName)
	}

	// 读取环境变量并设置前缀
	setupEnvironmentVariables()

	// 读取配置文件, 如果指定了配置文件名, 则使用指定的配置文件, 否则在注册的搜索路径中搜索
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Fail to read viper configuration file, err: %v", err)
	}

	// 打印当前的配置文件
	log.Printf("Using cofig file: %s", viper.ConfigFileUsed())
}

func setupEnvironmentVariables() {
	// 允许viper自动匹配环境变量
	viper.AutomaticEnv()

	// 当多个服务同是在同一台机器上运行时会出现环境变量名称冲突
	// 设置环境变量前缀
	viper.SetEnvPrefix("MINIBLOG")

	// 替换环境变量中的分割符, ".", "-"为"_"
	replacer := strings.NewReplacer(".", "_", "-", "_")
	viper.SetEnvKeyReplacer(replacer)
}

func filePath() string {
	return ""
}

// 返回默认配置的文件搜索目录
func searchDirs() []string {
	// 用户主目录
	homeDir, err := os.UserHomeDir()
	// 如果获取用户主目录失败, 则打印错误信息并退出
	cobra.CheckErr(err)
	// 返回当前目录.和$HOME/.miniblog目录
	return []string{filepath.Join(homeDir, defaultHomeDir), "."}
}
