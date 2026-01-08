// Package converter 提供将 GitHub 数据转换为 Markdown 的功能
package converter

import (
	"strings"
	"testing"

	"github.com/bigwhite/issue2md/internal/models"
)

// TestGenerateFrontMatter 测试 YAML Front Matter 生成
func TestGenerateFrontMatter(t *testing.T) {
	item := &models.Item{
		Type:       models.ContentTypeIssue,
		Number:     123,
		Title:      "Test Issue",
		State:      models.StateOpen,
		Author:     models.User{Login: "testuser"},
		CreatedAt:  "2024-01-01T00:00:00Z",
		UpdatedAt:  "2024-01-02T00:00:00Z",
		URL:        "https://github.com/owner/repo/issues/123",
		Repository: "owner/repo",
		Labels: []models.Label{
			{Name: "bug", Color: "d73a4a"},
			{Name: "enhancement", Color: "a2eeef"},
		},
	}

	want := `---
title: "Test Issue"
author: "testuser"
number: 123
type: "issue"
state: "open"
labels: ["bug", "enhancement"]
created_at: "2024-01-01T00:00:00Z"
updated_at: "2024-01-02T00:00:00Z"
url: "https://github.com/owner/repo/issues/123"
repository: "owner/repo"
---
`

	got := generateFrontMatter(item)
	if got != want {
		t.Errorf("generateFrontMatter() =\n%q\nwant\n%q", got, want)
	}
}

// TestGenerateMarkdown 测试完整的 Markdown 生成
func TestGenerateMarkdown(t *testing.T) {
	tests := []struct {
		name    string
		item    *models.Item
		want    string
		wantErr bool
	}{
		{
			name: "simple issue without comments",
			item: &models.Item{
				Type:       models.ContentTypeIssue,
				Number:     1,
				Title:      "Test Issue",
				Body:       "This is a test issue",
				State:      models.StateOpen,
				Author:     models.User{Login: "testuser"},
				CreatedAt:  "2024-01-01T00:00:00Z",
				UpdatedAt:  "2024-01-01T00:00:00Z",
				URL:        "https://github.com/owner/repo/issues/1",
				Repository: "owner/repo",
				Comments:   []models.Comment{},
			},
			wantErr: false,
		},
		{
			name: "issue with comments",
			item: &models.Item{
				Type:       models.ContentTypeIssue,
				Number:     1,
				Title:      "Test Issue",
				Body:       "This is a test issue",
				State:      models.StateOpen,
				Author:     models.User{Login: "testuser"},
				CreatedAt:  "2024-01-01T00:00:00Z",
				UpdatedAt:  "2024-01-01T00:00:00Z",
				URL:        "https://github.com/owner/repo/issues/1",
				Repository: "owner/repo",
				Comments: []models.Comment{
					{
						ID:       1,
						Author:   models.User{Login: "commenter1"},
						Body:     "First comment",
						CreatedAt: "2024-01-01T01:00:00Z",
					},
				},
			},
			wantErr: false,
		},
		{
			name: "issue with nested comments",
			item: &models.Item{
				Type:       models.ContentTypeIssue,
				Number:     1,
				Title:      "Test Issue",
				Body:       "This is a test issue",
				State:      models.StateOpen,
				Author:     models.User{Login: "testuser"},
				CreatedAt:  "2024-01-01T00:00:00Z",
				UpdatedAt:  "2024-01-01T00:00:00Z",
				URL:        "https://github.com/owner/repo/issues/1",
				Repository: "owner/repo",
				Comments: []models.Comment{
					{
						ID:       1,
						Author:   models.User{Login: "commenter1"},
						Body:     "First comment",
						CreatedAt: "2024-01-01T01:00:00Z",
						Replies: []models.Comment{
							{
								ID:        2,
								Author:    models.User{Login: "commenter2"},
								Body:      "Reply to first",
								CreatedAt: "2024-01-01T02:00:00Z",
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			markdown, err := GenerateMarkdown(tt.item)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateMarkdown() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// 检查基本元素存在
				if !strings.Contains(markdown, tt.item.Title) {
					t.Errorf("Generated markdown missing title")
				}
				if !strings.Contains(markdown, tt.item.Body) {
					t.Errorf("Generated markdown missing body")
				}
				if !strings.Contains(markdown, "---") {
					t.Errorf("Generated markdown missing front matter delimiter")
				}
			}
		})
	}
}

// TestFormatComment 测试评论格式化
func TestFormatComment(t *testing.T) {
	tests := []struct {
		name       string
		comment    models.Comment
		indent     int
		wantHeader string
	}{
		{
			name: "simple comment",
			comment: models.Comment{
				Author:    models.User{Login: "testuser"},
				Body:      "Test comment",
				CreatedAt: "2024-01-01T00:00:00Z",
			},
			indent:     0,
			wantHeader: "### @testuser - 2024-01-01 00:00:00 UTC",
		},
		{
			name: "comment with special characters",
			comment: models.Comment{
				Author:    models.User{Login: "test-user"},
				Body:      "Comment with **bold** and *italic*",
				CreatedAt: "2024-01-01T00:00:00Z",
			},
			indent:     0,
			wantHeader: "### @test-user - 2024-01-01 00:00:00 UTC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatComment(tt.comment, tt.indent)
			if !strings.Contains(result, tt.wantHeader) {
				t.Errorf("formatComment() result =\n%s\ndoes not contain header %q", result, tt.wantHeader)
			}
			if !strings.Contains(result, tt.comment.Body) {
				t.Errorf("formatComment() result =\n%s\ndoes not contain body %q", result, tt.comment.Body)
			}
		})
	}
}

// TestConvertEmoji 测试 Emoji 转换
func TestConvertEmoji(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple emoji",
			input:    ":smile:",
			expected: "😊",
		},
		{
			name:     "thumbs up",
			input:    ":thumbsup:",
			expected: "👍",
		},
		{
			name:     "text without emoji",
			input:    "plain text",
			expected: "plain text",
		},
		{
			name:     "mixed content",
			input:    "Great work! :tada:",
			expected: "Great work! 🎉",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertEmoji(tt.input)
			if got != tt.expected {
				t.Errorf("convertEmoji(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestFormatTimestamp 测试时间戳格式化
func TestFormatTimestamp(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "ISO 8601 timestamp",
			input:    "2024-01-15T10:30:00Z",
			expected: "2024-01-15 10:30:00 UTC",
		},
		{
			name:     "with milliseconds",
			input:    "2024-01-15T10:30:00.123Z",
			expected: "2024-01-15 10:30:00 UTC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatTimestamp(tt.input)
			if got != tt.expected {
				t.Errorf("formatTimestamp(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestGenerateFileName 测试文件名生成
func TestGenerateFileName(t *testing.T) {
	tests := []struct {
		name     string
		item     *models.Item
		expected string
	}{
		{
			name: "issue",
			item: &models.Item{
				Type:   models.ContentTypeIssue,
				Number: 123,
			},
			expected: "issue-123.md",
		},
		{
			name: "pull request",
			item: &models.Item{
				Type:   models.ContentTypePullRequest,
				Number: 456,
			},
			expected: "pr-456.md",
		},
		{
			name: "discussion",
			item: &models.Item{
				Type:   models.ContentTypeDiscussion,
				Number: 789,
			},
			expected: "discussion-789.md",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateFileName(tt.item)
			if got != tt.expected {
				t.Errorf("GenerateFileName() = %q, want %q", got, tt.expected)
			}
		})
	}
}
