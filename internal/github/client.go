// Package github 提供 GitHub API 客户端实现
package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// Client 表示 GitHub API 客户端
type Client struct {
	token     string
	userAgent string
}

// NewClient 创建一个新的 GitHub API 客户端
func NewClient(token string) *Client {
	return &Client{
		token:     token,
		userAgent: "issue2md",
	}
}

// ParseURLInfo 包含解析 URL 后的信息
type ParseURLInfo struct {
	Owner       string
	Repo        string
	Number      int
	ContentType string // "issue", "pr", "discussion"
}

// ParseURL 解析 GitHub URL 并提取信息
func ParseURL(inputURL string) (*ParseURLInfo, error) {
	if inputURL == "" {
		return nil, fmt.Errorf("empty URL")
	}

	// 移除 www. 前缀
	inputURL = strings.Replace(inputURL, "www.github.com", "github.com", 1)

	// 匹配 GitHub URL 格式
	// 支持: issues/123, pull/123, discussions/123
	pattern := `^https?://github\.com/([^/]+)/([^/]+)/(issues|pull|discussions)/(\d+)$`
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(inputURL)

	if len(matches) != 5 {
		return nil, fmt.Errorf("invalid GitHub URL: %s", inputURL)
	}

	owner := matches[1]
	repo := matches[2]
	contentType := matches[3]
	number, err := strconv.Atoi(matches[4])
	if err != nil {
		return nil, fmt.Errorf("invalid number in URL: %w", err)
	}

	// 将 URL 类型映射到 API 类型
	var apiType string
	switch contentType {
	case "issues":
		apiType = "issue"
	case "pull":
		apiType = "pr"
	case "discussions":
		apiType = "discussion"
	default:
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}

	return &ParseURLInfo{
		Owner:       owner,
		Repo:        repo,
		Number:      number,
		ContentType: apiType,
	}, nil
}

// BuildAPIURL 构建 GitHub API URL
func BuildAPIURL(owner, repo, contentType string, number int) string {
	baseURL := "https://api.github.com"

	var path string
	switch contentType {
	case "issue":
		path = fmt.Sprintf("/repos/%s/%s/issues/%d", owner, repo, number)
	case "pr":
		path = fmt.Sprintf("/repos/%s/%s/pulls/%d", owner, repo, number)
	case "discussion":
		path = fmt.Sprintf("/repos/%s/%s/discussions/%d", owner, repo, number)
	default:
		path = fmt.Sprintf("/repos/%s/%s/issues/%d", owner, repo, number)
	}

	return baseURL + path
}

// makeRequest 发送 HTTP 请求到 GitHub API
func (c *Client) makeRequest(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// 设置 User-Agent
	req.Header.Set("User-Agent", c.userAgent)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	// 如果提供了 token，添加认证头
	if c.token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("GitHub API returned status %d: %s", resp.StatusCode, string(body))
	}

	return resp, nil
}

// FetchIssue 获取 Issue 数据
func (c *Client) FetchIssue(owner, repo string, number int) (*IssueResponse, error) {
	url := BuildAPIURL(owner, repo, "issue", number)
	resp, err := c.makeRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var issue IssueResponse
	if err := json.Unmarshal(body, &issue); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &issue, nil
}

// FetchPullRequest 获取 PR 数据
func (c *Client) FetchPullRequest(owner, repo string, number int) (*PullRequestResponse, error) {
	url := BuildAPIURL(owner, repo, "pr", number)
	resp, err := c.makeRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var pr PullRequestResponse
	if err := json.Unmarshal(body, &pr); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &pr, nil
}

// FetchDiscussion 获取 Discussion 数据
func (c *Client) FetchDiscussion(owner, repo string, number int) (*DiscussionResponse, error) {
	url := BuildAPIURL(owner, repo, "discussion", number)
	resp, err := c.makeRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var discussion DiscussionResponse
	if err := json.Unmarshal(body, &discussion); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &discussion, nil
}

// FetchIssueComments 获取 Issue 评论列表
func (c *Client) FetchIssueComments(owner, repo string, number int) ([]CommentResponse, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d/comments", owner, repo, number)
	resp, err := c.makeRequest(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var comments []CommentResponse
	if err := json.Unmarshal(body, &comments); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return comments, nil
}

// FetchPRComments 获取 PR 评论列表（包含 Review Comments）
func (c *Client) FetchPRComments(owner, repo string, number int) ([]CommentResponse, []ReviewCommentResponse, error) {
	// 获取普通评论
	commentsURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/issues/%d/comments", owner, repo, number)
	commentsResp, err := c.makeRequest(commentsURL)
	if err != nil {
		return nil, nil, err
	}
	defer commentsResp.Body.Close()

	commentsBody, err := io.ReadAll(commentsResp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read comments response: %w", err)
	}

	var comments []CommentResponse
	if err := json.Unmarshal(commentsBody, &comments); err != nil {
		return nil, nil, fmt.Errorf("failed to parse comments JSON: %w", err)
	}

	// 获取 review comments
	reviewURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls/%d/comments", owner, repo, number)
	reviewResp, err := c.makeRequest(reviewURL)
	if err != nil {
		return nil, nil, err
	}
	defer reviewResp.Body.Close()

	reviewBody, err := io.ReadAll(reviewResp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read review response: %w", err)
	}

	var reviews []ReviewCommentResponse
	if err := json.Unmarshal(reviewBody, &reviews); err != nil {
		return nil, nil, fmt.Errorf("failed to parse review JSON: %w", err)
	}

	return comments, reviews, nil
}

// FetchDiff 获取 PR 的代码 diff
func (c *Client) FetchDiff(diffURL string) (string, error) {
	resp, err := c.makeRequest(diffURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read diff: %w", err)
	}

	return string(body), nil
}
