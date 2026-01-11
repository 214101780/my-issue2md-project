// Package models 定义核心数据模型
package models

// ResourceType 表示 GitHub 资源类型
type ResourceType int

const (
	// ResourceTypeIssue 表示 GitHub Issue
	ResourceTypeIssue ResourceType = iota
	// ResourceTypePullRequest 表示 GitHub Pull Request
	ResourceTypePullRequest
	// ResourceTypeDiscussion 表示 GitHub Discussion
	ResourceTypeDiscussion
)

// String 返回资源类型的字符串表示
func (r ResourceType) String() string {
	switch r {
	case ResourceTypeIssue:
		return "issue"
	case ResourceTypePullRequest:
		return "pull"
	case ResourceTypeDiscussion:
		return "discussion"
	default:
		return "unknown"
	}
}
