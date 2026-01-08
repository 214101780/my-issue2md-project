// Package converter 提供 GitHub 数据到 Markdown 的转换功能
package converter

import (
	"fmt"
	"strings"
	"time"

	"github.com/bigwhite/issue2md/internal/models"
)

// GenerateFileName 根据内容类型和编号生成文件名
func GenerateFileName(item *models.Item) string {
	var prefix string
	switch item.Type {
	case models.ContentTypeIssue:
		prefix = "issue"
	case models.ContentTypePullRequest:
		prefix = "pr"
	case models.ContentTypeDiscussion:
		prefix = "discussion"
	default:
		prefix = "item"
	}
	return fmt.Sprintf("%s-%d.md", prefix, item.Number)
}

// GenerateMarkdown 将 Item 转换为完整的 Markdown 文档
func GenerateMarkdown(item *models.Item) (string, error) {
	var sb strings.Builder

	// 1. Front Matter
	sb.WriteString(generateFrontMatter(item))
	sb.WriteString("\n")

	// 2. 标题
	sb.WriteString(fmt.Sprintf("# %s\n\n", item.Title))

	// 3. 元数据信息
	sb.WriteString(fmt.Sprintf("**作者:** @%s\n", item.Author.Login))
	sb.WriteString(fmt.Sprintf("**创建时间:** %s\n\n", formatTimestamp(item.CreatedAt)))

	// 4. 正文
	if item.Body != "" {
		sb.WriteString("## 正文\n\n")
		sb.WriteString(item.Body)
		sb.WriteString("\n\n")
	}

	// 5. PR 特有内容：代码变更
	if item.IsPR && item.Diff != "" {
		sb.WriteString("---\n\n")
		sb.WriteString("## 代码变更\n\n")
		sb.WriteString("```diff\n")
		sb.WriteString(item.Diff)
		sb.WriteString("\n```\n\n")
	}

	// 6. 评论
	if len(item.Comments) > 0 {
		sb.WriteString("---\n\n")
		sb.WriteString("## 评论\n\n")
		for _, comment := range item.Comments {
			sb.WriteString(formatComment(comment, 0))
			sb.WriteString("\n")
		}
	}

	return sb.String(), nil
}

// generateFrontMatter 生成 YAML Front Matter
func generateFrontMatter(item *models.Item) string {
	var sb strings.Builder

	sb.WriteString("---\n")
	sb.WriteString(fmt.Sprintf("title: \"%s\"\n", escapeYAMLString(item.Title)))
	sb.WriteString(fmt.Sprintf("author: \"%s\"\n", item.Author.Login))
	sb.WriteString(fmt.Sprintf("number: %d\n", item.Number))
	sb.WriteString(fmt.Sprintf("type: \"%s\"\n", item.Type))
	sb.WriteString(fmt.Sprintf("state: \"%s\"\n", item.State))

	// 标签列表
	if len(item.Labels) > 0 {
		sb.WriteString("labels: [")
		for i, label := range item.Labels {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("\"%s\"", label.Name))
		}
		sb.WriteString("]\n")
	}

	sb.WriteString(fmt.Sprintf("created_at: \"%s\"\n", item.CreatedAt))
	sb.WriteString(fmt.Sprintf("updated_at: \"%s\"\n", item.UpdatedAt))
	sb.WriteString(fmt.Sprintf("url: \"%s\"\n", item.URL))
	sb.WriteString(fmt.Sprintf("repository: \"%s\"\n", item.Repository))
	sb.WriteString("---\n")

	return sb.String()
}

// formatComment 格式化单条评论，支持嵌套
func formatComment(comment models.Comment, indent int) string {
	var sb strings.Builder
	indentStr := strings.Repeat("  ", indent)

	// 标题
	sb.WriteString(fmt.Sprintf("%s### @%s - %s\n\n", indentStr, comment.Author.Login, formatTimestamp(comment.CreatedAt)))

	// 内容（转换 emoji）
	body := convertEmoji(comment.Body)
	if body != "" {
		// 按行处理，添加引用标记
		lines := strings.Split(body, "\n")
		for _, line := range lines {
			sb.WriteString(fmt.Sprintf("%s%s\n", indentStr, line))
		}
		sb.WriteString("\n")
	}

	// 递归处理回复
	if len(comment.Replies) > 0 {
		for _, reply := range comment.Replies {
			// 嵌套回复使用引用格式
			replyContent := formatComment(reply, indent+1)
			// 为嵌套内容添加引用标记
			replyLines := strings.Split(replyContent, "\n")
			for _, line := range replyLines {
				if line != "" {
					sb.WriteString(fmt.Sprintf("%s> %s\n", indentStr, line))
				} else {
					sb.WriteString(indentStr + ">\n")
				}
			}
			sb.WriteString("\n")
		}
	}

	return sb.String()
}

// convertEmoji 将 GitHub emoji 代码转换为 Unicode 字符
func convertEmoji(text string) string {
	// 常用的 emoji 映射
	emojiMap := map[string]string{
		// 表情
		":smile:":     "😊",
		":laughing:":  "😆",
		":blush:":     "😊",
		":heart_eyes:": "😍",
		":cry:":       "😢",
		":sob:":       "😭",
		":angry:":     "😠",
		":thumbsup:":  "👍",
		":thumbsdown:": "👎",
		":clap:":      "👏",
		":wave:":      "👋",
		":heart:":     "❤️",
		":fire:":      "🔥",
		":star:":      "⭐",
		":tada:":      "🎉",
		":rocket:":    "🚀",
		":check:":     "✅",
		":x:":         "❌",
		":warning:":   "⚠️",
		":question:":  "❓",
		":+1:":        "👍",
		":-1:":        "👎",
	}

	result := text
	for code, unicode := range emojiMap {
		result = strings.ReplaceAll(result, code, unicode)
	}
	return result
}

// formatTimestamp 将 ISO 8601 时间戳格式化为可读格式
func formatTimestamp(isoTime string) string {
	// 解析 ISO 8601 时间
	t, err := time.Parse(time.RFC3339, isoTime)
	if err != nil {
		// 如果解析失败，返回原始字符串
		return isoTime
	}

	// 格式化为 "YYYY-MM-DD HH:MM:SS UTC"
	return t.Format("2006-01-02 15:04:05 MST")
}

// escapeYAMLString 转义 YAML 字符串中的特殊字符
func escapeYAMLString(s string) string {
	// 转义反斜杠和引号
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}
