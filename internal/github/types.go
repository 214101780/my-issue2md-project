// Package github 定义 GitHub API 响应结构
package github

// GraphQLResponse 是 GraphQL API 的通用响应包装
type GraphQLResponse struct {
	Data   *GraphQLData   `json:"data"`
	Errors []GraphQLError `json:"errors,omitempty"`
}

// GraphQLError 表示 GraphQL 错误
type GraphQLError struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// GraphQLData 包含所有可能的查询结果
type GraphQLData struct {
	Repository *Repository `json:"repository"`
}

// Repository 表示 GitHub 仓库
type Repository struct {
	Issue      *IssueNode       `json:"issue,omitempty"`
	PullRequest *PullRequestNode `json:"pullRequest,omitempty"`
	Discussion *DiscussionNode  `json:"discussion,omitempty"`
}

// IssueNode 表示 Issue 节点
type IssueNode struct {
	Title     string              `json:"title"`
	Number    int                 `json:"number"`
	URL       string              `json:"url"`
	State     string              `json:"state"`
	Body      string              `json:"body"`
	CreatedAt string              `json:"createdAt"`
	UpdatedAt string              `json:"updatedAt"`
	Author    *AuthorNode         `json:"author"`
	Reactions *ReactionConnection `json:"reactions"`
	Comments  *CommentConnection  `json:"comments"`
}

// PullRequestNode 表示 PR 节点
type PullRequestNode struct {
	Title     string              `json:"title"`
	Number    int                 `json:"number"`
	URL       string              `json:"url"`
	State     string              `json:"state"`
	Merged    bool                `json:"merged"`
	Body      string              `json:"body"`
	CreatedAt string              `json:"createdAt"`
	UpdatedAt string              `json:"updatedAt"`
	Author    *AuthorNode         `json:"author"`
	Reactions *ReactionConnection `json:"reactions"`
	Comments  *CommentConnection  `json:"comments"`
}

// DiscussionNode 表示 Discussion 节点
type DiscussionNode struct {
	Title     string              `json:"title"`
	Number    int                 `json:"number"`
	URL       string              `json:"url"`
	State     string              `json:"state"`
	Body      string              `json:"body"`
	CreatedAt string              `json:"createdAt"`
	UpdatedAt string              `json:"updatedAt"`
	Author    *AuthorNode         `json:"author"`
	Reactions *ReactionConnection `json:"reactions"`
	Comments  *CommentConnection  `json:"comments"`
	Answer    *CommentNode        `json:"answer"`
}

// AuthorNode 表示 API 返回的作者节点
type AuthorNode struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatarUrl"`
}

// ReactionConnection 表示 reactions 连接
type ReactionConnection struct {
	TotalCount int            `json:"totalCount"`
	Nodes      []ReactionNode `json:"nodes"`
}

// ReactionNode 表示单个 reaction
type ReactionNode struct {
	Content string `json:"content"`
}

// CommentConnection 表示评论连接
type CommentConnection struct {
	TotalCount int           `json:"totalCount"`
	Nodes      []CommentNode `json:"nodes"`
}

// CommentNode 表示单个评论节点
type CommentNode struct {
	ID        string              `json:"id"`
	Body      string              `json:"body"`
	CreatedAt string              `json:"createdAt"`
	UpdatedAt string              `json:"updatedAt"`
	Author    *AuthorNode         `json:"author"`
	Reactions *ReactionConnection `json:"reactions"`
}
