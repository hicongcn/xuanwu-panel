# 内置库示例

玄武面板提供名为 `xuanwu` 的内建助手库（Built-in SDK），支持 Python 和 Node.js。通过该库可在脚本中实现**消息推送**、**环境变量管理**和**任务执行控制**等功能。

---

## 安装

### Python

```bash
pip install xuanwu
```

### Node.js

```bash
npm install xuanwu
```

也可以通过命令行工具一键为所有语言环境安装：

```bash
xuanwu builtininstall
```

---

## 环境变量配置

根据需要调用的功能，在任务的「环境变量」或「机密」中配置以下 Key：

### 消息推送

| 环境变量 | 说明 | 必填 |
|:---|:---|:---|
| `XWPKG_NOTIFY_TOKEN` | 通知 Token，在「消息推送」->「脚本调用说明」页面获取 | 是 |
| `XWPKG_NOTIFY_CHANNEL` | 渠道 ID，在「消息推送」->「渠道列表」页面查看 | 是 |
| `XWPKG_NOTIFY_URL` | 通知 API 地址，默认 `http://localhost:8052/api/v1/notify/send` | 否 |

### 环境变量管理 & 任务控制

| 环境变量 | 说明 | 必填 |
|:---|:---|:---|
| `XWPKG_OPENAPI_TOKEN`（或 `OPENAPI_TOKEN`） | OpenAPI 鉴权 Token，在「系统设置」->「OpenAPI」页面生成 | 是 |
| `XWPKG_OPENAPI_URL`（或 `OPENAPI_URL`） | OpenAPI 基础地址，默认自动推断 | 否 |

---

## 消息通知

一行代码即可触发推送。

::: code-group

```python [Python]
import xuanwu

xuanwu.notify(
    title="Python 任务提醒",
    text="这是一条来自 Python 脚本的通知消息"
)
```

```javascript [Node.js]
const xuanwu = require('xuanwu');

xuanwu.notify(
    "Node.js 任务提醒",
    "这是一条来自 Node.js 脚本的通知消息"
);
```

:::

---

## 环境变量管理

支持对面板环境变量进行增删改查。

### API 方法

| 功能 | Python | Node.js |
|:---|:---|:---|
| 获取所有变量 | `get_envs()` | `getEnvs()` |
| 按名称查询 | `get_env(name)` | `getEnv(name)` |
| 添加变量 | `add_env(name, value, remark)` | `addEnv(name, value, remark)` |
| 批量添加 | `add_envs(envs_list)` | `addEnvs(envsList)` |
| 更新变量 | `update_env(id, name, value, remark)` | `updateEnv(id, name, value, remark)` |
| 删除变量 | `delete_env(id)` | `deleteEnv(id)` |
| 批量删除 | `delete_envs(ids)` | `deleteEnvs(ids)` |

### 代码示例

::: code-group

```python [Python]
import xuanwu

def main():
    # 获取全部环境变量
    envs = xuanwu.get_envs()
    print(f"当前共有 {len(envs)} 个环境变量")

    # 添加环境变量
    created = xuanwu.add_env(
        name="MY_API_KEY",
        value="secret-value",
        remark="测试变量"
    )
    print(f"创建成功: ID={created.get('id')}")

    # 查询
    env = xuanwu.get_env("MY_API_KEY")
    if env:
        print(f"查询到: {env['name']} = {env['value']}")

        # 更新
        xuanwu.update_env(env["id"], env["name"], "new-value")

        # 删除
        xuanwu.delete_env(env["id"])
        print("已删除")

if __name__ == "__main__":
    main()
```

```javascript [Node.js]
const xuanwu = require('xuanwu');

async function main() {
    // 获取全部环境变量
    const envs = await xuanwu.getEnvs();
    console.log(`当前共有 ${envs.length} 个环境变量`);

    // 添加环境变量
    const created = await xuanwu.addEnv(
        "MY_API_KEY",
        "secret-value",
        "测试变量"
    );
    console.log(`创建成功: ID=${created.id}`);

    // 查询
    const env = await xuanwu.getEnv("MY_API_KEY");
    if (env) {
        console.log(`查询到: ${env.name} = ${env.value}`);

        // 更新
        await xuanwu.updateEnv(env.id, env.name, "new-value");

        // 删除
        await xuanwu.deleteEnv(env.id);
        console.log("已删除");
    }
}

main();
```

:::

---

## 定时任务管理

支持查询任务列表、触发执行和获取执行结果。

### API 方法

| 功能 | Python | Node.js |
|:---|:---|:---|
| 获取任务列表 | `get_tasks()` | `getTasks()` |
| 获取任务详情 | `get_task(id)` | `getTask(id)` |
| 触发执行 | `execute_task(id)` | `executeTask(id)` |
| 停止任务 | `stop_task(log_id)` | `stopTask(logId)` |
| 获取执行结果 | `get_last_results()` | `getLastResults()` |
| 更新任务 | `update_task(id, ...)` | `updateTask(id, ...)` |
| 删除任务 | `delete_task(id)` | `deleteTask(id)` |

### 代码示例

::: code-group

```python [Python]
import xuanwu

def main():
    # 获取所有任务
    tasks = xuanwu.get_tasks()
    print(f"共 {len(tasks)} 个任务:")
    for t in tasks[:5]:
        print(f"  - [{t.get('id')}] {t.get('name')}")

    # 手动触发第一个任务
    if tasks:
        xuanwu.execute_task(tasks[0].get("id"))
        print("执行指令已发送")

    # 获取最近执行结果
    results = xuanwu.get_last_results()
    print(f"最近 {len(results)} 条执行记录")

if __name__ == "__main__":
    main()
```

```javascript [Node.js]
const xuanwu = require('xuanwu');

async function main() {
    // 获取所有任务
    const tasks = await xuanwu.getTasks();
    console.log(`共 ${tasks.length} 个任务:`);
    tasks.slice(0, 5).forEach(t => {
        console.log(`  - [${t.id}] ${t.name}`);
    });

    // 手动触发第一个任务
    if (tasks.length > 0) {
        await xuanwu.executeTask(tasks[0].id);
        console.log("执行指令已发送");
    }

    // 获取最近执行结果
    const results = await xuanwu.getLastResults();
    console.log(`最近 ${results.length} 条执行记录`);
}

main();
```

:::
