# API Sketch - issue2md Core Packages

**Version**: 1.0
**Status**: Design Draft
**Date**: 2026-01-08

---

## 1. internal/github 包

**职责**: 封装 GitHub API 交互逻辑，获取 Issue/PR/Discussion 数据。

### 1.1 数据类型

```go
// ResourceType 表示 GitHub 资源类型
type ResourceType int

const (
    ResourceIssue ResourceType = iota
    ResourcePullRequest
    ResourceDiscussion
)

// ParsedURL 表示解析后的 GitHub URL
type ParsedURL struct {
    Owner       string
    Repo        string
    Number      int
    Type        ResourceType
    OriginalURL string
}

// ReactionStats 表示反应统计
type ReactionStats struct {
    ThumbsUp    int `json:"thumbs_up"`
    ThumbsDown  int `json:"thumbs_down"`
    Laugh       int `json:"laugh"`
    Hooray      int `json:"hooray"`
    Confused    int `json:"confused"`
    Heart       int `json:"heart"`
    Rocket      int `json:"rocket"`
    Eyes        int `json:"eyes"`
}

// Comment 表示一条评论
type Comment struct {
    ID        string
    Author    string
    AuthorURL string
    Body      string
    CreatedAt time.Time
    UpdatedAt time.Time
    Reactions ReactionStats
    IsAnswer  bool // 仅 Discussion 使用
}

// GitHubData 表示统一的 GitHub 资源数据
type GitHubData struct {
    Title          string
    URL            string
    Author         string
    AuthorURL      string
    CreatedAt      time.Time
    UpdatedAt      time.Time
    Status         string // open/closed/merged
    Type           ResourceType
    Body           string
    Comments       []Comment
    Reactions      ReactionStats
    TotalComments  int
}
```

### 1.2 Client 结构体

```go
// Client 封装 GitHub API 客户端
type Client struct {
    token     string
    client    *http.Client
    baseURL   string
}

// NewClient 创建一个新的 GitHub API 客户端
// token 为可选的 GitHub Personal Access Token
func NewClient(token string) *Client

// FetchIssue 根据 owner/repo/number 获取 Issue 数据
// 返回统一的 GitHubData 结构或错误
func (c *Client) FetchIssue(ctx context.Context, owner, repo string, number int) (*GitHubData, error)

// FetchPullRequest 根据 owner/repo/number 获取 PR 数据
// 返回统一的 GitHubData 结构或错误
func (c *Client) FetchPullRequest(ctx context.Context, owner, repo string, number int) (*GitHubData, error)

// FetchDiscussion 根据 owner/repo/number 获取 Discussion 数据
// 返回统一的 GitHubData 结构或错误
func (c *Client) FetchDiscussion(ctx context.Context, owner, repo string, number int) (*GitHubData, error)

// Fetch 根据 ParsedURL 自动路由到对应的 Fetch 方法
func (c *Client) Fetch(ctx context.Context, url ParsedURL) (*GitHubData, error)
```

### 1.3 错误处理

```go
// 定义包级别的错误类型
var (
    ErrInvalidURL      = errors.New("invalid GitHub URL")
    ErrResourceNotFound = errors.New("resource not found")
    ErrRateLimitExceeded = errors.New("rate limit exceeded")
    ErrUnauthorized    = errors.New("unauthorized: invalid token")
)
```

---

## 2. internal/parser 包

**职责**: 解析和验证 GitHub URL，识别资源类型。

### 2.1 主要接口

```go
// ParseURL 解析 GitHub URL 并返回 ParsedURL 结构
// 支持的格式:
//   - https://github.com/{owner}/{repo}/issues/{number}
//   - https://github.com/{owner}/{repo}/pull/{number}
//   - https://github.com/{owner}/{repo}/discussions/{number}
func ParseURL(rawURL string) (*ParsedURL, error)

// Validate 验证 ParsedURL 的有效性
func (u *ParsedURL) Validate() error

// String 返回原始 URL 字符串
func (u *ParsedURL) String() string
```

---

## 3. internal/converter 包

**职责**: 将 GitHub 数据转换为 Markdown 格式。

### 3.1 配置选项

```go
// Options 表示转换器的配置选项
type Options struct {
    EnableReactions    bool  // 包含 reactions 统计
    EnableUserLinks    bool  // 将 @username 转换为链接
    OutputPath         string // 输出文件路径（空表示输出到 stdout）
}
```

### 3.2 主要接口

```go
// Converter 负责将 GitHub 数据转换为 Markdown
type Converter struct {
    opts Options
}

// NewConverter 创建一个新的转换器实例
func NewConverter(opts Options) *Converter

// Convert 将 GitHubData 转换为 Markdown 字符串
// 生成的 Markdown 包含:
//   - YAML frontmatter
//   - 标题和元数据
//   - 描述内容
//   - 所有评论
func (c *Converter) Convert(data *github.GitHubData) (string, error)

// ConvertToFile 转换并直接写入文件
func (c *Converter) ConvertToFile(data *github.GitHubData, filepath string) error

// formatFrontmatter 生成 YAML frontmatter
func (c *Converter) formatFrontmatter(data *github.GitHubData) string

// formatBody 生成 Markdown 正文
func (c *Converter) formatBody(data *github.GitHubData) string

// formatComments 格式化评论列表
func (c *Converter) formatComments(data *github.GitHubData) string
```

### 3.3 辅助函数

```go
// formatTimestamp 格式化时间戳为 UTC 字符串
func formatTimestamp(t time.Time) string

// formatReactions 格式化 reactions 统计
func formatReactions(r github.ReactionStats) string

// linkifyUsernames 将 @username 转换为 Markdown 链接
func linkifyUsernames(text string) string
```

---

## 4. internal/config 包

**职责**: 管理配置和选项。

### 4.1 配置结构

```go
// Config 表示应用配置
type Config struct {
    GitHubToken    string
    EnableReactions bool
    EnableUserLinks bool
    OutputFile     string
    URL           string
}

// LoadFromEnv 从环境变量加载配置
func LoadFromEnv() *Config

// LoadFromArgs 从命令行参数加载配置
func LoadFromArgs(args []string) (*Config, error)
```

---

## 5. internal/cli 包

**职责**: 命令行接口和用户交互。

### 5.1 主要接口

```go
// CLI 表示命令行应用
type CLI struct {
    config *config.Config
    github *github.Client
    converter *converter.Converter
    stdout io.Writer
    stderr io.Writer
}

// NewCLI 创建一个新的 CLI 实例
func NewCLI(cfg *config.Config) *CLI

// Run 执行 CLI 应用
// 返回退出码（0 表示成功，非 0 表示错误）
func (cli *CLI) Run(ctx context.Context) int

// printHelp 显示帮助信息
func (cli *CLI) printHelp()

// printVersion 显示版本信息
func (cli *CLI) printVersion()
```

### 5.2 命令行参数

```bash
issue2md [flags] <url> [output_file]

Flags:
  -enable-reactions    包含 reactions 统计信息 (默认: false)
  -enable-user-links   将用户名渲染为 GitHub 主页链接 (默认: false)
  -h, -help           显示帮助信息
  -v, -version        显示版本信息

Environment Variables:
  GITHUB_TOKEN        GitHub Personal Access Token (可选)
```

---

## 6. 包依赖关系

```
cmd/issue2md/
    └── internal/cli/
        ├── internal/config/
        ├── internal/parser/
        ├── internal/github/
        └── internal/converter/
            └── internal/github/
```

**依赖规则**:
- `github` 包：无内部依赖，仅使用标准库和 net/http
- `parser` 包：无内部依赖
- `converter` 包：依赖 `github` 包（使用 GitHubData 类型）
- `config` 包：无内部依赖
- `cli` 包：依赖所有其他包

---

## 7. 实现优先级

### Phase 1: Core Foundation
1. `internal/github` - Client 基础实现 + FetchIssue
2. `internal/parser` - URL 解析
3. `internal/converter` - 基础 Markdown 转换

### Phase 2: CLI Integration
4. `internal/config` - 配置管理
5. `internal/cli` - 命令行接口
6. `cmd/issue2md` - 主入口

### Phase 3: Extended Features
7. `internal/github` - FetchPullRequest + FetchDiscussion
8. `internal/converter` - Reactions + UserLinks

---

## 8. 测试策略

### 单元测试（表格驱动）
- `internal/parser`: 测试各种 URL 格式解析
- `internal/github`: 测试 API 响应解析
- `internal/converter`: 测试 Markdown 生成

### 集成测试
- 端到端测试：URL → Parser → GitHub → Converter → Markdown
- 使用真实的 GitHub API（或 mock server）

---

**下一步**: 根据 API Sketch 开始实现 Phase 1 的核心功能。
