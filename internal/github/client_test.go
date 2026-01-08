package github

import (
	"strings"
	"testing"
)

// TestParseURL 测试 URL 解析功能
func TestParseURL(t *testing.T) {
	tests := []struct {
		name          string
		inputURL      string
		wantOwner     string
		wantRepo      string
		wantNumber    int
		wantType      string
		wantErr       bool
		errContains   string
	}{
		{
			name:     "valid issue URL",
			inputURL: "https://github.com/owner/repo/issues/123",
			wantOwner: "owner",
			wantRepo: "repo",
			wantNumber: 123,
			wantType: "issue",
			wantErr: false,
		},
		{
			name:     "valid PR URL",
			inputURL: "https://github.com/owner/repo/pull/456",
			wantOwner: "owner",
			wantRepo: "repo",
			wantNumber: 456,
			wantType: "pr",
			wantErr: false,
		},
		{
			name:     "valid discussion URL",
			inputURL: "https://github.com/owner/repo/discussions/789",
			wantOwner: "owner",
			wantRepo: "repo",
			wantNumber: 789,
			wantType: "discussion",
			wantErr: false,
		},
		{
			name:     "URL with www",
			inputURL: "https://www.github.com/owner/repo/issues/123",
			wantOwner: "owner",
			wantRepo: "repo",
			wantNumber: 123,
			wantType: "issue",
			wantErr: false,
		},
		{
			name:        "invalid URL - missing number",
			inputURL:    "https://github.com/owner/repo/issues",
			wantErr:     true,
			errContains: "invalid GitHub URL",
		},
		{
			name:        "invalid URL - wrong domain",
			inputURL:    "https://gitlab.com/owner/repo/issues/123",
			wantErr:     true,
			errContains: "invalid GitHub URL",
		},
		{
			name:        "invalid URL - malformed",
			inputURL:    "not-a-url",
			wantErr:     true,
			errContains: "invalid GitHub URL",
		},
		{
			name:        "unknown type",
			inputURL:    "https://github.com/owner/repo/unknown/123",
			wantErr:     true,
			errContains: "invalid GitHub URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info, err := ParseURL(tt.inputURL)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.errContains)
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("expected error containing %q, got %q", tt.errContains, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if info.Owner != tt.wantOwner || info.Repo != tt.wantRepo || info.Number != tt.wantNumber || info.ContentType != tt.wantType {
				t.Errorf("got (%q, %q, %d, %q), want (%q, %q, %d, %q)",
					info.Owner, info.Repo, info.Number, info.ContentType, tt.wantOwner, tt.wantRepo, tt.wantNumber, tt.wantType)
			}
		})
	}
}

// TestBuildAPIURL 测试 API URL 构建功能
func TestBuildAPIURL(t *testing.T) {
	tests := []struct {
		name       string
		owner      string
		repo       string
		contentType string
		number     int
		wantURL    string
	}{
		{
			name:       "issue API URL",
			owner:      "owner",
			repo:       "repo",
			contentType: "issue",
			number:     123,
			wantURL:    "https://api.github.com/repos/owner/repo/issues/123",
		},
		{
			name:       "PR API URL",
			owner:      "owner",
			repo:       "repo",
			contentType: "pr",
			number:     456,
			wantURL:    "https://api.github.com/repos/owner/repo/pulls/456",
		},
		{
			name:       "discussion API URL",
			owner:      "owner",
			repo:       "repo",
			contentType: "discussion",
			number:     789,
			wantURL:    "https://api.github.com/repos/owner/repo/discussions/789",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL := BuildAPIURL(tt.owner, tt.repo, tt.contentType, tt.number)
			if gotURL != tt.wantURL {
				t.Errorf("got %q, want %q", gotURL, tt.wantURL)
			}
		})
	}
}

// TestValidateURL 测试 URL 验证功能
func TestValidateURL(t *testing.T) {
	tests := []struct {
		name        string
		inputURL    string
		wantValid   bool
	}{
		{
			name:      "valid GitHub URL",
			inputURL:  "https://github.com/golang/go/issues/12345",
			wantValid: true,
		},
		{
			name:      "empty URL",
			inputURL:  "",
			wantValid: false,
		},
		{
			name:      "non-GitHub URL",
			inputURL:  "https://example.com",
			wantValid: false,
		},
		{
			name:      "malformed URL",
			inputURL:  "://github.com/owner/repo/issues/123",
			wantValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 使用我们的 ParseURL 来验证 URL 格式
			_, err := ParseURL(tt.inputURL)
			valid := err == nil

			if valid != tt.wantValid {
				t.Errorf("ParseURL(%q) valid = %v, want %v (err: %v)", tt.inputURL, valid, tt.wantValid, err)
			}
		})
	}
}

// TestNewClient 测试客户端创建
func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "client without token",
			token:   "",
			wantErr: false,
		},
		{
			name:    "client with token",
			token:   "ghp_test_token",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := NewClient(tt.token)
			if client == nil {
				t.Error("NewClient() returned nil")
			}
			if tt.token != "" && client.token != tt.token {
				t.Errorf("client.token = %q, want %q", client.token, tt.token)
			}
		})
	}
}
