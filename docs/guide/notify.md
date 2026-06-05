# 消息中心

消息中心集成了灵活的消息分发引擎，支持多种推送渠道和事件绑定规则，实现任务状态和系统事件的自动通知。

## 消息通道

面板支持以下消息推送渠道：

| 渠道 | 类型标识 | 说明 |
| :--- | :--- | :--- |
| **企业微信** | `qyweixin` | 企业微信群机器人 Webhook |
| **钉钉** | `dtalk` | 钉钉群机器人 Webhook |
| **飞书** | `feishu` | 飞书群机器人 Webhook |
| **Telegram** | `telegram` | Telegram Bot API |
| **Bark** | `bark` | Bark 推送服务（支持自建） |
| **VoceChat** | `vocechat` | VoceChat 推送（支持自建） |
| **邮件** | `email` | SMTP 邮件发送 |
| **自定义 Webhook** | `custom` | 自定义 HTTP POST 回调 |
| **PushPlus** | `pushplus` | PushPlus 推送服务 |
| **PushMe** | `pushme` | PushMe 推送服务 |
| **WxPusher** | `wxpusher` | 微信推送服务 |
| **Gotify** | `gotify` | Gotify 推送服务 |
| **ntfy** | `ntfy` | ntfy 推送服务 |
| **阿里云短信** | `aliyun_sms` | 阿里云短信服务 |

每个渠道支持独立的启停控制，可在渠道列表中进行配置、测试和删除。

## 事件绑定规则

通过事件绑定机制，将系统事件与推送渠道关联，实现自动化通知。

### 系统事件

| 事件 | 说明 |
| :--- | :--- |
| `user_login` | 用户登录（成功/失败） |
| `brute_force_login` | 密码多次错误（暴力破解检测） |
| `password_changed` | 密码修改 |

### 任务事件

| 事件 | 说明 |
| :--- | :--- |
| `task_success` | 任务执行成功 |
| `task_failed` | 任务执行失败 |
| `task_timeout` | 任务执行超时 |

任务事件支持绑定到特定任务（通过 `data_id` 指定任务 ID），实现不同任务使用不同的通知渠道。

## 推送路径

### 路径一：任务绑定通知

在任务编辑页面的「通知配置」中绑定推送渠道和触发条件：

1. 选择推送渠道（可多选）
2. 勾选触发时机：成功时 / 失败时 / 超时时
3. 可选开启日志推送，并设置日志截取长度（默认 1000 字符）

保存后，任务每次执行结束都会按配置自动推送通知。

### 路径二：内建 SDK 调用

通过玄武面板内建的 Python/Node.js 助手库（`xuanwu` 包）在脚本中主动发送通知。

**前置条件**：先执行 `xuanwu builtininstall` 安装助手库。

配置环境变量：
- `XWPKG_NOTIFY_TOKEN`：通知 Token（在消息中心的「脚本调用说明」中获取）
- `XWPKG_NOTIFY_CHANNEL`：目标渠道 ID
- `XWPKG_NOTIFY_URL`：通知 API 地址（默认 `http://localhost:8052/api/v1/notify/send`）

Python 示例：
```python
import xuanwu
xuanwu.notify("任务标题", "通知正文内容")
```

Node.js 示例：
```javascript
const xuanwu = require('xuanwu');
xuanwu.notify("任务标题", "通知正文内容");
```

### 路径三：原始 HTTP API

通过标准 HTTP POST 请求调用通知接口，适用于 Shell 或其他语言。

```bash
curl -X POST "http://localhost:8052/api/v1/notify/send" \
  -H "notify-token: 您的_NOTIFY_TOKEN" \
  -d '{"channel_id":"渠道ID", "title":"标题", "text":"内容"}'
```

认证方式：使用 `notify-token` Header 传递通知 Token（在系统设置中生成）。

## 消息前缀

所有推送消息的标题前会自动添加全局消息前缀，默认为 `[玄武面板]`。可在系统设置中自定义前缀内容。

## 消息模板

支持为每种事件类型自定义推送消息的标题和正文模板，使用 `{{变量名}}` 语法引用动态内容：

模板变量包括：`{{task_id}}`、`{{task_name}}`、`{{status}}`、`{{start_time}}`、`{{duration}}`、`{{error}}`、`{{output}}`、`{{username}}`、`{{ip}}` 等。

## 消息日志审计

所有通过面板发送的消息均会记录到推送日志中，包括：
- 发送时间
- 目标渠道名称和类型
- 消息标题和内容
- 发送结果（成功/失败）和错误信息

推送日志支持在「消息日志」页面查看，支持按条件筛选和自动清理。
