# 仓库同步

仓库同步允许玄武面板以 Git 仓库的形式管理和更新脚本库，支持大规模脚本的分发与自动化部署。

## 青龙指令兼容

如果您曾使用青龙面板，可以直接粘贴类似的仓库同步指令：

```
ql repo <url> <whitelist> <blacklist> <dependence> <branch>
```

系统将自动解析指令中的各项参数（仓库地址、白名单、黑名单、依赖文件、分支），并转换为玄武面板的仓库同步配置。

> [!IMPORTANT]
> **依赖管理说明**：由于玄武面板采用基于 Mise 的多版本语言管理系统，与青龙的全局环境不同，系统 **无法通过 `dependence` 字段自动安装依赖**。需要手动前往「语言依赖」页面或在终端中安装脚本所需的依赖包。

## Git 源管理

### 支持的代码托管平台

- **GitHub**
- **GitLab**
- **Gitee**

### 访问方式

| 方式 | 说明 |
| :--- | :--- |
| **公开仓库** | 无需认证，直接克隆 |
| **Token 认证** | 通过 `--auth-token` 参数或认证 Token 访问私有仓库 |
| **SSH 认证** | 在环境变量中配置 SSH 密钥 |

## 自动解析规则

同步完成后，系统自动扫描仓库中的脚本文件并解析任务配置：

### 脚本名解析

- `new Env('任务名称')`：解析 JavaScript 脚本中定义的任务展示名

### Cron 注释解析

- `cron "0 0 * * *"`：自动提取文件头部的 Cron 注释规则
- 支持青龙（QL）格式的脚本注释解析（通过 `--commenttotask true` 开启）

### 自动注册

解析到的 Cron 规则会自动创建对应的定时任务，并增量同步到主程序调度器。

## 白名单 / 黑名单过滤

同步时支持以下过滤规则：

| 参数 | 说明 | 分隔符 |
| :--- | :--- | :--- |
| **白名单路径** | 同步时受保护不被清理的路径，也用于脚本筛选 | 逗号或竖线 |
| **黑名单** | 包含关键字的文件将被过滤删除 | 竖线 `\|` |
| **依赖文件** | 包含关键字的文件强制保留 | 竖线 `\|` |
| **文件后缀** | 仅保留指定后缀的文件，不符则删除 | 竖线 `\|` |

关键字支持正则表达式匹配（默认忽略大小写），匹配失败时回退到全小写包含判断。

过滤执行顺序：依赖文件保留 → 后缀检查 → 黑名单删除 → 白名单检查

## 同步模式

### Git 同步

- **首次克隆**：`git clone --depth 1`，支持指定分支
- **增量更新**：`git fetch --all` + `git reset --hard origin/<branch>`
- **浅克隆**：默认使用 `--depth 1` 浅克隆，减少下载量

### URL 下载

支持通过 HTTP 直接下载单个文件，适用于 Raw 文件 URL（自动检测 GitHub/GitLab/Gitee 的 Raw URL 格式）。

## 高级功能

### 分支切换

通过 `--branch` 参数指定同步的 Git 分支，默认检测远程仓库的默认分支（main/master）。

### 稀疏检出

通过 `--path` 参数仅同步仓库中的特定子目录或文件：
```bash
xuanwu reposync --source-url <url> --target-path <path> --path "scripts/daily"
```

### 单文件模式

通过 `--single-file true` + `--path` 仅下载仓库中的指定单个文件。

### 代理加速

支持多种 GitHub 加速代理类型：

| 代理类型 | 说明 |
| :--- | :--- |
| `none` | 不使用代理 |
| `ghproxy` | 使用 gh-proxy.com 加速 |
| `mirror` | 使用 mirror.ghproxy.com 加速 |
| `custom` | 自定义代理地址 |

### 路径保护

同步时通过 `--whitelist-paths` 指定的路径会先被移到临时目录，同步完成后再恢复，确保用户数据不被覆盖。

### 前置/后置命令

同步前后可执行自定义命令，工作目录自动设置为仓库目录，并注入 `CURR_REPO_DIR` 环境变量指向仓库物理路径。

## CLI 命令

```bash
xuanwu reposync --source-url <仓库地址> --target-path <目标路径> [其他参数]
```

支持的参数（共 16 个）：`--source-type`、`--source-url`、`--target-path`、`--branch`、`--path`、`--single-file`、`--proxy`、`--proxy-url`、`--auth-token`、`--http-proxy`、`--whitelist-paths`、`--blacklist`、`--dependence`、`--extensions`、`--commenttotask`、`--task-id`、`--task-langs`、`--repo-task-id`、`--task-timeout`、`--pre-command`、`--post-command`

详见 [命令行工具 - reposync](cli.md#xuanwu-reposync)。
