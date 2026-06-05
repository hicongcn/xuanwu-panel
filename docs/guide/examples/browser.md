# 浏览器示例

`example/playwright/` 提供了一组远程浏览器脚本示例，演示如何连接 **Browserless** 进行浏览器自动化。

> [!NOTE]
> 本示例以 Browserless 为演示对象。你也可以使用其他支持 CDP 协议的浏览器服务（如原生 `headless-shell`）。

---

## 部署方式对比

### 远程部署（推荐）

单独部署 Browserless 服务，由玄武中的脚本通过网络连接：

- **性能隔离**：浏览器的负载波动不会影响玄武面板的稳定性
- **开箱即用**：专业的浏览器镜像包含所有底层依赖和沙箱安全配置
- **可视化调试**：支持通过 VNC 同步查看浏览器画面
- **弹性伸缩**：支持多会话、多实例模式

### 本地部署（不推荐）

在玄武容器内直接安装浏览器：

- 浏览器极度消耗内存/CPU，容易导致 OOM
- Docker 镜像体积暴增 500MB+
- 精简镜像中安装浏览器常遇到缺失 `.so` 库文件的错误

---

## 准备工作

### 1. 安装语言依赖

前往「语言依赖」页面安装对应包：

| 脚本 | 依赖包 | 推荐版本 |
|:---|:---|:---|
| `playwright.js` | `puppeteer-core` | Node.js 环境 |
| `playwright.py` | `playwright` | Python 3.11 |

::: tip
Python 建议使用 **3.11** 版本，更高版本（如 3.12+）可能导致部分 Playwright 依赖编译失败。
:::

### 2. Browserless 连接要点

- 不要执行 `playwright install`（远程浏览器不需要本地安装）
- 不需要额外下载 Chromium / Firefox / WebKit
- 直接使用 `connect_over_cdp` 连接远程浏览器

---

## Docker Compose 部署

```yaml
version: "3.8"

services:
  browser:
    image: ghcr.io/browserless/chromium:latest
    container_name: browser
    restart: unless-stopped
    environment:
      MAX_CONCURRENT_SESSIONS: 5
      MAX_QUEUE_LENGTH: 20
      CONNECTION_TIMEOUT: 300000
      DEFAULT_LAUNCH_ARGS: '["--no-sandbox","--disable-setuid-sandbox","--disable-dev-shm-usage"]'
      TOKEN: your-secret-token
      ENABLE_DEBUGGER: "false"
    shm_size: "1gb"
    mem_limit: 2g
    cpus: 2
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:3000/healthz"]
      interval: 30s
      timeout: 10s
      retries: 3

  xuanwu:
    image: ghcr.io/hicongcn/xuanwu:latest
    container_name: xuanwu
    restart: unless-stopped
    ports:
      - "8052:8052"
    volumes:
      - ./data:/app/data
    environment:
      - TZ=Asia/Shanghai
      - XW_SERVER_PORT=8052
      - XW_SERVER_HOST=0.0.0.0
      - XW_DB_TYPE=sqlite
      - XW_DB_PATH=/app/data/xuanwu.db
      - XW_DB_TABLE_PREFIX=xuanwu_
    depends_on:
      - browser
```

---

## 代码示例

### Python (Playwright)

```python
from playwright.sync_api import sync_playwright

def main():
    with sync_playwright() as p:
        browser = p.chromium.connect_over_cdp(
            "http://browser:3000?token=your-secret-token"
        )
        page = browser.new_page()
        page.goto("https://www.baidu.com", wait_until="networkidle", timeout=30000)
        page.screenshot(path="baidu.png", full_page=True)
        print("截图已保存为 baidu.png")
        browser.close()

if __name__ == "__main__":
    main()
```

### Node.js (Puppeteer)

```javascript
const puppeteer = require('puppeteer-core');

(async () => {
  const browser = await puppeteer.connect({
    browserWSEndpoint: 'ws://browser:3000?token=your-secret-token'
  });
  const page = await browser.newPage();
  await page.goto('https://www.baidu.com', {
    waitUntil: 'networkidle2',
    timeout: 30000
  });
  await page.screenshot({ path: 'baidu.png', fullPage: true });
  console.log('截图完成：baidu.png');
  await browser.close();
})();
```

---

## 使用步骤

1. 启动 `browser` 和 `xuanwu` 服务
2. 在玄武「语言依赖」中安装对应包
3. 修改脚本中的 Browserless 地址和 Token
4. 在玄武中创建任务并运行脚本

## 常见问题

### 找不到模块

请确认已在「语言依赖」中安装了对应包（`puppeteer-core` 或 `playwright`）。

### 连接不上 Browserless

检查以下几点：

- Browserless 服务是否正常启动
- Token 是否与 Browserless 配置一致
- 玄武与 Browserless 是否在同一网络中
- 地址是否写成了当前运行环境可访问的地址（如 Docker 内部用服务名 `browser`，外部用 `localhost`）

### 页面打开超时

尝试：更换更稳定的目标站点、调大超时时间、增加 Browserless 容器的 `shm_size`、检查容器资源是否充足。
