// issue2md - 将 GitHub Issues/PRs/Discussions 转换为 Markdown
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bigwhite/issue2md/internal/converter"
	"github.com/bigwhite/issue2md/internal/github"
	"github.com/bigwhite/issue2md/internal/models"
)

const version = "0.1.0"

func main() {
	// 定义命令行参数
	var (
		outputPath string
		token      string
		showVersion bool
	)

	flag.StringVar(&outputPath, "o", "", "指定输出文件路径（可选）")
	flag.StringVar(&outputPath, "output", "", "指定输出文件路径（可选）")
	flag.StringVar(&token, "t", "", "GitHub Personal Access Token（可选）")
	flag.StringVar(&token, "token", "", "GitHub Personal Access Token（可选）")
	flag.BoolVar(&showVersion, "v", false, "显示版本信息")
	flag.BoolVar(&showVersion, "version", false, "显示版本信息")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "issue2md v%s - 将 GitHub Issues/PRs/Discussions 转换为 Markdown\n\n", version)
		fmt.Fprintf(os.Stderr, "用法:\n  issue2md [OPTIONS] <GITHUB_URL>\n\n")
		fmt.Fprintf(os.Stderr, "参数:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n示例:\n")
		fmt.Fprintf(os.Stderr, "  issue2md https://github.com/owner/repo/issues/123\n")
		fmt.Fprintf(os.Stderr, "  issue2md -o custom.md https://github.com/owner/repo/pull/456\n")
		fmt.Fprintf(os.Stderr, "  issue2md -t ghp_xxx https://github.com/owner/repo/discussions/789\n")
	}

	flag.Parse()

	// 显示版本
	if showVersion {
		fmt.Printf("issue2md version %s\n", version)
		os.Exit(0)
	}

	// 检查 URL 参数
	if flag.NArg() < 1 {
		flag.Usage()
		fmt.Fprintln(os.Stderr, "\n错误: 请提供 GitHub URL")
		os.Exit(1)
	}

	url := flag.Arg(0)

	// 执行转换
	if err := run(url, outputPath, token); err != nil {
		fmt.Fprintf(os.Stderr, "错误: %v\n", err)
		os.Exit(1)
	}
}

// run 执行主要的转换逻辑
func run(url, outputPath, token string) error {
	// 1. 解析 URL
	info, err := github.ParseURL(url)
	if err != nil {
		return fmt.Errorf("解析 URL 失败: %w", err)
	}

	// 2. 创建 GitHub 客户端
	client := github.NewClient(token)

	// 3. 获取数据并转换为内部模型
	var item *models.Item

	switch info.ContentType {
	case "issue":
		issueResp, err := client.FetchIssue(info.Owner, info.Repo, info.Number)
		if err != nil {
			return fmt.Errorf("%s", github.FormatErrorMessage(err, url))
		}
		comments, err := client.FetchIssueComments(info.Owner, info.Repo, info.Number)
		if err != nil {
			return fmt.Errorf("获取 Issue 评论失败: %w", err)
		}
		item = github.ToModel(issueResp, comments, url, info)

	case "pr":
		prResp, err := client.FetchPullRequest(info.Owner, info.Repo, info.Number)
		if err != nil {
			return fmt.Errorf("%s", github.FormatErrorMessage(err, url))
		}
		comments, reviews, err := client.FetchPRComments(info.Owner, info.Repo, info.Number)
		if err != nil {
			return fmt.Errorf("获取 PR 评论失败: %w", err)
		}
		var diff string
		if prResp.DiffURL != "" {
			diff, err = client.FetchDiff(prResp.DiffURL)
			if err != nil {
				return fmt.Errorf("获取 PR diff 失败: %w", err)
			}
		}
		item = github.ToModelPR(prResp, comments, reviews, diff, url, info)

	case "discussion":
		discussionResp, err := client.FetchDiscussion(info.Owner, info.Repo, info.Number)
		if err != nil {
			return fmt.Errorf("%s", github.FormatErrorMessage(err, url))
		}
		item = github.ToModelDiscussion(discussionResp, url, info)

	default:
		return fmt.Errorf("不支持的内容类型: %s", info.ContentType)
	}

	// 4. 转换为 Markdown
	markdown, err := converter.GenerateMarkdown(item)
	if err != nil {
		return fmt.Errorf("生成 Markdown 失败: %w", err)
	}

	// 5. 确定输出文件名
	if outputPath == "" {
		outputPath = converter.GenerateFileName(item)
	}

	// 6. 检查文件是否存在
	if _, err := os.Stat(outputPath); err == nil {
		fmt.Printf("文件 %s 已存在。是否覆盖? [y/N]: ", outputPath)
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			return fmt.Errorf("操作已取消")
		}
	}

	// 7. 写入文件
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}

	if err := os.WriteFile(outputPath, []byte(markdown), 0644); err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}

	fmt.Printf("✓ 已保存到: %s\n", outputPath)
	return nil
}
