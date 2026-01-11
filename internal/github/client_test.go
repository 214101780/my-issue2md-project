// Package github 的单元测试
package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetIssue_Success 测试成功获取 Issue
func TestGetIssue_Success(t *testing.T) {
	// 创建 Mock Server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求方法
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		// 验证 Content-Type
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}

		// 返回模拟的 GraphQL 响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockIssueResponse))
	}))
	defer mockServer.Close()

	// 创建客户端并配置使用 Mock Server
	client := NewClient("test-token")
	client.SetBaseURL(mockServer.URL)
	client.SetHTTPClient(mockServer.Client())

	// 调用 GetIssue
	ctx := context.Background()
	issue, err := client.GetIssue(ctx, "bigwhite", "issue2md", 1)

	// 验证无错误
	if err != nil {
		t.Fatalf("GetIssue() error = %v", err)
	}

	// 验证返回的数据
	if issue == nil {
		t.Fatal("GetIssue() returned nil issue")
	}

	// 验证标题
	if issue.Title != "Test Issue Title" {
		t.Errorf("expected title 'Test Issue Title', got '%s'", issue.Title)
	}

	// 验证编号
	if issue.Number != 1 {
		t.Errorf("expected number 1, got %d", issue.Number)
	}

	// 验证状态
	if issue.State != "OPEN" {
		t.Errorf("expected state 'OPEN', got '%s'", issue.State)
	}

	// 验证作者
	if issue.Author == nil {
		t.Error("expected author, got nil")
	} else {
		if issue.Author.Login != "testuser" {
			t.Errorf("expected author login 'testuser', got '%s'", issue.Author.Login)
		}
	}

	// 验证评论
	if issue.Comments == nil {
		t.Error("expected comments, got nil")
	} else if issue.Comments.TotalCount != 2 {
		t.Errorf("expected 2 comments, got %d", issue.Comments.TotalCount)
	}
}

// TestGetIssue_NotFound 测试 Issue 不存在
func TestGetIssue_NotFound(t *testing.T) {
	// 创建返回 404 的 Mock Server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": {"repository": {"issue": null}}}`))
	}))
	defer mockServer.Close()

	client := NewClient("test-token")
	client.SetBaseURL(mockServer.URL)

	ctx := context.Background()
	issue, err := client.GetIssue(ctx, "bigwhite", "issue2md", 999)

	// 应该返回错误
	if err == nil {
		t.Error("expected error for non-existent issue, got nil")
	}

	// Issue 应该为 nil
	if issue != nil {
		t.Errorf("expected nil issue for not found, got %+v", issue)
	}
}

// TestGetIssue_APIError 测试 API 返回错误
func TestGetIssue_APIError(t *testing.T) {
	// 创建返回 GraphQL 错误的 Mock Server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"data": null,
			"errors": [
				{
					"message": "Resource not found",
					"type": "NOT_FOUND"
				}
			]
		}`))
	}))
	defer mockServer.Close()

	client := NewClient("test-token")
	client.SetBaseURL(mockServer.URL)

	ctx := context.Background()
	issue, err := client.GetIssue(ctx, "bigwhite", "issue2md", 999)

	// 应该返回错误
	if err == nil {
		t.Error("expected error for API error, got nil")
	}

	// Issue 应该为 nil
	if issue != nil {
		t.Errorf("expected nil issue for API error, got %+v", issue)
	}
}

// mockIssueResponse 是模拟的 Issue GraphQL 响应
const mockIssueResponse = `{
	"data": {
		"repository": {
			"issue": {
				"title": "Test Issue Title",
				"number": 1,
				"url": "https://github.com/bigwhite/issue2md/issues/1",
				"state": "OPEN",
				"body": "This is a test issue body.",
				"createdAt": "2024-01-01T10:00:00Z",
				"updatedAt": "2024-01-02T15:30:00Z",
				"author": {
					"login": "testuser",
					"avatarUrl": "https://github.com/images/error/testuser_happy.gif"
				},
				"reactions": {
					"totalCount": 5,
					"nodes": [
						{"content": "THUMBS_UP"},
						{"content": "THUMBS_UP"},
						{"content": "LAUGH"},
						{"content": "HOORAY"},
						{"content": "HEART"}
					]
				},
				"comments": {
					"totalCount": 2,
					"nodes": [
						{
							"id": "comment_1",
							"body": "First comment",
							"createdAt": "2024-01-01T11:00:00Z",
							"updatedAt": "2024-01-01T11:00:00Z",
							"author": {
								"login": "commenter1",
								"avatarUrl": "https://github.com/images/error/commenter1_happy.gif"
							},
							"reactions": {
								"totalCount": 1,
								"nodes": [{"content": "THUMBS_UP"}]
							}
						},
						{
							"id": "comment_2",
							"body": "Second comment",
							"createdAt": "2024-01-01T12:00:00Z",
							"updatedAt": "2024-01-01T12:00:00Z",
							"author": {
								"login": "commenter2",
								"avatarUrl": "https://github.com/images/error/commenter2_happy.gif"
							},
							"reactions": {
								"totalCount": 0,
								"nodes": []
							}
						}
					]
				}
			}
		}
	}
}`
