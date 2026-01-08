# issue2md 功能规格说明书

## 项目概述

**issue2md** 是一个命令行工具，用于将 GitHub Issues、Pull Requests 和 Discussions 转换为 Markdown 文件，便于离线阅读和归档。

### 核心使命
提供一个简单、可靠的工具，让开发者能够将 GitHub 上的讨论内容保存为本地 Markdown 文档。

---

## 功能需求

### 1. 支持的内容类型

| 类型 | 支持状态 | 说明 |
|------|----------|------|
| GitHub Issues | ✅ | 公开仓库的 Issue |
| GitHub Pull Requests | ✅ | 公开仓库的 PR（含代码变更） |
| GitHub Discussions | ✅ | 公开仓库的 Discussion |

### 2. 输出格式规范

#### 2.1 文件结构
```markdown
---
title: "Issue/PR/Discussion 标题"
author: "作者名"
number: 123
type: "issue|pr|discussion"
state: "open|closed|merged"
labels: ["label1", "label2"]
created_at: "2024-01-01T00:00:00Z"
updated_at: "2024-01-02T00:00:00Z"
url: "https://github.com/owner/repo/issues/123"
repository: "owner/repo"
---

# [Issue/PR/Discussion] 标题

**作者:** @username
**创建时间:** 2024-01-01 00:00:00 UTC

## 正文

[原始正文内容，保留 Markdown 格式]

---

## 评论

### @username - 2024-01-01 12:00:00 UTC

[评论内容]

> ### @another_user - 2024-01-01 13:00:00 UTC
>
> [回复内容]

---

### PR 代码变更

[如果是 PR，在此嵌入完整的 diff]
```

#### 2.2 格式化规则

| 内容类型 | 处理方式 |
|----------|----------|
| 元数据 | YAML Front Matter |
| 代码块 | 保留语言标识符（如 ` ```python `） |
| Emoji | 转换为 Unicode 字符（如 `😊`） |
| 图片 | 保留原始 GitHub 链接 |
| 嵌套评论 | 使用 Markdown 引用语法（`>`）嵌套 |
| Review Comments | 与普通评论按时间混合展示 |
| PR Diff | 完整嵌入 |

#### 2.3 文件命名规则

```
issue-<number>.md      # Issue
pr-<number>.md         # Pull Request
discussion-<number>.md # Discussion
```

### 3. 命令行接口

#### 3.1 基本用法

```bash
issue2md <GITHUB_URL>
```

#### 3.2 参数选项

```bash
issue2md [OPTIONS] <URL>

OPTIONS:
  -o, --output string    指定输出文件路径（可选）
  -t, --token string     GitHub Personal Access Token（可选，用于提高 API 限制）
  -h, --help             显示帮助信息
  -v, --version          显示版本信息
```

#### 3.3 输入 URL 格式支持

```
https://github.com/owner/repo/issues/123
https://github.com/owner/repo/pull/456
https://github.com/owner/repo/discussions/789
```

### 4. 行为规范

| 场景 | 行为 |
|------|------|
| 文件已存在 | 交互式询问是否覆盖 |
| 部分内容获取失败 | 立即失败并显示错误 |
| 无 Token | 使用匿名访问（受 GitHub API 限制） |
| 私有仓库 | V1 不支持（明确提示用户） |
| 无效 URL | 显示错误并退出 |

---

## 技术约束

### 遵循项目宪法

1. **简单性原则**
   - 使用 Go 标准库 `net/http` 处理 HTTP 请求
   - 不引入非必需的第三方依赖
   - 优先使用简单的函数而非复杂抽象

2. **测试先行**
   - 所有功能从失败的测试开始
   - 优先使用表格驱动测试
   - 优先编写集成测试，避免 Mock

3. **明确性原则**
   - 所有错误必须显式处理
   - 使用 `fmt.Errorf("...: %w", err)` 包装错误
   - 无全局变量，依赖通过参数显式传递

### 技术栈

- **语言**: Go >= 1.24
- **HTTP 客户端**: `net/http` (标准库)
- **GitHub API**: REST API v3
- **命令行解析**: `flag` 包 (标准库)

---

## 非功能性需求

| 需求 | 说明 |
|------|------|
| 性能 | 单个 Issue/PR/Discussion 转换时间 < 10 秒 |
| 兼容性 | 支持 Windows、macOS、Linux |
| 可维护性 | 代码覆盖率 >= 80% |
| 用户体验 | 错误信息清晰，提供解决建议 |

---

## V1 范围界定

### 包含
- ✅ Issue、PR、Discussion 的基本转换
- ✅ 评论及嵌套回复
- ✅ PR 代码 diff
- ✅ 元数据提取
- ✅ YAML Front Matter

### 不包含（未来版本）
- ❌ 私有仓库支持
- ❌ 批量转换
- ❌ 图片本地化下载
- ❌ 自定义模板
- ❌ Markdown 渲染优化（如折叠长代码块）
- ❌ GitHub Actions 集成

---

## 验收标准

一个 Issue/PR/Discussion 被认为成功转换当：

1. ✅ 元数据完整且准确（标题、作者、时间、标签等）
2. ✅ 正文内容保留原始格式
3. ✅ 所有评论按正确顺序展示
4. ✅ 嵌套评论层级清晰
5. ✅ PR 包含完整代码 diff
6. ✅ 文件可以正常渲染为 Markdown
7. ✅ Emoji 正确显示

---

## 开发里程碑

1. **Phase 1**: 核心数据结构定义 + GitHub API 客户端
2. **Phase 2**: Issue 转换功能（最简单）
3. **Phase 3**: PR 转换功能（含 diff）
4. **Phase 4**: Discussion 转换功能（结构最复杂）
5. **Phase 5**: 命令行接口与文件输出
6. **Phase 6**: 错误处理与边界情况
