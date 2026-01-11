---
description: 根据当前暂存区的代码变更，生成一条符合Conventional Commits规范的Commit Message。
allowed-tools: Bash(git add:*)
---

你是一位Git专家。请根据以下提供的 **Git上下文信息**，按照逻辑判断并生成内容，为我生成一条符合Conventional Commits规范的、高质量的`git commit`中文消息。

### Git上下文信息

**1.当前分支 (Current Branch):**
!`git branch --show-current`

**2.文件状态概览(File Status):**
!`git status --short`

**3.暂存区变更 (Staged Changes - 优先):**
!`git diff --staged`

**4.工作区未暂存变更 (Unstaged Changes - 备用)**
!`git diff`

### 生成逻辑与要求
请按照以下优先级顺序进行判断和输出：
**情况A：如果「暂存区变更」不为空**
1. 忽略「工作区未暂存变更」。
2. 仅根据「暂存区变更」的内容，生成一条符合Conventional Commits 规范的 Commit Message。
3. **输出格式**：直接输出 Commit Message，不需要任何前缀或解释。

**情况 B：如果「暂存区变更」为空，但「工作区未暂存变更」或「文件状态概览」不为空**
1.  这说明用户忘记执行 `git add` 了。
2.  根据「工作区未暂存变更」的内容，生成一条 **建议的** Commit Message。
3.  **重要**：在输出 Commit Message 之前，必须先用一行加粗文字提示用户：**"暂存区为空，以下是根据工作区变更生成的建议（请先 git add）："**

**情况 C：如果所有变更都为空**
1.  直接回复："当前工作区是干净的，没有检测到任何代码变更。"

### Commit Message 规范
- 格式：`<type>(<scope>): <subject>`
- 常用Type：feat, fix, docs, style, refactor, perf, test, build, ci, chore, revert
- 语言：请使用中文（除非项目主要语言为英文）

请只输出commit message本身，不要有任何额外的解释。