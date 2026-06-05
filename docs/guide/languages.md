# 语言依赖

玄武面板深度集成了 **Mise** 运行时管理器，提供多版本语言环境的灵活管理和隔离。

## 脚本运行环境

玄武面板原生支持以下脚本的定时执行：
- **Python3**、**Node.js**、**Bash**（标准版镜像内置环境）

通过 **Mise** 可扩展支持几乎所有主流编程语言的动态安装与切换。

> [!TIP]
> **Minimal 镜像注意**：如果您使用的是 `minimal` 标签的镜像，系统初始不包含 Python 和 Node.js。您需要进入「编程语言」页面手动安装所需的运行时。

## Mise 插件列表

面板内置以下 25 个主流 Mise 插件，支持一键安装：

| 插件 | 语言/工具 | 插件 | 语言/工具 |
| :--- | :--- | :--- | :--- |
| `python` | Python | `node` | Node.js |
| `go` | Go | `rust` | Rust |
| `ruby` | Ruby | `java` | Java |
| `php` | PHP | `deno` | Deno |
| `bun` | Bun | `zig` | Zig |
| `dotnet` | .NET | `elixir` | Elixir |
| `erlang` | Erlang | `crystal` | Crystal |
| `lua` | Lua | `julia` | Julia |
| `nim` | Nim | `perl` | Perl |
| `scala` | Scala | `kotlin` | Kotlin |
| `clojure` | Clojure | `dart` | Dart |
| `flutter` | Flutter | `terraform` | Terraform |

> [!NOTE]
> 通过 `mise ls-remote <plugin>` 命令或面板界面，可以查看每个插件的所有可用版本。面板会自动获取最新 300 个版本供选择。

## 依赖管理

系统内置了高度集成的跨语言依赖管理器，支持以下语言的自动化依赖安装：

| 语言 | 包管理器 | 功能说明 |
| :--- | :--- | :--- |
| **Python** | pip | 使用内置虚拟环境，支持镜像源 |
| **Node.js** | npm | 全局安装模式，自动配置镜像 |
| **Go** | go install | 安装二进制工具到 GOPATH |
| **Rust** | cargo | 通过 `cargo install` 安装 |
| **Ruby** | gem | 支持 `gem install` 安装 |
| **Bun** | bun | 支持 `bun add -g` 全局模式 |
| **PHP** | composer | 支持 `composer global require` |
| **Deno** | deno | 支持 `deno install -g` |
| **.NET** | dotnet | 支持 `dotnet tool install -g` |
| **Elixir/Erlang** | mix | 支持 `mix archive.install` |
| **Lua** | luarocks | 通过 `luarocks` 管理 |
| **Nim** | nimble | 支持 `nimble install` |
| **Dart/Flutter** | pub | 支持 `pub global activate` |
| **Perl** | cpanm | 支持 `cpanm` 安装 |
| **Crystal** | shards | 支持 `shards` 安装 |

## 使用方法

### 1. 安装语言环境

进入「编程语言」页面，从插件列表中选择所需的语言及版本，点击安装即可。

### 2. 依赖管理

在已安装列表中点击「依赖管理」，输入包名称（可选版本号），系统会自动在对应语言环境中完成安装。

### 3. 多版本切换

对于复杂的项目，可以为不同任务配置不同的语言版本。系统基于 `mise exec` 实现环境隔离，不同版本的依赖包互不冲突。

### 4. 全局版本设置

支持将某个语言版本设置为全局默认版本（`mise use -g`），在未指定版本的任务中自动使用该版本。

## 内建助手库

玄武面板提供了内建的 Python 和 Node.js 助手库（`xuanwu` 包），用于简化通知推送、环境变量管理和任务控制等操作。

安装助手库到所有已安装的运行时版本：

```bash
xuanwu builtininstall
```

该命令会遍历所有 Mise 管理的 Python 和 Node.js 版本，在每个版本中安装内建助手库。

## 环境变量管理

Mise 支持设置全局环境变量，通过面板的「编程语言」页面可以：
- 查看当前全局环境变量列表
- 设置新的环境变量
- 删除不需要的环境变量

## 隔离机制说明

- 玄武面板通过动态注入 `PATH` 环境和 `mise shims` 将语言环境暴露给系统
- 每个任务在执行前根据配置自动加载对应的运行时环境变量
- `MISE_DATA_DIR` 等环境变量指向持久化目录，确保护版本数据在容器重启后不丢失
