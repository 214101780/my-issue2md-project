// Package models 定义 issue2md 的核心数据模型
package models

// ContentType 表示 GitHub 内容的类型
type ContentType string

const (
	ContentTypeIssue       ContentType = "issue"
	ContentTypePullRequest ContentType = "pr"
	ContentTypeDiscussion  ContentType = "discussion"
)

// ItemState 表示内容的状态
type ItemState string

const (
	StateOpen   ItemState = "open"
	StateClosed ItemState = "closed"
	StateMerged ItemState = "merged" // 仅用于 PR
)

// Label 表示 GitHub 的标签
type Label struct {
	Name        string
	Color       string
	Description string
}

// User 表示 GitHub 用户
type User struct {
	Login     string
	ID        int64
	AvatarURL string
	Type      string // User, Bot, etc.
}

// Comment 表示一条评论或回复
type Comment struct {
	ID        int64
	Author    User
	Body      string
	CreatedAt string // ISO 8601 格式
	UpdatedAt string // ISO 8601 格式
	// 用于嵌套评论
	ParentID  *int64 // 父评论 ID，如果是顶级评论则为 nil
	Replies   []Comment
}

// Item 表示一个 Issue、PR 或 Discussion
// 使用统一的结构简化处理
type Item struct {
	// 基础元数据
	Type        ContentType
	Number      int
	Title       string
	Body        string
	State       ItemState
	Author      User
	CreatedAt   string // ISO 8601 格式
	UpdatedAt   string // ISO 8601 格式
	ClosedAt    *string // ISO 8601 格式，如果未关闭则为 nil
	MergedAt    *string // ISO 8601 格式，仅 PR 使用
	URL         string
	Repository  string // "owner/repo" 格式

	// 额外元数据
	Labels      []Label
	Milestone   *string
	Assignees   []User

	// PR 特有字段
	IsPR        bool
	Additions   int    // 代码增加行数
	Deletions   int    // 代码删除行数
	ChangedFiles int   // 修改文件数
	Diff        string // 完整的 diff 内容

	// Discussion 特有字段
	Category    *string // Discussion 分类

	// 评论（包含普通评论和 Review Comments）
	Comments    []Comment
}

// RepositoryInfo 从 URL 解析出的仓库信息
type RepositoryInfo struct {
	Owner string
	Repo  string
}

// ConvertConfig 转换配置选项
type ConvertConfig struct {
	OutputPath string // 输出文件路径，为空则自动生成
	Token      string // GitHub Personal Access Token（可选）
}
