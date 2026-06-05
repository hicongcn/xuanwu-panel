# 命令行工具

玄武面板内置了同名的 `xuanwu` 命令行工具，支持在终端中执行系统级别的操作。

## 命令一览

| 命令 | 说明 |
| :--- | :--- |
| `xuanwu server` | 启动面板服务 |
| `xuanwu resetpwd` | 重置管理员密码 |
| `xuanwu reposync` | 同步远程 Git 仓库（支持 16+ 个参数） |
| `xuanwu restore <file>` | 从备份文件恢复数据 |
| `xuanwu builtininstall` | 安装内建助手库 |
| `xuanwu task` | 任务管理子命令集 |

使用 `xuanwu <命令> --help` 查看具体命令的参数说明。

---

## xuanwu server

启动玄武面板的后台服务进程。

```bash
# Docker 容器内自动启动
xuanwu server

# 手动部署时后台运行
nohup ./xuanwu server > /dev/null 2>&1 &
```

---

## xuanwu resetpwd

交互式重置管理员密码，密码丢失时可通过终端执行。

```bash
# Docker 容器内
docker exec -it xuanwu xuanwu resetpwd

# 直接执行
xuanwu resetpwd
```

根据提示输入新的管理员密码即可重置成功。

---

## xuanwu reposync

面板核心的仓库同步命令，支持从 Git 仓库或 URL 同步脚本到本地。

### 参数列表

| 参数名 | 默认值 | 说明 |
| :--- | :--- | :--- |
| `--source-type` | `git` | 同步源类型：`git` 或 `url` |
| `--source-url` | | 源地址（必填） |
| `--target-path` | | 目标路径（支持 `$SCRIPTS_DIR$` 变量替换） |
| `--branch` | | Git 分支名（留空自动检测默认分支） |
| `--path` | | 稀疏检出路径或单文件相对路径 |
| `--single-file` | `false` | 单文件模式 |
| `--proxy` | `none` | 代理类型：`none`/`ghproxy`/`mirror`/`custom` |
| `--proxy-url` | | 自定义代理地址 |
| `--auth-token` | | 认证 Token |
| `--http-proxy` | | HTTP 代理地址 |
| `--whitelist-paths` | | 白名单路径（逗号或竖线分隔） |
| `--blacklist` | | 黑名单关键字（竖线分隔） |
| `--dependence` | | 依赖文件关键字（竖线分隔） |
| `--extensions` | | 允许的文件后缀（竖线分隔） |
| `--commenttotask` | `false` | 启用青龙格式注释解析 |
| `--task-id` | | 内部任务 ID |
| `--task-langs` | | 任务语言配置（JSON） |
| `--repo-task-id` | | 原始任务 ID |
| `--task-timeout` | `30` | 超时时间（分钟） |
| `--pre-command` | | 同步前执行的命令 |
| `--post-command` | | 同步后执行的命令 |

### 使用示例

```bash
# 基础同步
xuanwu reposync --source-url https://github.com/example/repo.git --target-path $SCRIPTS_DIR$/repo1

# 启用代理并过滤后缀
xuanwu reposync --source-url https://github.com/example/repo.git \
  --target-path $SCRIPTS_DIR$/repo1 --proxy ghproxy --extensions ".js|.py"

# 稀疏检出
xuanwu reposync --source-url https://github.com/example/repo.git \
  --target-path $SCRIPTS_DIR$/repo1 --path "scripts/daily"

# 青龙格式注释解析
xuanwu reposync --source-url https://github.com/example/repo.git \
  --target-path $SCRIPTS_DIR$/repo1 --commenttotask true
```

---

## xuanwu restore

从本地 ZIP 备份文件恢复系统数据，全量覆盖现有数据库和脚本文件。

```bash
xuanwu restore /app/data/backup-20231027.zip
```

> [!WARNING]
> 该操作会全量覆盖现有数据库和脚本文件，请谨慎操作。恢复完成后可能需要重启服务。

---

## xuanwu builtininstall

为所有 Mise 管理的 Python 和 Node.js 环境安装内建助手库（`xuanwu` 包）。

```bash
xuanwu builtininstall
```

该命令会：
1. 检测系统中所有已安装的 Node.js 和 Python 版本
2. 在每个版本中安装对应的内建助手库
3. Node.js 使用 `npm i -g` 全局安装
4. Python 使用 `pip install --force-reinstall` 安装

---

## xuanwu task

任务管理子命令集，支持通过终端直接管理定时任务。

### 子命令列表

| 子命令 | 说明 |
| :--- | :--- |
| `list` | 查询并展示任务列表 |
| `run` | 手动立即触发执行任务 |
| `enable` | 启用任务 |
| `disable` | 禁用任务 |
| `status` | 查看任务最近一次执行状态 |
| `history` | 查看任务近期执行历史 |

所有子命令支持传入任务 ID、任务名称（精确/模糊匹配）或快捷字面量 `repo`。

### list - 任务列表

```bash
# 默认展示前 20 条
xuanwu task list

# 按名称筛选，查看第 2 页
xuanwu task list -name "签到" -page 2 -size 10

# 按类型筛选
xuanwu task list -type repo
```

参数：`-name`（名称关键词）、`-type`（任务类型）、`-page`（页码）、`-size`（每页条数）

### run - 立即执行

```bash
xuanwu task run <任务ID/名称/repo>

# 示例
xuanwu task run a1b2c3d4
xuanwu task run "自动签到"
xuanwu task run repo
```

### enable / disable - 启停控制

```bash
xuanwu task enable <任务ID/名称/repo>
xuanwu task disable <任务ID/名称/repo>
```

### status - 查看执行状态

```bash
# 查看最近一条日志
xuanwu task status <任务ID/名称/repo>

# 查看指定历史日志
xuanwu task status <任务ID> <日志ID>
```

### history - 执行历史

```bash
# 默认展示最近 10 条
xuanwu task history <任务ID/名称/repo>

# 展示最近 20 条
xuanwu task history repo -limit 20
```

> [!TIP]
> 所有 `xuanwu task` 子命令均支持 `--help` 参数获取详细的选项说明和示例。
