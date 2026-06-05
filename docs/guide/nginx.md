# 反向代理

通过 Nginx 反向代理，可以为玄武面板配置域名、HTTPS 以及子路径部署。

## WebSocket 映射

在 `nginx.conf` 的 `http` 块中添加以下 WebSocket 升级映射（放在所有 `server` 块之外）：

```nginx
map $http_upgrade $connection_upgrade {
    default upgrade;
    '' close;
}
```

---

## 标准反向代理

适用于根路径部署（默认端口 `8052`）：

```nginx
server {
    listen 443 ssl http2;
    server_name example.com;

    ssl_certificate     /etc/letsencrypt/live/example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/example.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers off;

    access_log /var/log/nginx/xuanwu.access.log;
    error_log  /var/log/nginx/xuanwu.error.log warn;

    location / {
        proxy_pass http://127.0.0.1:8052;
        proxy_set_header Host $http_host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket 支持（终端、实时日志等功能必需）
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection $connection_upgrade;

        proxy_buffering off;
        proxy_read_timeout 60s;
    }
}

# HTTP 自动跳转 HTTPS（可选）
server {
    listen 80;
    server_name example.com;
    return 301 https://$server_name$request_uri;
}
```

---

## 子路径部署

当需要将玄武面板部署在域名的子路径下（如 `https://example.com/xuanwu/`）时：

### 1. 设置环境变量

```bash
XW_SERVER_URL_PREFIX=/xuanwu
```

或在 `config.ini` 中配置：

```ini
[server]
url_prefix = /xuanwu
```

### 2. Nginx 配置

```nginx
location /xuanwu/ {
    proxy_pass http://127.0.0.1:8052/xuanwu/;
    proxy_set_header Host $http_host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;

    # WebSocket 支持
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection $connection_upgrade;

    proxy_buffering off;
    proxy_read_timeout 60s;
}
```

::: warning
子路径部署时，`proxy_pass` 末尾的 `/` 不可省略，否则会导致路径转发异常。
:::

---

## 验证配置

```bash
# 检查 Nginx 配置语法
nginx -t

# 重新加载配置
nginx -s reload
```
