// Package parser 的表格驱动测试
package parser

import (
	"errors"
	"testing"

	"github.com/bigwhite/issue2md/internal/models"
)

// TestParseURL 使用表格驱动测试验证 ParseURL 函数
func TestParseURL(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantResult  *ParseResult
		wantErr     error
		description string
	}{
		{
			name:        "valid_issue_url",
			input:       "https://github.com/bigwhite/issue2md/issues/1",
			wantResult: &ParseResult{
				Owner:  "bigwhite",
				Repo:   "issue2md",
				Number: 1,
				Type:   models.ResourceTypeIssue,
			},
			wantErr:     nil,
			description: "合法的 Issue URL",
		},
		{
			name:        "valid_pull_request_url",
			input:       "https://github.com/golang/go/pull/42",
			wantResult: &ParseResult{
				Owner:  "golang",
				Repo:   "go",
				Number: 42,
				Type:   models.ResourceTypePullRequest,
			},
			wantErr:     nil,
			description: "合法的 Pull Request URL",
		},
		{
			name:        "valid_discussion_url",
			input:       "https://github.com/bigwhite/issue2md/discussions/123",
			wantResult: &ParseResult{
				Owner:  "bigwhite",
				Repo:   "issue2md",
				Number: 123,
				Type:   models.ResourceTypeDiscussion,
			},
			wantErr:     nil,
			description: "合法的 Discussion URL",
		},
		{
			name:        "invalid_url_format",
			input:       "invalid-url",
			wantResult:  nil,
			wantErr:     ErrInvalidURL,
			description: "无效的 URL（格式错误）",
		},
		{
			name:        "unsupported_url_type",
			input:       "https://github.com/bigwhite/issue2md",
			wantResult:  nil,
			wantErr:     ErrUnsupportedURL,
			description: "不支持的 URL 类型（仓库主页）",
		},
		{
			name:        "url_with_trailing_slash",
			input:       "https://github.com/bigwhite/issue2md/issues/1/",
			wantResult:  nil,
			wantErr:     ErrInvalidURL,
			description: "带尾部斜杠的 URL",
		},
		{
			name:        "missing_owner",
			input:       "https://github.com/issues/1",
			wantResult:  nil,
			wantErr:     ErrInvalidURL,
			description: "缺少所有者的 URL",
		},
		{
			name:        "missing_repo",
			input:       "https://github.com/bigwhite/issues/1",
			wantResult:  nil,
			wantErr:     ErrInvalidURL,
			description: "缺少仓库名的 URL",
		},
		{
			name:        "missing_number",
			input:       "https://github.com/bigwhite/issue2md/issues",
			wantResult:  nil,
			wantErr:     ErrInvalidURL,
			description: "缺少编号的 URL",
		},
		{
			name:        "invalid_number",
			input:       "https://github.com/bigwhite/issue2md/issues/abc",
			wantResult:  nil,
			wantErr:     ErrInvalidURL,
			description: "编号为非数字的 URL",
		},
		{
			name:        "empty_string",
			input:       "",
			wantResult:  nil,
			wantErr:     ErrInvalidURL,
			description: "空字符串",
		},
		{
			name:        "not_github_domain",
			input:       "https://gitlab.com/bigwhite/issue2md/issues/1",
			wantResult:  nil,
			wantErr:     ErrInvalidURL,
			description: "非 GitHub 域名",
		},
		{
			name:        "missing_protocol",
			input:       "github.com/bigwhite/issue2md/issues/1",
			wantResult:  nil,
			wantErr:     ErrInvalidURL,
			description: "缺少协议的 URL",
		},
		{
			name:        "special_characters_in_owner",
			input:       "https://github.com/big_white/issue2md/issues/1",
			wantResult:  nil,
			wantErr:     ErrInvalidURL,
			description: "所有者包含特殊字符",
		},
		{
			name:        "special_characters_in_repo",
			input:       "https://github.com/bigwhite/issue_2md/issues/1",
			wantResult:  nil,
			wantErr:     ErrInvalidURL,
			description: "仓库名包含特殊字符",
		},
		{
			name:        "hyphen_allowed_in_owner",
			input:       "https://github.com/big-white/issue2md/issues/1",
			wantResult: &ParseResult{
				Owner:  "big-white",
				Repo:   "issue2md",
				Number: 1,
				Type:   models.ResourceTypeIssue,
			},
			wantErr:     nil,
			description: "所有者包含连字符（合法）",
		},
		{
			name:        "hyphen_allowed_in_repo",
			input:       "https://github.com/bigwhite/issue-2md/issues/1",
			wantResult: &ParseResult{
				Owner:  "bigwhite",
				Repo:   "issue-2md",
				Number: 1,
				Type:   models.ResourceTypeIssue,
			},
			wantErr:     nil,
			description: "仓库名包含连字符（合法）",
		},
		{
			name:        "number_with_leading_zero",
			input:       "https://github.com/bigwhite/issue2md/issues/01",
			wantResult: &ParseResult{
				Owner:  "bigwhite",
				Repo:   "issue2md",
				Number: 1,
				Type:   models.ResourceTypeIssue,
			},
			wantErr:     nil,
			description: "编号前导零（应解析为 1）",
		},
		{
			name:        "zero_number",
			input:       "https://github.com/bigwhite/issue2md/issues/0",
			wantResult:  nil,
			wantErr:     ErrInvalidURL,
			description: "编号为零（无效）",
		},
		{
			name:        "negative_number",
			input:       "https://github.com/bigwhite/issue2md/issues/-1",
			wantResult:  nil,
			wantErr:     ErrInvalidURL,
			description: "负数编号",
		},
		{
			name:        "url_with_query_params",
			input:       "https://github.com/bigwhite/issue2md/issues/1?q=1",
			wantResult:  nil,
			wantErr:     ErrInvalidURL,
			description: "带查询参数的 URL",
		},
		{
			name:        "url_with_fragment",
			input:       "https://github.com/bigwhite/issue2md/issues/1#section",
			wantResult:  nil,
			wantErr:     ErrInvalidURL,
			description: "带片段的 URL",
		},
		{
			name:        "case_insensitive_issues_path",
			input:       "https://github.com/bigwhite/issue2md/ISSUES/1",
			wantResult:  nil,
			wantErr:     ErrUnsupportedURL,
			description: "大写 ISSUES 路径（应支持或明确拒绝）",
		},
		{
			name:        "uppercase_host",
			input:       "https://GITHUB.COM/bigwhite/issue2md/issues/1",
			wantResult:  nil,
			wantErr:     ErrInvalidURL,
			description: "大写域名（应标准化为小写处理）",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotResult, gotErr := ParseURL(tt.input)

			// 验证错误
			if tt.wantErr != nil {
				if gotErr == nil {
					t.Errorf("[%s] ParseURL(%q) expected error %v, but got nil", tt.description, tt.input, tt.wantErr)
				} else if !errors.Is(gotErr, tt.wantErr) {
					t.Errorf("[%s] ParseURL(%q) expected error %v, but got %v", tt.description, tt.input, tt.wantErr, gotErr)
				}
				return
			}

			// 无错误情况，验证结果
			if gotErr != nil {
				t.Errorf("[%s] ParseURL(%q) unexpected error: %v", tt.description, tt.input, gotErr)
				return
			}

			if gotResult == nil {
				t.Errorf("[%s] ParseURL(%q) expected result %v, but got nil", tt.description, tt.input, tt.wantResult)
				return
			}

			// 验证各个字段
			if gotResult.Owner != tt.wantResult.Owner {
				t.Errorf("[%s] ParseURL(%q).Owner = %q, want %q", tt.description, tt.input, gotResult.Owner, tt.wantResult.Owner)
			}
			if gotResult.Repo != tt.wantResult.Repo {
				t.Errorf("[%s] ParseURL(%q).Repo = %q, want %q", tt.description, tt.input, gotResult.Repo, tt.wantResult.Repo)
			}
			if gotResult.Number != tt.wantResult.Number {
				t.Errorf("[%s] ParseURL(%q).Number = %d, want %d", tt.description, tt.input, gotResult.Number, tt.wantResult.Number)
			}
			if gotResult.Type != tt.wantResult.Type {
				t.Errorf("[%s] ParseURL(%q).Type = %v, want %v", tt.description, tt.input, gotResult.Type, tt.wantResult.Type)
			}
		})
	}
}
