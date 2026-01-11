// Package parser 提供 GitHub URL 解析功能
package parser

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/bigwhite/issue2md/internal/models"
)

// 定义解析错误
var (
	ErrInvalidURL     = errors.New("invalid GitHub URL")
	ErrUnsupportedURL = errors.New("unsupported URL type")
)

// ParseResult 表示解析后的 GitHub 资源 URL
type ParseResult struct {
	Owner  string              // 仓库所有者
	Repo   string              // 仓库名称
	Number int                 // Issue/PR/Discussion 编号
	Type   models.ResourceType // 资源类型
}

// pathParts 存储 URL 分割后的各部分
type pathParts struct {
	domain       string
	owner        string
	repo         string
	resourceType string
	numberStr    string
}

// ParseURL 解析 GitHub URL 并返回资源信息
func ParseURL(urlStr string) (*ParseResult, error) {
	parts, err := splitAndValidatePath(urlStr)
	if err != nil {
		return nil, err
	}

	if err := validateOwnerRepo(parts.owner, parts.repo); err != nil {
		return nil, err
	}

	resType, err := parseResourceType(parts.resourceType)
	if err != nil {
		return nil, err
	}

	number, err := parseNumber(parts.numberStr)
	if err != nil {
		return nil, err
	}

	return &ParseResult{
		Owner:  parts.owner,
		Repo:   parts.repo,
		Number: number,
		Type:   resType,
	}, nil
}

// splitAndValidatePath 分割 URL 路径并进行基础验证
func splitAndValidatePath(urlStr string) (*pathParts, error) {
	if urlStr == "" {
		return nil, fmt.Errorf("empty URL: %w", ErrInvalidURL)
	}

	// 检查协议前缀
	if !strings.HasPrefix(urlStr, "https://") && !strings.HasPrefix(urlStr, "http://") {
		return nil, fmt.Errorf("missing protocol: %w", ErrInvalidURL)
	}

	// 移除协议前缀
	afterProto := strings.TrimPrefix(urlStr, "https://")
	afterProto = strings.TrimPrefix(afterProto, "http://")

	// 分割路径部分
	parts := strings.Split(afterProto, "/")
	if len(parts) < 5 {
		// 如果只有 3 个部分（domain/owner/repo），可能是仓库主页
		if len(parts) == 3 {
			// 检查第二部分是否是资源类型关键字，如果不是则是仓库主页（不支持）
			if parts[1] != "issues" && parts[1] != "pull" && parts[1] != "discussions" {
				return nil, fmt.Errorf("repository homepage not supported: %w", ErrUnsupportedURL)
			}
		}
		return nil, fmt.Errorf("invalid URL structure: %w", ErrInvalidURL)
	}

	// 检查是否有额外的路径段（如尾部斜杠、查询参数等）
	if len(parts) > 5 {
		return nil, fmt.Errorf("extra path segments: %w", ErrInvalidURL)
	}

	// 验证域名
	if parts[0] != "github.com" {
		return nil, fmt.Errorf("invalid domain %q: %w", parts[0], ErrInvalidURL)
	}

	return &pathParts{
		domain:       parts[0],
		owner:        parts[1],
		repo:         parts[2],
		resourceType: parts[3],
		numberStr:    parts[4],
	}, nil
}

// validateOwnerRepo 验证 owner 和 repo 名称格式
func validateOwnerRepo(owner, repo string) error {
	if owner == "" || !isValidOwnerRepo(owner) {
		return fmt.Errorf("invalid owner %q: %w", owner, ErrInvalidURL)
	}
	if repo == "" || !isValidOwnerRepo(repo) {
		return fmt.Errorf("invalid repo %q: %w", repo, ErrInvalidURL)
	}
	return nil
}

// parseResourceType 解析资源类型
func parseResourceType(resourceType string) (models.ResourceType, error) {
	switch resourceType {
	case "issues":
		return models.ResourceTypeIssue, nil
	case "pull":
		return models.ResourceTypePullRequest, nil
	case "discussions":
		return models.ResourceTypeDiscussion, nil
	default:
		return 0, fmt.Errorf("unknown resource type %q: %w", resourceType, ErrUnsupportedURL)
	}
}

// parseNumber 解析编号
func parseNumber(numberStr string) (int, error) {
	number, err := strconv.Atoi(numberStr)
	if err != nil {
		return 0, fmt.Errorf("invalid number %q: %w", numberStr, ErrInvalidURL)
	}
	if number <= 0 {
		return 0, fmt.Errorf("number must be positive, got %d: %w", number, ErrInvalidURL)
	}
	return number, nil
}

// isValidOwnerRepo 验证 owner/repo 名称格式
// 只允许字母、数字和连字符
func isValidOwnerRepo(name string) bool {
	if name == "" {
		return false
	}
	for _, r := range name {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-') {
			return false
		}
	}
	return true
}
