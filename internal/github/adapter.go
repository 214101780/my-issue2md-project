package github

import (
	"fmt"
	"strings"

	"github.com/bigwhite/issue2md/internal/models"
)

// ToModel 将 IssueResponse 转换为内部模型
func ToModel(issue *IssueResponse, comments []CommentResponse, url string, info *ParseURLInfo) *models.Item {
	item := &models.Item{
		Type:        models.ContentTypeIssue,
		Number:      issue.Number,
		Title:       issue.Title,
		Body:        issue.Body,
		State:       models.ItemState(issue.State),
		Author:      toModelUser(issue.User),
		CreatedAt:   issue.CreatedAt,
		UpdatedAt:   issue.UpdatedAt,
		URL:         url,
		Repository:  fmt.Sprintf("%s/%s", info.Owner, info.Repo),
		IsPR:        issue.PullRequest != nil,
	}

	// 处理 ClosedAt
	if issue.ClosedAt != nil {
		item.ClosedAt = issue.ClosedAt
	}

	// 处理标签
	if len(issue.Labels) > 0 {
		item.Labels = make([]models.Label, len(issue.Labels))
		for i, label := range issue.Labels {
			item.Labels[i] = models.Label{
				Name:        label.Name,
				Color:       label.Color,
				Description: label.Description,
			}
		}
	}

	// 处理 Milestone
	if issue.Milestone != nil {
		milestone := issue.Milestone.Title
		item.Milestone = &milestone
	}

	// 处理 Assignees
	if len(issue.Assignees) > 0 {
		item.Assignees = make([]models.User, len(issue.Assignees))
		for i, assignee := range issue.Assignees {
			item.Assignees[i] = toModelUser(assignee)
		}
	}

	// 处理评论
	item.Comments = convertComments(comments)

	return item
}

// ToModelPR 将 PRResponse 转换为内部模型
func ToModelPR(pr *PullRequestResponse, comments []CommentResponse, reviews []ReviewCommentResponse, diff string, url string, info *ParseURLInfo) *models.Item {
	item := &models.Item{
		Type:        models.ContentTypePullRequest,
		Number:      pr.Number,
		Title:       pr.Title,
		Body:        pr.Body,
		State:       models.ItemState(pr.State),
		Author:      toModelUser(pr.User),
		CreatedAt:   pr.CreatedAt,
		UpdatedAt:   pr.UpdatedAt,
		URL:         url,
		Repository:  fmt.Sprintf("%s/%s", info.Owner, info.Repo),
		IsPR:        true,
		Additions:   pr.Additions,
		Deletions:   pr.Deletions,
		ChangedFiles: pr.ChangedFiles,
		Diff:        diff,
	}

	// 处理 ClosedAt
	if pr.ClosedAt != nil {
		item.ClosedAt = pr.ClosedAt
	}

	// 处理 MergedAt
	if pr.MergedAt != nil {
		item.MergedAt = pr.MergedAt
		item.State = models.StateMerged
	}

	// 处理标签
	if len(pr.Labels) > 0 {
		item.Labels = make([]models.Label, len(pr.Labels))
		for i, label := range pr.Labels {
			item.Labels[i] = models.Label{
				Name:        label.Name,
				Color:       label.Color,
				Description: label.Description,
			}
		}
	}

	// 处理 Milestone
	if pr.Milestone != nil {
		milestone := pr.Milestone.Title
		item.Milestone = &milestone
	}

	// 处理 Assignees
	if len(pr.Assignees) > 0 {
		item.Assignees = make([]models.User, len(pr.Assignees))
		for i, assignee := range pr.Assignees {
			item.Assignees[i] = toModelUser(assignee)
		}
	}

	// 合并普通评论和 Review Comments
	allComments := convertComments(comments)
	for _, review := range reviews {
		allComments = append(allComments, models.Comment{
			ID:        review.ID,
			Author:    toModelUser(review.User),
			Body:      review.Body,
			CreatedAt: review.CreatedAt,
			UpdatedAt: review.UpdatedAt,
		})
	}

	item.Comments = allComments

	return item
}

// ToModelDiscussion 将 DiscussionResponse 转换为内部模型
func ToModelDiscussion(discussion *DiscussionResponse, url string, info *ParseURLInfo) *models.Item {
	item := &models.Item{
		Type:       models.ContentTypeDiscussion,
		Number:     discussion.Number,
		Title:      discussion.Title,
		Body:       discussion.Body,
		Author:     models.User{Login: discussion.Author.Login},
		CreatedAt:  discussion.CreatedAt,
		UpdatedAt:  discussion.UpdatedAt,
		URL:        url,
		Repository: fmt.Sprintf("%s/%s", info.Owner, info.Repo),
	}

	// 处理 ClosedAt
	if discussion.ClosedAt != nil {
		item.ClosedAt = discussion.ClosedAt
		item.State = models.StateClosed
	} else {
		item.State = models.StateOpen
	}

	// 处理 Category
	if discussion.Category.Slug != "" {
		item.Category = &discussion.Category.Slug
	}

	// 处理标签
	if len(discussion.Labels) > 0 {
		item.Labels = make([]models.Label, len(discussion.Labels))
		for i, label := range discussion.Labels {
			item.Labels[i] = models.Label{
				Name:        label.Name,
				Color:       label.Color,
				Description: label.Description,
			}
		}
	}

	// Discussion 评论的处理比较复杂，因为它们是嵌套的
	// 这里简化处理，V1 版本可以只获取顶级评论
	// TODO: 后续版本完善 Discussion 评论的嵌套处理

	return item
}

// toModelUser 将 GitHub User 转换为内部模型
func toModelUser(user User) models.User {
	return models.User{
		Login:     user.Login,
		ID:        user.ID,
		AvatarURL: user.AvatarURL,
		Type:      user.Type,
	}
}

// convertComments 将 CommentResponse 数组转换为内部模型
func convertComments(comments []CommentResponse) []models.Comment {
	if len(comments) == 0 {
		return []models.Comment{}
	}

	result := make([]models.Comment, len(comments))
	for i, comment := range comments {
		result[i] = models.Comment{
			ID:        comment.ID,
			Author:    toModelUser(comment.User),
			Body:      comment.Body,
			CreatedAt: comment.CreatedAt,
			UpdatedAt: comment.UpdatedAt,
			Replies:   []models.Comment{}, // 简化处理，暂不处理嵌套
		}
	}
	return result
}

// IsPrivateRepoURL 检查是否为私有仓库 URL
func IsPrivateRepoURL(url string) bool {
	// V1 不支持私有仓库，此函数预留
	// 可以通过 API 调用失败来判断
	return false
}

// FormatErrorMessage 格式化错误消息
func FormatErrorMessage(err error, url string) string {
	if strings.Contains(err.Error(), "404") {
		return fmt.Sprintf("未找到内容 (%s)\n可能原因:\n- URL 不正确\n- 内容是私有的（V1 不支持私有仓库）\n- 内容已被删除", url)
	}
	if strings.Contains(err.Error(), "403") {
		return fmt.Sprintf("访问被拒绝 (%s)\n可能原因:\n- API 速率限制已超出，请提供 GitHub Token\n- 内容是私有的（V1 不支持私有仓库）", url)
	}
	return err.Error()
}
