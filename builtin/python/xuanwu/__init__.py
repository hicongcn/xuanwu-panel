import os
from .notify import notify as _notify
from .env import (
    get_envs,
    get_env,
    add_envs,
    add_env,
    update_env,
    delete_envs,
    delete_env
)
from .task import (
    get_tasks,
    get_task,
    update_task,
    delete_task,
    execute_task,
    stop_task,
    get_last_results
)

def notify(title, text):
    """
    发送内建通知。
    会在调用时校验环境变量：XWPKG_NOTIFY_TOKEN, XWPKG_NOTIFY_CHANNEL
    """
    _TOKEN = os.environ.get("XWPKG_NOTIFY_TOKEN")
    _CHANNEL = os.environ.get("XWPKG_NOTIFY_CHANNEL")

    if not _TOKEN or not _CHANNEL:
        missing = []
        if not _TOKEN: missing.append("XWPKG_NOTIFY_TOKEN")
        if not _CHANNEL: missing.append("XWPKG_NOTIFY_CHANNEL")
        
        error_msg = f"缺少必要的环境变量以使用 xuanwu 模块: {', '.join(missing)}。请在玄武面板的任务设置中配置指定的 Key。"
        raise RuntimeError(error_msg)
    
    return _notify(title, text)

__all__ = [
    'notify',
    'get_envs',
    'get_env',
    'add_envs',
    'add_env',
    'update_env',
    'delete_envs',
    'delete_env',
    'get_tasks',
    'get_task',
    'update_task',
    'delete_task',
    'execute_task',
    'stop_task',
    'get_last_results'
]
