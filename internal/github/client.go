// Package github 提供 GitHub API 客户端
package github

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Client 表示 GitHub API 客户端
type Client struct {
	httpClient *http.Client
	token      string
	baseURL    string
}

// NewClient 创建新的 GitHub API 客户端
func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{},
		token:      token,
		baseURL:    "https://api.github.com/graphql",
	}
}

// SetHTTPClient 设置自定义 HTTP 客户端（用于测试）
func (c *Client) SetHTTPClient(client *http.Client) {
	c.httpClient = client
}

// SetBaseURL 设置自定义基础 URL（用于测试）
func (c *Client) SetBaseURL(url string) {
	c.baseURL = url
}

// GetIssue 获取 Issue 信息
func (c *Client) GetIssue(ctx context.Context, owner, repo string, number int) (*IssueNode, error) {
	// 构建 GraphQL 查询
	query := fmt.Sprintf(`{
		"query": "query GetIssue($owner: String!, $repo: String!, $number: Int!) {
			repository(owner: $owner, name: $repo) {
				issue(number: $number) {
					title
					number
					url
					state
					body
					createdAt
					updatedAt
					author {
						login
						avatarUrl
					}
					reactions {
						totalCount
						nodes {
							content
						}
					}
					comments(first: 100) {
						totalCount
						nodes {
							id
							body
							createdAt
							updatedAt
							author {
								login
								avatarUrl
							}
							reactions {
								totalCount
								nodes {
									content
								}
							}
						}
					}
				}
			}
		}",
		"variables": {
			"owner": %q,
			"repo": %q,
			"number": number
		}
	}`, owner, repo)

	// 创建请求
	reqBody := []byte(query)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	// 发送请求
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var graphqlResp GraphQLResponse
	if err := json.NewDecoder(resp.Body).Decode(&graphqlResp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	// 检查 GraphQL 错误
	if len(graphqlResp.Errors) > 0 {
		return nil, fmt.Errorf("GraphQL error: %s", graphqlResp.Errors[0].Message)
	}

	// 检查 Issue 是否存在
	if graphqlResp.Data == nil || graphqlResp.Data.Repository == nil || graphqlResp.Data.Repository.Issue == nil {
		return nil, fmt.Errorf("issue not found")
	}

	return graphqlResp.Data.Repository.Issue, nil
}
