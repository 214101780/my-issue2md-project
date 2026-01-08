// Package github 提供 GitHub API 的响应类型定义
package github

// IssueResponse represents a GitHub Issue API response
type IssueResponse struct {
	ID        int64 `json:"id"`
	Number    int   `json:"number"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	State     string `json:"state"`
	User      User   `json:"user"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	ClosedAt  *string `json:"closed_at"`
	URL       string `json:"html_url"`
	Labels    []Label `json:"labels"`
	Milestone *Milestone `json:"milestone"`
	Assignees []User `json:"assignees"`
	PullRequest *PullRequestInfo `json:"pull_request"`
}

// PullRequestResponse represents a GitHub PR API response
type PullRequestResponse struct {
	ID        int64 `json:"id"`
	Number    int   `json:"number"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	State     string `json:"state"`
	User      User   `json:"user"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	ClosedAt  *string `json:"closed_at"`
	MergedAt  *string `json:"merged_at"`
	URL       string `json:"html_url"`
	Labels    []Label `json:"labels"`
	Milestone *Milestone `json:"milestone"`
	Assignees []User `json:"assignees"`
	Additions int `json:"additions"`
	Deletions int `json:"deletions"`
	ChangedFiles int `json:"changed_files"`
	DiffURL   string `json:"diff_url"`
}

// DiscussionResponse represents a GitHub Discussion API response
type DiscussionResponse struct {
	ID        int64 `json:"id"`
	Number    int   `json:"number"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	ClosedAt  *string `json:"closed_at"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	URL       string `json:"url"`
	Author    Author `json:"author"`
	Category  Category `json:"category"`
	Labels    []Label `json:"labels"`
}

// CommentResponse represents a GitHub Issue/PR comment
type CommentResponse struct {
	ID        int64 `json:"id"`
	User      User `json:"user"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ReviewCommentResponse represents a GitHub PR review comment
type ReviewCommentResponse struct {
	ID        int64 `json:"id"`
	User      User `json:"user"`
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	Position  *int `json:"position,omitempty"`
	Path      *string `json:"path,omitempty"`
}

// User represents a GitHub user
type User struct {
	Login     string `json:"login"`
	ID        int64  `json:"id"`
	AvatarURL string `json:"avatar_url"`
	Type      string `json:"type"`
}

// Author represents a Discussion author (slightly different from User)
type Author struct {
	Login string `json:"login"`
}

// Label represents a GitHub label
type Label struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

// Milestone represents a GitHub milestone
type Milestone struct {
	Title string `json:"title"`
}

// PullRequestInfo indicates if an issue is a PR
type PullRequestInfo struct {
	URL string `json:"url"`
}

// Category represents a Discussion category
type Category struct {
	ID string `json:"id"`
	Slug string `json:"slug"`
}

// CommentThreadResponse represents nested comments in Discussions
type CommentThreadResponse struct {
	Comment CommentResponse `json:"comment"`
	Replies []CommentResponse `json:"replies"`
}
