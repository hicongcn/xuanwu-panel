# 玄武面板 (Xuanwu Panel)

极致轻量、高性能的自动化任务调度平台。Go + Vue3 架构，深度集成 Mise 运行时管理，支持几乎所有主流编程语言的动态安装与依赖管理。

## 特色

- **轻量高性能** — 极低的 CPU 和内存占用
- **多语言支持** — 通过 Mise 动态安装 Python、Node.js、Go、Rust、PHP 等 15+ 语言
- **统一依赖管理** — 跨语言依赖安装、卸载、版本管理
- **任务调度** — 标准 Cron 表达式，支持并发控制和超时管理
- **脚本管理** — 在线代码编辑器，支持文件上传和压缩包解压
- **在线终端** — WebSocket 实时终端
- **消息推送** — 内置十余种主流渠道（企业微信、钉钉、飞书、Telegram、Bark 等）
- **PWA 支持** — 响应式设计，深色/浅色主题，移动端适配
- **仓库同步** — 兼容青龙仓库格式，自动同步脚本并创建任务

## 快速开始

### 方式一：Docker 部署（推荐）

```bash
docker run -d \
  --name xuanwu \
  -p 8052:8052 \
  -v $(pwd)/xuanwu/data:/app/data \
  -e TZ=Asia/Shanghai \
  --restart unless-stopped \
  ghcr.io/hicongcn/xuanwu:latest
```

启动后访问 `http://localhost:8052`，用户名 `admin`，密码查看日志：

```bash
docker logs xuanwu | grep "密码"
```

### 方式二：Docker Compose 部署

创建 `docker-compose.yml`：

```yaml
services:
  xuanwu:
    image: ghcr.io/hicongcn/xuanwu:latest
    container_name: xuanwu
    ports:
      - "8052:8052"
    volumes:
      - ./xuanwu/data:/app/data
    environment:
      - TZ=Asia/Shanghai
    restart: unless-stopped
```

```bash
docker compose up -d
docker compose logs xuanwu  # 查看管理员密码
```

### 方式三：下载二进制运行

从 [Releases](https://github.com/hicongcn/xuanwu-panel/releases/latest) 下载对应平台的二进制文件：

| 平台          | 文件                         |
| ----------- | -------------------------- |
| Linux amd64 | `xuanwu-linux-amd64`       |
| Linux arm64 | `xuanwu-linux-arm64`       |
| macOS amd64 | `xuanwu-darwin-amd64`      |
| macOS arm64 | `xuanwu-darwin-arm64`      |
| Windows     | `xuanwu-windows-amd64.exe` |

```bash
chmod +x xuanwu-linux-amd64
./xuanwu-linux-amd64 server
```

### 方式四：源码编译

**前置要求：** Go 1.25+, Node.js 22+

```bash
# 克隆仓库
git clone https://github.com/hicongcn/xuanwu-panel.git
cd xuanwu-panel

# 方式 A：一键构建（前端 + 后端，当前平台）
make release
./bin/xuanwu server

# 方式 B：全平台构建
make release-all
ls bin/release/

# 方式 C：仅后端开发模式（需要单独启动前端 dev server）
go run main.go server
```

### Docker 镜像标签

| 标签                | 基础镜像      | 说明                               |
| ----------------- | --------- | -------------------------------- |
| `latest`          | Debian 12 | 默认推荐，内置 Python 3.13 + Node.js 23 |
| `latest-debian13` | Debian 13 | 基于 Debian Trixie                 |
| `latest-minimal`  | Debian 13 | 最小化版，不含预装语言环境                    |

## 数据目录

所有用户数据统一存放在 `xuanwu/data/` 目录下（青龙风格），方便备份和管理：

```
xuanwu/data/
├── config/               # 配置文件
│   └── config.ini        # 应用配置（可选）
├── db/                   # 数据库
│   └── xuanwu.db         # SQLite 数据库
├── scripts/              # 脚本文件
├── bak/                  # 备份文件
├── deps/                 # 语言运行时
│   └── mise/             # Mise 管理的语言和依赖
└── log/                  # 运行日志
```

## 配置

支持环境变量和配置文件两种方式，环境变量优先级更高。

### 常用环境变量

| 环境变量             | 说明                   | 默认值              |
| ---------------- | -------------------- | ---------------- |
| `XW_SERVER_PORT` | 服务端口                 | 8052             |
| `XW_SERVER_HOST` | 监听地址                 | 0.0.0.0          |
| `XW_DB_TYPE`     | 数据库类型 (sqlite/mysql) | sqlite           |
| `XW_DB_HOST`     | 数据库地址                | localhost        |
| `XW_DB_PORT`     | 数据库端口                | 3306             |
| `XW_DB_USER`     | 数据库用户                | root             |
| `XW_DB_PASSWORD` | 数据库密码                | -                |
| `XW_DB_NAME`     | 数据库名称                | xw\_panel        |
| `XW_DB_PATH`     | SQLite 文件路径          | ./data/xuanwu.db |

### 配置文件

创建 `xuanwu/data/config/config.ini`：

```ini
[server]
port = 8052
host = 0.0.0.0

[database]
type = sqlite
path = ./data/db/xuanwu.db
table_prefix = xuanwu_
```

### MySQL 部署

```bash
docker run -d --name xuanwu -p 8052:8052 \
  -v $(pwd)/xuanwu/data:/app/data \
  -e TZ=Asia/Shanghai \
  -e XW_DB_TYPE=mysql \
  -e XW_DB_HOST=your-mysql-host \
  -e XW_DB_PORT=3306 \
  -e XW_DB_USER=root \
  -e XW_DB_PASSWORD=your-password \
  -e XW_DB_NAME=xuanwu \
  --restart unless-stopped \
  ghcr.io/hicongcn/xuanwu:latest
```

## CLI 命令

```bash
xuanwu server             # 启动面板服务
xuanwu resetpwd           # 重置管理员密码
xuanwu reposync           # 同步远程 Git 仓库
xuanwu restore <file>     # 从备份文件恢复数据
```

## 支持的语言

| 语言      | 安装方式        | 语言     | 安装方式     |
| ------- | ----------- | ------ | -------- |
| Python  | pip         | Ruby   | gem      |
| Node.js | npm         | Bun    | bun      |
| Go      | go install  | PHP    | composer |
| Rust    | cargo       | Deno   | deno     |
| .NET    | dotnet tool | Elixir | mix      |
| Lua     | luarocks    | Nim    | nimble   |
| Dart    | pub         | Perl   | cpanm    |
| Crystal | shards      | <br /> | <br />   |

## Nginx 反向代理

```nginx
server {
    listen 443 ssl http2;
    server_name example.com;

    ssl_certificate     /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    location / {
        proxy_pass http://127.0.0.1:8052;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

子路径部署：添加 `-e XW_SERVER_URL_PREFIX=/xuanwu` 环境变量，访问地址变为 `https://example.com/xuanwu/`。

## 致谢

本项目的开发离不开以下优秀的开源项目：

- **[baihu-panel](https://github.com/engigu/baihu-panel)** — 后端框架架构参考，部分代码基于白虎面板改进
- **[qinglong](https://github.com/whyour/qinglong)** — 功能设计参考，定时任务管理、环境变量、订阅管理等核心功能借鉴自青龙面板

感谢以上项目作者的贡献！

## 免责声明

玄武面板仅作为任务托管与调度平台，不提供任何第三方脚本。用户自行添加的脚本需自行审核安全性。本项目按"原样"提供，不保证无 Bug 或漏洞。

## 许可证

本项目采用 [Apache License 2.0](LICENSE) 协议。

***
