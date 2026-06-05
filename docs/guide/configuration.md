# 系统配置

玄武面板支持通过**环境变量**和**配置文件**两种方式调整系统参数。配置文件存在时优先使用配置文件，不存在时自动回退到环境变量。

## 配置加载优先级

```
配置文件 (config.ini) > 环境变量 > 内置默认值
```

::: tip
启动时读取的环境变量会被自动清除（`Unsetenv`），避免在子进程环境中泄露敏感信息。
:::

---

## 环境变量配置

适用于 Docker 容器环境、CI/CD 流水线等场景。当配置文件不存在时，系统自动读取以下环境变量。

### 服务配置

| 环境变量 | 说明 | 默认值 |
|:---|:---|:---|
| `XW_SERVER_PORT` | 服务监听端口 | `8052` |
| `XW_SERVER_HOST` | 监听地址 | `0.0.0.0` |
| `XW_SERVER_URL_PREFIX` | URL 前缀，用于子路径部署 | 空 |
| `XW_SERVER_PPROF` | 是否启用 pprof 性能分析 | `false` |

### 数据库配置

| 环境变量 | 说明 | 默认值 |
|:---|:---|:---|
| `XW_DB_TYPE` | 数据库类型（`sqlite` / `mysql`） | `sqlite` |
| `XW_DB_HOST` | 数据库地址 | `localhost` |
| `XW_DB_PORT` | 数据库端口 | `3306` |
| `XW_DB_USER` | 数据库用户名 | `root` |
| `XW_DB_PASSWORD` | 数据库密码 | 空 |
| `XW_DB_NAME` | 数据库名称 | `github.com/hicongcn/xuanwu-panel` |
| `XW_DB_PATH` | SQLite 文件路径 | `./data/db/xuanwu.db` |
| `XW_DB_TABLE_PREFIX` | 数据库表前缀 | `xuanwu_` |
| `XW_DB_DSN` | 数据库 DSN（仅 MySQL，优先级高于上述独立参数） | 空 |

### 安全配置

| 环境变量 | 说明 | 默认值 |
|:---|:---|:---|
| `XW_SECRET` | JWT 密钥与密码盐值 | 空 |

### 其他

| 环境变量 | 说明 | 默认值 |
|:---|:---|:---|
| `XW_CONFIG_PATH` | 自定义配置文件路径 | `{DataDir}/config/config.ini` |
| `XW_DEMO_MODE` | 演示模式（`true` / `1` 启用，拦截所有任务执行） | `false` |

---

## 配置文件

如果需要更精细的控制，可以创建 INI 格式的配置文件。

### 文件路径

默认路径为 `{应用根目录}/data/config/config.ini`，可通过 `XW_CONFIG_PATH` 环境变量自定义。

### 目录结构

```
data/
├── config/
│   └── config.ini    # 配置文件
├── db/
│   └── xuanwu.db     # SQLite 数据库
├── scripts/           # 脚本工作目录
├── log/               # 日志目录
└── bak/               # 备份目录
```

### 配置文件示例

```ini
[server]
port = 8052
host = 0.0.0.0
# url_prefix = /xuanwu
# pprof_enabled = false

[database]
type = sqlite
# host = localhost
# port = 3306
# user = root
# password =
# dbname = xuanwu
path = ./data/db/xuanwu.db
table_prefix = xuanwu_

[security]
# secret = your-secret-here
```

### MySQL 配置示例

```ini
[database]
type = mysql
host = 127.0.0.1
port = 3306
user = root
password = your_password
dbname = xuanwu
table_prefix = xuanwu_
```

也可以使用 DSN 方式连接：

```ini
[database]
type = mysql
dsn = user:password@tcp(127.0.0.1:3306)/xuanwu?charset=utf8mb4&parseTime=True&loc=Local
table_prefix = xuanwu_
```

### Docker 挂载配置文件

```yaml
volumes:
  - ./data/config/config.ini:/app/data/config/config.ini
```

---

## 调度设置

任务调度器采用异步队列 + Worker Pool 架构，可在「系统设置 > 调度设置」页面进行在线调整。

| 参数 | 说明 | 默认值 |
|:---|:---|:---|
| Worker 数量 | 同时并发运行的任务进程数 | `4` |
| 队列大小 | 待处理任务队列的最大容量 | `100` |
| 速率间隔 | 两个任务启动之间的最小间隔（毫秒） | `200` |

### 工作原理

- 多个 Worker 协程从共享队列中消费任务
- 每次消费前经过速率限制器，确保任务间有最小间隔
- 队列满时，非严格模式下任务会降级为直接执行
- 支持运行时热重载配置（无需重启服务）

---

## 目录说明

玄武面板的根目录通过以下优先级确定：

1. 当前工作目录及其上级目录中查找 `data/config/config.ini`
2. 当前可执行文件路径及其上级目录中查找
3. 回退到当前工作目录

| 目录 | 说明 |
|:---|:---|
| `data/config/` | 配置文件存放目录 |
| `data/db/` | SQLite 数据库文件目录 |
| `data/scripts/` | 脚本工作目录 |
| `data/log/` | 运行日志目录 |
| `data/bak/` | 备份数据目录 |
| `data/deps/` | 语言依赖安装目录 |
| `data/bin/` | 程序软链接目录 |
