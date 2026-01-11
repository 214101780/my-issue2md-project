// Package config 管理应用配置
package config

import "time"

// Options 表示命令行选项
type Options struct {
	EnableReactions bool // 包含 reactions 统计信息
	EnableUserLinks bool // 将 @username 转换为链接
}

// Config 是应用配置的完整结构
type Config struct {
	GitHubToken string        // GitHub Personal Access Token
	TargetURL   string        // 目标 GitHub URL
	OutputFile  string        // 输出文件路径，空字符串表示 stdout
	Options     Options       // 功能开关
	APITimeout  time.Duration // API 请求超时时间
}

// NewConfig 创建默认配置
func NewConfig() *Config {
	return &Config{
		Options: Options{
			EnableReactions: false,
			EnableUserLinks: false,
		},
		APITimeout: 30 * time.Second,
	}
}
