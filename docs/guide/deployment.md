# 部署指南

## Docker 部署（推荐）

### 镜像选择

| 标签 (Tag) | 基础镜像 | 说明 |
| :--- | :--- | :--- |
| `latest` | Debian 12 | **默认推荐**：内置 Python 3.13 + Node.js 23，开箱即用 |
| `latest-debian13` | Debian 13 | 基于 Debian Trixie 的尝鲜版本 |
| `latest-minimal` | Debian 13 | **最小化版**：不预置语言环境，仅内置 Mise，极致纯净 |

::: tip 提示
默认使用 `latest` 标签。如需切换，将镜像名后的 `latest` 替换为 `latest-debian13` 或 `latest-minimal` 即可。
:::

### Docker Run

最简单的部署方式，一条命令即可启动：

```bash
docker run -d \
  --name xuanwu \
  -p 8052:8052 \
  -v $(pwd)/xuanwu/data:/app/data \
  -e TZ=Asia/Shanghai \
  --restart unless-stopped \
  ghcr.io/hicongcn/xuanwu:latest
```

启动后通过 `http://localhost:8052` 访问面板。详见 [快速开始](./getting-started.md) 获取默认账号信息。

### 完整环境变量配置

如需自定义更多配置，可通过环境变量覆盖默认值：

```bash
docker run -d \
  --name xuanwu \
  -p 8052:8052 \
  -v $(pwd)/xuanwu/data:/app/data \
  -e TZ=Asia/Shanghai \
  -e XW_SERVER_PORT=8052 \
  -e XW_SERVER_HOST=0.0.0.0 \
  -e XW_DB_TYPE=sqlite \
  -e XW_DB_PATH=/app/data/db/xuanwu.db \
  -e XW_DB_TABLE_PREFIX=xuanwu_ \
  --restart unless-stopped \
  ghcr.io/hicongcn/xuanwu:latest
```

### Docker Compose

推荐的生产环境部署方式：

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
    logging:
      driver: json-file
      options:
        max-size: "10m"
        max-file: "3"
    restart: unless-stopped
```

保存为 `docker-compose.yml` 后执行：

```bash
docker compose up -d
```

### MySQL 部署

默认使用 SQLite，如需使用 MySQL，设置 `XW_DB_TYPE=mysql` 并提供连接信息：

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
      - XW_DB_TYPE=mysql
      - XW_DB_HOST=mysql-server
      - XW_DB_PORT=3306
      - XW_DB_USER=root
      - XW_DB_PASSWORD=your_password
      - XW_DB_NAME=xuanwu
      - XW_DB_TABLE_PREFIX=xuanwu_
    restart: unless-stopped

  mysql:
    image: mysql:8.0
    container_name: xuanwu-mysql
    environment:
      - MYSQL_ROOT_PASSWORD=your_password
      - MYSQL_DATABASE=xuanwu
    volumes:
      - mysql_data:/var/lib/mysql
    restart: unless-stopped

volumes:
  mysql_data:
```

::: warning 注意
使用 MySQL 时，请确保 `mysql` 服务已就绪后再启动 `xuanwu` 服务，避免数据库连接失败。
:::

### 数据目录结构

容器启动后，所有持久化数据存储在挂载的 `/app/data` 目录下：

```
/data/
├── config/    配置文件 (config.ini)
├── db/        数据库文件 (xuanwu.db)
├── scripts/   用户脚本
├── bak/       备份文件
├── deps/      运行时依赖 (Mise)
└── log/       运行日志
```

### 自动更新镜像

使用 Watchtower 可实现镜像自动更新：

```bash
docker run -d \
  --name watchtower \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -e TZ=Asia/Shanghai \
  -e WATCHTOWER_SCHEDULE="0 0 3 * * *" \
  -e WATCHTOWER_CLEANUP=true \
  --restart unless-stopped \
  containrrr/watchtower \
  xuanwu
```

每天凌晨 3 点自动检查并更新名为 `xuanwu` 的容器。由于数据通过持久化挂载，更新不会造成数据丢失。

---

## 二进制部署

### 下载地址

从 [GitHub Releases](https://github.com/hicongcn/xuanwu-panel/releases) 下载对应平台的二进制文件。

### 支持平台

| 平台 | 文件名 |
| :--- | :--- |
| Linux amd64 | `xuanwu-linux-amd64` |
| Linux arm64 | `xuanwu-linux-arm64` |
| macOS amd64 | `xuanwu-darwin-amd64` |
| macOS arm64 | `xuanwu-darwin-arm64` |
| Windows amd64 | `xuanwu-windows-amd64.exe` |

### 启动服务

```bash
# Linux / macOS
chmod +x xuanwu-linux-amd64
./xuanwu-linux-amd64 server

# Windows
xuanwu-windows-amd64.exe server
```

服务启动后默认监听 `8052` 端口，访问 `http://localhost:8052` 进入面板。

---

## 源码编译

### 前置要求

- **Go** 1.25+
- **Node.js** 22+

### 编译步骤

```bash
git clone https://github.com/hicongcn/xuanwu-panel.git
cd xuanwu-panel
make release
```

编译产物位于 `bin/xuanwu`，启动方式：

```bash
./bin/xuanwu server
```

::: tip 提示
`make release` 会自动构建前端并嵌入到 Go 二进制文件中，生成的单文件即可独立运行。
:::
